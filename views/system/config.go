package system

import (
	"daosuan/core/auth"
	"daosuan/utils"
	"github.com/kataras/iris"
)

// 重载配置
func HeavyLoadConfig(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckAdmin()
	utils.InitGlobal()
	ctx.JSON(iris.Map {
		"status": "重载成功",
	})
}
