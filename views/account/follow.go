package account

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/exceptions/account"
	"daosuan/logics/account"
	"daosuan/models/db"
	"daosuan/utils/params"
	"github.com/kataras/iris"
)

// 关注
func Following(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {
	auth.CheckLogin()

	if aid == auth.AccountModel().Id {
		panic(accountException.FollowingFail())
	}

	var t db.AccountFollow
	if err := db.Driver.Where("source_id = ? and target_id = ?", auth.AccountModel().Id, aid).First(&t).Error; err == nil && t.Id != 0 {
		ctx.JSON(iris.Map {
			"id": t.Id,
		})
		return
	}

	logic := accountLogic.NewAccountLogic(auth, aid)
	follow := db.AccountFollow{
		SourceId: auth.AccountModel().Id,
		Target: *logic.AccountModel(),
	}
	if err := db.Driver.Create(&follow).Error; err == nil {
		// 清理缓存
		cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.FollowingModel, auth.AccountModel().Id))
		cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.FollowerModel, aid))
	} else {
		panic(accountException.FollowingFail())
	}

	ctx.JSON(iris.Map {
		"id": aid,
	})
}

// 取消关注
func CancelFollowing(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {
	auth.CheckLogin()

	if aid == auth.AccountModel().Id {
		panic(accountException.CancelFollowingFail())
	}

	var t db.AccountFollow
	if err := db.Driver.Where("source_id = ? and target_id = ?", auth.AccountModel().Id, aid).First(&t).Error; err == nil && t.Id != 0 {
		if err := db.Driver.Delete(&t).Error; err == nil {
			// 清理缓存
			cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.FollowingModel, auth.AccountModel().Id))
			cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.FollowerModel, aid))
		} else {
			panic(accountException.CancelFollowingFail())
		}
	}

	ctx.JSON(iris.Map {
		"status": "success",
	})
}

