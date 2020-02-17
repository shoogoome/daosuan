package main

import (
	"daosuan/core/cache"
	"daosuan/core/queue"
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
	app.UseGlobal(middlewares.AbnormalHandle, middlewares.RequestLogHandle)
	hero.Register(viewbase.DaoSuanView)
	// 注册路由
	initRouter(app)
	// 初始化配置
	utils.InitGlobal()
	// 初始化数据库
	db.InitDB()
	// 初始化缓存
	cache.InitDijan()
	// 初始化任务队列
	queue.InitTaskQueue()
	// 启动系统
	app.Run(iris.Addr(":80"), iris.WithoutServerError(iris.ErrServerClosed))
}
