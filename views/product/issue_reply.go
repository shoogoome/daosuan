package product

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/enums/product"
	productException "daosuan/exceptions/product"
	"daosuan/logics/product"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils/hash"
	"daosuan/utils/params"
	"encoding/json"
	"github.com/kataras/iris"
)

// 回复
func ReplyIssue(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid, iid int) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := productLogic.NewIssueLogic(auth, pid, iid)

	reply := db.IssueReply{
		AuthorId: auth.AccountModel().Id,
		IssueId: logic.IssueModel().Id,
		Reply: params.Str("reply", "回复"),
	}
	if params.Has("at_account_id") {
		aid := params.Int("at_account_id", "@回复人id")
		if aid == auth.AccountModel().Id {
			panic(productException.IsNotAllowAtSelf())
		}
		var atAccount db.Account
		if err := db.Driver.GetOne("account", aid, &atAccount); err != nil || atAccount.Id == 0 {
			panic(productException.AtAccountIsNotFound())
		}
		reply.AtAccountId = aid
	}
	tx := db.Driver.Begin()
	if err := tx.Create(&reply).Error; err != nil {
		tx.Rollback()
		panic(productException.ReplyFail())
	}
	logic.IssueModel().ReplyNumber += 1
	if err := tx.Save(logic.IssueModel()).Error; err != nil {
		tx.Rollback()
		panic(productException.ReplyFail())
	}
	tx.Commit()
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.ProductIssueReplyModel, pid, iid))
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.DbModel, "issue", pid, iid))
	ctx.JSON(iris.Map {
		"id": reply.Id,
	})
}

// 删除回复
func DeleteReply(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid, iid, rid int) {
	auth.CheckLogin()

	logic := productLogic.NewIssueLogic(auth, pid, iid)
	var reply db.IssueReply
	if err := db.Driver.GetOne("issue_Reply", rid, &reply); err != nil {
		panic(productException.ReplyIsNotExists())
	}
	if !auth.IsAdmin() && reply.AuthorId != auth.AccountModel().Id {
		panic(productException.NoPermission())
	}
	tx := db.Driver.Begin()
	if err := tx.Delete(reply).Error; err != nil {
		tx.Rollback()
		panic(productException.ReplyDeleteFail())
	}
	logic.IssueModel().ReplyNumber -= 1
	if err := tx.Save(logic.IssueModel()).Error; err != nil {
		tx.Rollback()
		panic(productException.ReplyDeleteFail())
	}
	tx.Commit()
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.ProductIssueReplyModel, pid, iid))
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.DbModel, "issue", pid, iid))
	ctx.JSON(iris.Map {
		"id": rid,
	})
}

// 获取回复列表
func GetReply(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid, iid int) {

	var result []struct{
		db.IssueReply
		AuthorName string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
		AtAccountName string `json:"at_account_name"`
	}
	// 读取缓存
	if re, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.ProductIssueReplyModel, pid, iid)); err == nil && re != nil {
		if err := json.Unmarshal(re, &result); err == nil {
			ctx.JSON(result)
			return
		}
	}
	db.Driver.
		Table("issue as i, product as p, account as a, issue_reply as r left join account as at on at.id = r.at_account_id").
		Select("r.*, a.nickname as author_name, a.avator as author_avator, at.nickname as at_account_name").
		Where("r.issue_id = i.id and p.id = i.product_id and p.id = ? and p.status = ?", pid, productEnums.StatusReleased).
		Where("r.author_id = a.id").
		Where("i.id = ?", iid).
		Order("create_time").
		Find(&result)

	for i := 0; i < len(result); i++ {
		if len(result[i].AuthorAvator) > 0 {
			result[i].AuthorAvator = resourceLogic.GenerateToken(result[i].AuthorAvator, -1, -1)
		}
	}

	// 缓存
	if re, err := json.Marshal(result); err == nil {
		v := hash.RandInt64(240., 240 * 5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.ProductIssueReplyModel, pid, iid), re, int(v) * 60 * 60)
	}
	ctx.JSON(result)
}

