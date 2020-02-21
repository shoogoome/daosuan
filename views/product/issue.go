package product

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/enums/product"
	"daosuan/exceptions/product"
	"daosuan/logics/product"
	"daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils/hash"
	"daosuan/utils/params"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris"
)


// 发起提问
func CreateIssue(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	logic := productLogic.NewProductLogic(auth, pid)
	if logic.ProductModel().Status != productEnums.StatusReleased {
		panic(productException.NoPermission())
	}

	issue := db.Issue{
		AuthorId: auth.AccountModel().Id,
		ProductId: logic.ProductModel().Id,
		Title: params.Str("title", "标题"),
		Content: params.Str("content", "问题正文"),
	}

	if err := db.Driver.Create(&issue).Error; err != nil {
		panic(productException.CreateIssueFail())
	}
	ctx.JSON(iris.Map {
		"id": issue.Id,
	})
}

// 获取提问内容
func GetIssueInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, iid int) {

	var result []struct {
		db.Issue
		AuthorName string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
	}
	// 尝试读取缓存
	if re, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.DbModel, "issue", pid, iid)); err == nil && re != nil {
		if err := json.Unmarshal(re, &result); err == nil {
			ctx.JSON(result[0])
			return
		}
	}
	if err := db.Driver.
		Table("issue as i, product as p, account as a").
		Select("i.*, a.nickname as author_name, a.avator as author_avator").
		Where("i.product_id = p.id and i.author_id = a.id").
		Where("i.product_id = ? and i.id = ?", pid, iid).
		Find(&result).Error; err != nil || len(result) != 1 {
		panic(productException.IssueIsNotExists())
	}
	if len(result[0].AuthorAvator) > 0 {
		result[0].AuthorAvator = resourceLogic.GenerateToken(result[0].AuthorAvator, -1, -1)
	}
	// 缓存
	if re, err := json.Marshal(result); err == nil {
		v := hash.RandInt64(240, 240 * 5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.DbModel, "issue", pid, iid), re, int(v) * 60 * 60)
	}
	ctx.JSON(result[0])
}

// 删除提问
func DeleteIssue(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, iid int) {
	auth.CheckLogin()

	var issue db.Issue
	if err := db.Driver.Where("product_id = ? and id = ?", pid, iid).First(&issue).Error; err != nil {
		panic(productException.IssueIsNotExists())
	}

	if !auth.IsAdmin() && auth.AccountModel().Id != issue.AuthorId {
		panic(productException.NoPermission())
	}
	db.Driver.Delete(issue)
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.ProductIssueReplyModel, pid, iid))
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.DbModel, "issue", pid, iid))
	ctx.JSON(iris.Map {
		"id": iid,
	})
}

// 批量获取issue
func MgetIssue(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	ids := params.List("ids", "issue id列表")

	var result []struct {
		db.Issue
		AuthorName   string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
	}
	db.Driver.
		Table("issue as i, product as p, account as a").
		Select("i.*, a.nickname as author_name, a.avator as author_avator").
		Where("i.product_id = p.id and i.author_id = a.id").
		Where("i.product_id = ? and i.id in (?)", pid, ids).
		Find(&result)
	for i := 0; i < len(result); i++ {
		if len(result[i].AuthorAvator) > 0 {
			result[i].AuthorAvator = resourceLogic.GenerateToken(result[i].AuthorAvator, -1, -1)
		}
	}
	ctx.JSON(result)
}

// 获取issue列表
func GetIssueList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()

	var lists []struct{
		Id int `json:"id"`
		UpdateTime int64 `json:"update_time"`
		CreateTime int64 `json:"create_time"`
		Title string `json:"title"`
		ReplyNumber int `json:"reply_number"`
		AuthorName string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
	}
	var count int
	table := db.Driver.Debug().
		Select(`
			a.nickname as author_name, a.avator as author_avator,
			i.create_time as create_time, i.update_time as update_time, 
			i.title as title, i.id as id, i.reply_number as reply_number
`).
		Table("account as a, issue as i, product as p").
		Where("i.author_id = a.id and i.product_id = p.id and p.id = ?", pid)

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if key := ctx.URLParam("key"); len(key) > 0 {
		keyString := fmt.Sprintf("%%%s%%", key)
		table = table.Where("title like ?", keyString)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Order("-i.create_time").Find(&lists)
	for i := 0; i < len(lists); i++ {
		if len(lists[i].AuthorAvator) > 0 {
			lists[i].AuthorAvator = resourceLogic.GenerateToken(lists[i].AuthorAvator, -1, -1)
		}
	}

	ctx.JSON(iris.Map {
		"issues": lists,
		"total": count,
		"limit": limit,
		"page": page,
	})
}