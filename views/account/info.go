package account

import (
	authbase "daosuan/core/auth"
	accountException "daosuan/exceptions/account"
	accountLogic "daosuan/logics/account"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils/durl"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"fmt"
	"github.com/kataras/iris"
)

// 获取用户信息
func GetAccount(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {
	auth.CheckLogin()

	logic := accountLogic.NewAccountLogic(auth, aid)
	ctx.JSON(logic.GetAccountInfo())
}

// 获取用户列表
func GetAccountList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckAdmin()

	var lists []struct{
		Id int `json:"id"`
		UpdateTime int64 `json:"update_time"`
	}
	var count int
	table := db.Driver.Table("account")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	// 条件过滤
	if key := ctx.URLParam("key");len(key) > 0 {
		keyString := fmt.Sprintf("%%%s%%", key)
		table = table.Where(
			"nickname like ? or email like ? or realname like ?",
			keyString, keyString, keyString)
	}
	if idCode, err := ctx.URLParamBool("id_code"); err == nil {
		table = table.Where("id_code = ?", idCode)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, update_time").Find(&lists)
	ctx.JSON(iris.Map {
		"accounts": lists,
		"total": count,
		"limit": limit,
		"page": page,
	})
}

// 批量获取信息
func MgetAccountInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := accountLogic.NewAccountLogic(auth)

	ids := params.List("ids", "id列表")
	var data []interface{}
	var accounts []db.Account
	db.Driver.Table("account").Where("id in (?)", ids).Find(&accounts)
	for _, account := range accounts {
		logic.SetAccountModel(account)
		func(data *[]interface{}) {
			*data = append(*data, logic.GetAccountInfo())
			defer func() {
				recover()
			}()
		}(&data)
	}
	if len(data) == 0 {
		data = make([]interface{}, 0)
	}
	ctx.JSON(data)
}


// 修改用户信息
func PutAccount(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {
	auth.CheckLogin()
	data := paramsUtils.RequestJsonInterface(ctx)
	params := paramsUtils.NewParamsParser(data)
	logic := accountLogic.NewAccountLogic(auth, aid)
	account := logic.AccountModel()

	params.Diff(account)
	account.Nickname = params.Str("nickname", "昵称")
	account.Motto = params.Str("motto", "一句话签名")
	account.Realname = params.Str("realname", "真实姓名")
	account.IdCode = params.Str("id_code", "身份证")

	if params.Has("role") && auth.IsAdmin() {
		account.Role = int16(params.Int("role", "角色"))
	}

	if params.Has("new_password") {
		newPassword := params.Str("new_password", "旧密码")
		if !auth.IsAdmin() {
			oldPassword := params.Str("old_password", "旧密码")
			if hash.PasswordSignature(oldPassword) != account.Password {
				panic(accountException.OldPasswordIsNotTrue())
			}
		}
		account.Password = hash.PasswordSignature(newPassword)
	}
	// TODO: 特殊字段处理

	if params.Has("id_code") {
		// TODO: 身份证验证
	}

	if params.Has("avator") {
		logic := resourceLogic.NewReousrcesLocalStorage("account_avator")
		account.Avator = logic.SaveFile(fmt.Sprintf("%d/%s", account.Id, "avator.jpg"), durl.DataUrlParser(params.Str("avator", "头像")), true)
	}

	db.Driver.Save(&account)
	ctx.JSON(iris.Map {
		"id": account.Id,
	})
}



