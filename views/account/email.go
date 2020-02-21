package account

import (
	"daosuan/constants"
	"daosuan/core/cache"
	accountException "daosuan/exceptions/account"
	"daosuan/utils/hash"
	"daosuan/utils/mail"
	"daosuan/utils/params"
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
	if err := mailUtils.Send(target, token); err != nil {
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