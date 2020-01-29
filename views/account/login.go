package account

import (
	authbase "daosuan/core/auth"
	accountException "daosuan/exceptions/account"
	accountLogic "daosuan/logics/account"
	"daosuan/models/db"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"github.com/kataras/iris"
)

// 普通登录
func CommonLogin(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	data := paramsUtils.RequestJsonInterface(ctx)
	params := paramsUtils.NewParamsParser(data)

	var account db.Account
	// function - 2 邮箱   - 1 电话
	if function, err := ctx.URLParamInt("function"); err == nil && function == 1 {
		db.Driver.Where(
			"password = ? and phone = ? and phone_validated = true",
			hash.PasswordSignature(params.Str("password", "密码")),
			params.Str("key", "帐号"),
		).First(&account)
	} else {
		db.Driver.Where(
			"password = ? and email = ? and email_validated = true",
			hash.PasswordSignature(params.Str("password", "密码")),
			params.Str("key", "帐号"),
		).First(&account)
	}

	if account.Id == 0 {
		panic(accountException.VerificationFail())
	}
	auth.SetSession(account.Id)
	if params.Has("remember") && params.Bool("remember", "记住登录") {
		auth.SetCookie(account.Id)
	}
	ctx.JSON(iris.Map {
		"id": account.Id,
	})
}

// 登出
func Logout(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.SetSession(0)
	auth.SetCookie(0)
	ctx.JSON(iris.Map {
		"status": "success",
	})
}

// 判断是否登录
func CheckLogin(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	if auth.IsLogin() {
		logic := accountLogic.NewAccountLogic(auth)
		logic.SetAccountModel(*auth.AccountModel())
		ctx.JSON(iris.Map {
			"account": logic.GetAccountInfo(),
			"status": true,
		})
	} else {
		ctx.JSON(iris.Map {
			"status": false,
		})
	}
}