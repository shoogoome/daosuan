package views

import (
	"daosuan/views/resource"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func RegisterResourceRouters(app *iris.Application) {

	// 资源路由
	resourceRouter := app.Party("/resources")

	resourceRouter.Get("/local/{token:string}", hero.Handler(resource.LocalDownload))
	resourceRouter.Get("/qiniu/upload_token", hero.Handler(resource.GetQiNiuUploadToken))
}
