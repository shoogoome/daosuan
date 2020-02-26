package views

import (
	"daosuan/views/account"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"

	//"daosuan/views/account"
)

func RegisterAccountRouters(app *iris.Application) {

	// 账户路由
	accountRouter := app.Party("/accounts")

	accountRouter.Post("", hero.Handler(account.MgetAccountInfo))
	accountRouter.Get("/list", hero.Handler(account.GetAccountList))
	accountRouter.Post("/_mget", hero.Handler(account.MgetAccountInfo))
	accountRouter.Post("/register/email", hero.Handler(account.EmailRegister))
	accountRouter.Get("/logout", hero.Handler(account.Logout))
	accountRouter.Get("/unbind/oauth", hero.Handler(account.UnBindOauth))
	accountRouter.Post("/login/common", hero.Handler(account.CommonLogin))
	accountRouter.Get("/{aid:int}", hero.Handler(account.GetAccount))
	accountRouter.Put("/{aid:int}", hero.Handler(account.PutAccount))
	accountRouter.Get("/{aid:int}/dashboard", hero.Handler(account.Dashboard))
	accountRouter.Get("/check/login", hero.Handler(account.CheckLogin))
	accountRouter.Get("/check/nickname/{name:string}", hero.Handler(account.CheckNicknameExists))

	accountRouter.Get("/{aid:int}/follow", hero.Handler(account.Following))
	accountRouter.Get("/{aid:int}/cancel_follow", hero.Handler(account.CancelFollowing))

	// forget
	accountRouter.Post("/forget/email", hero.Handler(account.ForGetPasswordEmail))
	// oauth验证
	// github
	accountRouter.Get("/oauth/github/auth_url", hero.Handler(account.GitHubGetAuthUrl))
	accountRouter.Get("/oauth/github/callback", hero.Handler(account.GitHubCallback))
	// wechat
	accountRouter.Get("/oauth/wechat/auth_url", hero.Handler(account.WeChatGetAuthUrl))
	accountRouter.Get("/oauth/wechat/callback", hero.Handler(account.WeChatCallback))

	// 验证
	accountRouter.Get("/v/email/send", hero.Handler(account.SendMail))
	accountRouter.Post("/v/email", hero.Handler(account.VerificationMail))
}
