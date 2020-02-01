package product

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/enums/product"
	"daosuan/exceptions/product"
	"daosuan/logics/product"
	"daosuan/models/db"
	"daosuan/utils/params"
	"github.com/kataras/iris"
)

// 产品审核
func Examine(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckAdmin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := productLogic.NewProductLogic(auth, pid)
	product := logic.ProductModel()

	// 产品未提交请求
	if product.Status != int16(productEnums.StatusExamine) {
		panic(productException.ProductStatusIsNotExamine())
	}

	record := db.ProductExamineRecord{
		ProductId: product.Id,
		Adopt: params.Bool("adopt", "通过与否"),
		Reply: params.Str("reply", "回复", ""),
	}
	db.Driver.Create(&record)
	if record.Adopt {
		product.Status = productEnums.StatusReleased
	} else {
		product.Status = productEnums.StatusReject
	}
	db.Driver.Save(&product)
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.ProductExamineRecordModel, pid))
	ctx.JSON(iris.Map {
		"id": pid,
	})
}

// 获取产品审核信息列表
func GetExamineInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	logic := productLogic.NewProductLogic(auth, pid)
	if !auth.IsAdmin() && logic.ProductModel().AuthorId != auth.AccountModel().Id {
		panic(productException.NoPermission())
	}
	ctx.JSON(logic.GetExamineRecord())
}
