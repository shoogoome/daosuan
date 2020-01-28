package viewbase

import (
	authbase "daosuan/core/auth"
	"github.com/kataras/iris"
)

func DaoSuanView(ctx iris.Context) authbase.DaoSuanAuthAuthorization {
	return authbase.NewAuthAuthorization(&ctx)
}
