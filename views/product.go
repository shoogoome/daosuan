package views

import (
	"daosuan/views/product"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func RegisterProductRouters(app *iris.Application) {

	// 产品路由
	productRouter := app.Party("/products")

	productRouter.Post("", hero.Handler(product.CreateProduct))
	productRouter.Get("/list", hero.Handler(product.GetProductList))
	productRouter.Post("/_mget", hero.Handler(product.MgetProduct))
	productRouter.Get("/{pid:int}", hero.Handler(product.GetProductInfo))
	productRouter.Put("/{pid:int}", hero.Handler(product.PutProduct))
	productRouter.Delete("/{pid:int}", hero.Handler(product.DeleteProduct))
	productRouter.Get("/{pid:int}/star", hero.Handler(product.Star))
	productRouter.Get("/{pid:int}/cancel_star", hero.Handler(product.CancelStar))

	// 账户版本路由
	productVersionRouter := app.Party("/products/{pid:int}/versions")
	productVersionRouter.Post("", hero.Handler(product.CreateProductVersion))
	productVersionRouter.Get("/check/version_name", hero.Handler(product.CheckVersionName))
	productVersionRouter.Get("/{vid:int}", hero.Handler(product.GetVersion))
	productVersionRouter.Put("/{vid:int}", hero.Handler(product.PutProductVersionInfo))
	productVersionRouter.Delete("/{vid:int}", hero.Handler(product.DeleteVersion))
	productVersionRouter.Get("/list", hero.Handler(product.GetProductVersionList))
	productVersionRouter.Get("/{vid:int}/set/master", hero.Handler(product.SetMaster))

}

