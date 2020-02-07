package main

import (
	"daosuan/core/cache"
	viewbase "daosuan/core/view"
	"daosuan/models/db"
	"daosuan/utils"
	"daosuan/utils/middlewares"
	"daosuan/views"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
	//"daosuan/views"
)

func initRouter(app *iris.Application) {
	views.RegisterAccountRouters(app)
	views.RegisterTagRouters(app)
	views.RegisterResourceRouters(app)
	views.RegisterProductRouters(app)
	views.RegisterSystemRouters(app)
}

func main() {
	app := iris.New()
	// 注册控制器
	app.UseGlobal(middlewares.AbnormalHandle)
	hero.Register(viewbase.DaoSuanView)
	initRouter(app)
	utils.InitGlobal()
	db.InitDB()
	cache.InitDijan()
	app.Run(iris.Addr(":80"), iris.WithoutServerError(iris.ErrServerClosed))
}
