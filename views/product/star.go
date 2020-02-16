package product

import (
	"bytes"
	"daosuan/constants"
	authbase "daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/core/elasticsearch"
	"daosuan/exceptions/product"
	productLogic "daosuan/logics/product"
	"daosuan/models/db"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"encoding/json"
	"github.com/kataras/iris"
)

// Star
func Star(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()
	logic := productLogic.NewProductLogic(auth, pid)
	if logic.IsStar() {
		panic(productException.ReStar())
	}
	product := logic.ProductModel()
	// 写入数据库
	tx := db.Driver.Begin()
	product.Star += 1
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		panic(productException.StarFail())
	}

	star := db.AccountStar{
		AccountId: auth.AccountModel().Id,
		ProductId: pid,
	}
	if err := tx.Create(&star).Error; err != nil {
		tx.Rollback()
		panic(productException.StarFail())
	}
	tx.Commit()
	// 缓存
	v := hash.RandInt64(240, 240 * 5)
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.StarModel, pid, auth.AccountModel().Id), []byte("star"), int(v) * 60 * 60)
	// 修改索引信息
	doc := map[string]interface{}{
		"doc": product,
	}
	if re, err := json.Marshal(doc); err == nil {
		elasticsearch.Update("product", "root", product.Id, bytes.NewBuffer(re))
	}
	ctx.JSON(iris.Map {
		"id": pid,
	})
}

// 取消Star
func CancelStar(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()
	logic := productLogic.NewProductLogic(auth, pid)
	if !logic.IsStar() {
		ctx.JSON(iris.Map {
			"id": pid,
		})
		return
	}

	// 删除数据
	product := logic.ProductModel()
	tx := db.Driver.Begin()
	product.Star -= 1
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		panic(productException.CancelStarFail())
	}

	if err := tx.Exec("delete from `account_star` where product_id = ? and account_id = ?", pid, auth.AccountModel().Id).Error; err != nil {
		tx.Rollback()
		panic(productException.CancelStarFail())
	}
	tx.Commit()

	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.StarModel, pid, auth.AccountModel().Id))
	// 修改索引信息
	doc := map[string]interface{}{
		"doc": product,
	}
	if re, err := json.Marshal(doc); err == nil {
		elasticsearch.Update("product", "root", product.Id, bytes.NewBuffer(re))
	}
	ctx.JSON(iris.Map {
		"id": pid,
	})
}
