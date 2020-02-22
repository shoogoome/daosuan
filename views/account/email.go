package account

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	accountException "daosuan/exceptions/account"
	"daosuan/models/db"
	"daosuan/utils/hash"
	"daosuan/utils/log"
	mailUtils "daosuan/utils/mail"
	paramsUtils "daosuan/utils/params"
	"github.com/kataras/iris"
)

// 发送邮件验证
// 注册时 绑定验证
// 忘记密码 验证
func SendMail(ctx iris.Context) {

	target := ctx.URLParam("target")
	// 邮箱验证
	if !mailUtils.CheckMailFormat(target) {
		panic(accountException.EmailFormatError())

	}

	// ip发送邮件频率检查
	ip :=  ctx.RemoteAddr()
	if re, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.AccountVerificationEmailTime, ip)); err == nil && re != nil {
		panic(accountException.EmailSendFrequently())
	}

	token := hash.GetRandomString(6)
	if err := mailUtils.Send(token, target); err != nil {
		logUtils.Println(err)
		panic(accountException.EmailSendFail())
	}
	if err := cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.AccountVerificationEmail, target), []byte(token), 60 * 60 * 2); err != nil {
		panic(accountException.EmailSendFail())
	}
	// 记住ip
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.AccountVerificationEmailTime, ip), []byte(ip), 60)
	ctx.JSON(iris.Map {
		"status": "success",
	})
}

// 验证邮箱
func VerificationMail(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	email := params.Str("email", "邮箱")
	token := params.Str("token", "验证码")

	if re, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.AccountVerificationEmail, email)); err == nil && re != nil {
		if token != string(re) {
			panic(accountException.OauthVerificationFail())
		}
	} else {
		panic(accountException.OauthVerificationFail())
	}
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.AccountVerificationEmail, email))

	account := auth.AccountModel()
	account.Email = email
	account.EmailValidated = true
	db.Driver.Save(&account)

	ctx.JSON(iris.Map {
		"id": account.Id,
	})
}