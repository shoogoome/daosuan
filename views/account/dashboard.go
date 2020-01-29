package account

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/logics/account"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/models/dto"
	"github.com/kataras/iris"
)



// 个人空间
func Dashboard(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {

	logic := accountLogic.NewAccountLogic(auth, aid)

	var lists []dto.ProductList
	db.Driver.Table("product").
		Where("author_id = ?", aid).
		Select("p.id, p.update_time, p.cover, p.create_time, p.description, p.name, p.status").
		Order("update_time desc").Find(&lists)
	for i := 0; i < len(lists); i++ {
		lists[i].Cover = resourceLogic.GenerateToken(lists[i].Cover, -1, constants.DaoSuanSessionExpires)
	}

	ctx.JSON(iris.Map {
		"info": logic.GetAccountInfo(),
		"following": logic.GetFollowing(),
		"followers": logic.GetFollowers(),
		"product": lists,
		"stars": logic.GetStars(),
	})
}

