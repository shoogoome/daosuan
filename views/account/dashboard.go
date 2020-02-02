package account

import (
	"daosuan/core/auth"
	"daosuan/logics/account"
	"github.com/kataras/iris"
)

// 个人空间
func Dashboard(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {

	logic := accountLogic.NewAccountLogic(auth, aid)

	ctx.JSON(iris.Map {
		"info": logic.GetAccountInfo(),
		"following": logic.GetFollowing(),
		"followers": logic.GetFollowers(),
		"product": logic.GetProduct(),
		"stars": logic.GetStars(),
	})
}

