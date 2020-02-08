package views

import (
	"daosuan/views/system"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"

	//"daosuan/views/account"
)

func RegisterSystemRouters(app *iris.Application) {

	// 账户路由
	systemRouter := app.Party("/system")

	// 重建缓存
	systemRouter.Get("/reset_cache", hero.Handler(system.ResetCache))
	// 重载配置
	systemRouter.Get("/load_config", hero.Handler(system.HeavyLoadConfig))
}
