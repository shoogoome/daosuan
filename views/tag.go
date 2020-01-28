package views

import (
	"daosuan/views/tag"
	"github.com/kataras/iris"
	"github.com/kataras/iris/hero"
)

func RegisterTagRouters(app *iris.Application) {

	tagRouter := app.Party("/tags")

	tagRouter.Post("", hero.Handler(tag.CreateTag))
	tagRouter.Get("/list", hero.Handler(tag.GetTagList))
	tagRouter.Delete("/{tid:int}", hero.Handler(tag.DeleteTag))
}
