package account

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	accountException "daosuan/exceptions/account"
	"daosuan/models/db"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"github.com/kataras/iris"
)

// 忘记密码（邮箱验证）
func ForGetPasswordEmail(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	email := params.Str("email", "邮箱")
	password := params.Str("password", "密码")
	token := params.Str("token", "验证码")

	if re, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.AccountVerificationEmail, email)); err == nil && re != nil {
		if token != string(re) {
			panic(accountException.OauthVerificationFail())
		}
	} else {
		panic(accountException.OauthVerificationFail())
	}
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.AccountVerificationEmail, email))

	var account db.Account
	if err := db.Driver.Where("email = ?", email).First(&account).Error; err != nil || account.Id == 0 {
		panic(accountException.AccountIsNotExists())
	}
	account.Password = hash.PasswordSignature(password)
	db.Driver.Save(&account)
	ctx.JSON(iris.Map {
		"id": account.Id,
	})
}
