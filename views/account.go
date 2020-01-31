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
	accountRouter.Post("/register", hero.Handler(account.Register))
	accountRouter.Get("/logout", hero.Handler(account.Logout))
	accountRouter.Post("/login/common", hero.Handler(account.CommonLogin))
	accountRouter.Get("/{aid:int}", hero.Handler(account.GetAccount))
	accountRouter.Put("/{aid:int}", hero.Handler(account.PutAccount))
	accountRouter.Get("/{aid:int}/dashboard", hero.Handler(account.Dashboard))
	accountRouter.Get("/check/login", hero.Handler(account.CheckLogin))
	accountRouter.Get("/check/nickname/{name:string}", hero.Handler(account.CheckNicknameExists))

	accountRouter.Get("/{aid:int}/follow", hero.Handler(account.Following))
	accountRouter.Get("/{aid:int}/cancel_follow", hero.Handler(account.CancelFollowing))

}
