package viewbase

import (
	authbase "daosuan/core/auth"
	"github.com/kataras/iris"
)

func DaoSuanView(ctx iris.Context) authbase.DaoSuanAuthAuthorization {
	if model := ctx.GetHeader("X-DaoSuan-Auth-Model"); len(model) > 0 && model == "client" {
		return authbase.NewClientAuthAuthorization(&ctx)
	}
	return authbase.NewAuthAuthorization(&ctx)
}
