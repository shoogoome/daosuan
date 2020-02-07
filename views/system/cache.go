package system

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/logics/account"
	"daosuan/logics/product"
	"daosuan/utils/hash"
	"daosuan/utils/params"
	"fmt"
	"github.com/shoogoome/godijan"

	"daosuan/models/db"
	"encoding/json"
	"github.com/kataras/iris"
)

// 缓存重建
func ResetCache(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckAdmin()

	go func() {
		// 账户表缓存
		fmt.Println("[*] 开始重建账户表及账户附属信息缓存...")
		var accounts []db.Account
		db.Driver.Find(&accounts)

		cmds := make([]*godijan.Cmd, 0, len(accounts)*2)
		for _, account := range accounts {
			if re, err := json.Marshal(account); err == nil {
				// 添加缓存指令(本体缓存)
				v := hash.RandInt64(240, 240*5)
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.DbModel, "account", account.Id),
					Value: re,
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
				// 昵称缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.NicknameModel, account.Nickname),
					Value: []byte(account.Nickname),
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
			}
			// 账户附加信息缓存
			logic := accountLogic.NewAccountLogic(auth)
			logic.SetAccountModel(account)
			logic.GetStars()
			logic.GetFollowers()
			logic.GetFollowing()
			logic.GetProduct()
		}
		// 写入缓存
		cache.Dijan.PipelinedRun(cmds)
		fmt.Println("[*] 账户表及附属信息缓存完毕")

		// 产品表缓存
		fmt.Println("[*] 开始重建产品表及附属信息缓存...")
		var products []db.Product
		db.Driver.Preload("Author").Find(&products)

		cmds = make([]*godijan.Cmd, 0, len(products)*2)
		for _, product := range products {
			if re, err := json.Marshal(products); err == nil {
				v := hash.RandInt64(240, 240*5)
				// 表本体缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.DbModel, "product", product.Id),
					Value: re,
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
				// 产品名缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.ProductNameModel, product.Name),
					Value: []byte(product.Name),
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
			}
			// 产品附属信息缓存
			logic := productLogic.NewProductLogic(auth)
			logic.SetProductModel(product)
			logic.GetVersionInfo()
			logic.GetExamineRecord()
		}
		// 写入缓存
		cache.Dijan.PipelinedRun(cmds)
		fmt.Println("[*] 产品表及附属信息缓存完毕")

		// 用户产品点赞表缓存
		fmt.Println("[*] 开始重建用户点赞表及附属信息缓存...")
		var stars []db.AccountStar
		db.Driver.Find(&stars)

		cmds = make([]*godijan.Cmd, 0, len(stars) * 2)
		for _, star := range stars {
			if re, err := json.Marshal(star); err == nil {
				v := hash.RandInt64(240, 240*5)
				// 表本体缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.DbModel, "account_star", star.Id),
					Value: re,
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
				// 点赞缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.StarModel, star.ProductId, star.AccountId),
					Value: []byte("star"),
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
			}
		}
		// 写入缓存
		cache.Dijan.PipelinedRun(cmds)
		fmt.Println("[*] 产品表及附属信息缓存完毕")

		// 产品issue表缓存
		fmt.Println("[*] 开始重建issue表...")
		var issues []db.Issue
		db.Driver.Find(&issues)

		cmds = make([]*godijan.Cmd, 0, len(issues))
		for _, issue := range issues {
			if re, err := json.Marshal(issue); err == nil {
				v := hash.RandInt64(240, 240*5)
				// 表本体缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.DbModel, "issue", issue.Id),
					Value: re,
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
			}
		}
		// 写入缓存
		cache.Dijan.PipelinedRun(cmds)
		fmt.Println("[*] issue表缓存完毕")


		// 产品issue回复表缓存
		fmt.Println("[*] 开始问答回复表...")
		var replies []db.IssueReply
		db.Driver.Find(&replies)

		cmds = make([]*godijan.Cmd, 0, len(replies))
		for _, reply := range replies {
			if re, err := json.Marshal(reply); err == nil {
				v := hash.RandInt64(240, 240*5)
				// 表本体缓存
				cmds = append(cmds, &godijan.Cmd{
					Key:   paramsUtils.CacheBuildKey(constants.DbModel, "issue_reply", reply.Id),
					Value: re,
					Name:  "set",
					TTL:   int(v) * 60 * 60,
				})
			}
		}
		// 写入缓存
		cache.Dijan.PipelinedRun(cmds)
		fmt.Println("[*] 问答回复表缓存完毕")

		// 缓存结束
		fmt.Println("[*] 缓存结束")
	}()

	ctx.JSON(iris.Map{
		"status": "已启动后台缓存重建。详情可查看日志",
	})
}
