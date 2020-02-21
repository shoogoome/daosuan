package account

import (
	"daosuan/core/auth"
	"daosuan/exceptions/account"
	"daosuan/models/db"
	"github.com/kataras/iris"
)

// 解除第三方绑定
func UnBind(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()

	model, err := ctx.URLParamInt("model")
	if err != nil {
		panic(accountException.OperationFail())
	}
	var oauth db.AccountOauth
	if err := db.Driver.Where("model = ? and account_id = ?", model, auth.AccountModel().Id).First(&oauth).Error; err != nil {
		panic(accountException.OperationFail())
	}
	db.Driver.Delete(oauth)
	ctx.JSON(iris.Map {
		"status": "success",
	})
}
