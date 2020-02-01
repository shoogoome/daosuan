package resource

import (
	"daosuan/core/auth"
	"daosuan/utils/qiniu"
	"github.com/kataras/iris"
)

// 获取七牛上传token
func GetQiNiuUploadToken(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()

	ctx.JSON(iris.Map {
		"token": qiniuUtils.GetUploadToken(),
	})
}
