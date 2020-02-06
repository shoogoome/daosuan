package product

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/enums/product"
	"daosuan/exceptions/product"
	"daosuan/logics/product"
	"daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils/params"
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

	if err := db.Driver.Create(&issue); err != nil {
		panic(productException.CreateIssueFail())
	}
	ctx.JSON(iris.Map {
		"id": issue.Id,
	})
}

// 获取提问内容
func GetIssueInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, iid int) {

	result := struct {
		db.Issue
		ProductName string `json:"product_name"`
		AuthorName string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
	}{}
	if err := db.Driver.
		Table("issue as i, product as p, account as a").
		Select("i.*, a.nickname as author_name, a.avator as author_avator, p.name as product_name").
		Where("i.product_id = p.id and i.author_id = a.id").
		Where(
			"i.product_id = ? and i.id = ?",
			pid, iid).
		First(&result).Error; err != nil || result.Id == 0 {
		panic(productException.IssueIsNotExists())
	}
	if len(result.AuthorAvator) > 0 {
		result.AuthorAvator = resourceLogic.GenerateToken(result.AuthorAvator, -1, constants.DaoSuanSessionExpires)
	}
	ctx.JSON(result)
}

// 删除提问
func DeleteIssue(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, iid int) {
	auth.CheckLogin()

	var issue db.Issue
	if err := db.Driver.Where("product_id = ? and id = ?", pid, iid).First(&issue).Error; err != nil || issue.Id == 0 {
		panic(productException.IssueIsNotExists())
	}

	if !auth.IsAdmin() && auth.AccountModel().Id != issue.AuthorId {
		panic(productException.NoPermission())
	}

	db.Driver.Delete(issue)
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
		ProductName  string `json:"product_name"`
		AuthorName   string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
	}
	db.Driver.
		Table("issue as i, product as p, account as a").
		Select("i.*, a.nickname as author_name, a.avator as author_avator, p.name as product_name").
		Where("i.product_id = p.id and i.author_id = a.id").
		Where("i.product_id = ? and i.id in (?)", pid, ids).
		Find(&result)
	for i := 0; i < len(result); i++ {
		if len(result[i].AuthorAvator) > 0 {
			result[i].AuthorAvator = resourceLogic.GenerateToken(result[i].AuthorAvator, -1, constants.DaoSuanSessionExpires)
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
		Title int16 `json:"title"`
		ReplyNumber int `json:"reply_number"`
		ProductName string `json:"product_name"`
		AuthorName string `json:"author_name"`
		AuthorAvator string `json:"author_avator"`
	}
	var count int
	table := db.Driver.
		Select(`
			a.nickname as author_name, a.avator as author_avator,
			p.name as product_name, i.create_time, i.update_time,
			i.title, i.id, i.reply_number
`).
		Table("account as a, issue as i, product as p").
		Where("i.author_id = a.id, i.product_id = p.id")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	if key := ctx.URLParam("key"); len(key) > 0 {
		keyString := fmt.Sprintf("%%%s%%", key)
		table = table.Where("title like ?", keyString)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Order("-create_time").Find(&lists)
	ctx.JSON(iris.Map {
		"issues": lists,
		"total": count,
		"limit": limit,
		"page": page,
	})
}