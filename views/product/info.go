package product

import (
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/entity"
	"daosuan/enums/product"
	"daosuan/exceptions/product"
	"daosuan/logics/product"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/models/dto"
	"daosuan/utils/durl"
	"daosuan/utils/params"
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"strings"
)

// 创建产品
func CreateProduct(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	if productLogic.IskNameExists(params.Str("name", "名称")) {
		panic(productException.NameIsExist())
	}
	product := db.Product{
		Name:          params.Str("name", "名称"),
		AuthorId:      auth.AccountModel().Id,
		Description:   params.Str("description", "简介", ""),
		Details:       params.Str("details", "详情页"),
		MasterVersion: "v1.0.0",
	}
	tx := db.Driver.Begin()
	putProductInfo(params, &product, true, tx)
	version := db.ProductVersion{
		ProductId:   product.Id,
		Details:     product.Details,
		Additional:  product.Additional,
		VersionName: "v1.0.0",
	}
	if err := tx.Create(&version).Error; err != nil {
		tx.Rollback()
		panic(productException.ProductCreateFail())
	}
	tx.Commit()
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.AccountProductModel, auth.AccountModel().Id))
	ctx.JSON(iris.Map{
		"id": product.Id,
	})
}

// 获取产品信息
func GetProductInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {

	logic := productLogic.NewProductLogic(auth, pid)
	if !auth.IsAdmin() && logic.ProductModel().AuthorId != auth.AccountModel().Id && logic.ProductModel().Status != productEnums.StatusReleased {
		panic(productException.ProductIsNotExists())
	}

	ctx.JSON(logic.GetProductInfo())
}

// 修改产品信息
func PutProduct(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()
	logic := productLogic.NewProductLogic(auth, pid)
	product := logic.ProductModel()
	if product.AuthorId != auth.AccountModel().Id {
		panic(productException.NoPermission())
	}

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	params.Diff(*product)
	// 检查名称冲突
	name := params.Str("name", "名称")
	var t db.Product
	if err := db.Driver.Where("name = ? and id != ?", name, product.Id).First(&t).Error; err == nil && t.Id != 0 {
		panic(productException.NameIsExist())
	}
	product.Name = name
	product.Description = params.Str("description", "简介", "")
	product.Details = params.Str("details", "详情页")

	tx := db.Driver.Begin()
	putProductInfo(params, product, false, tx)
	var version db.ProductVersion
	if err := db.Driver.Where("product_id = ? and version_name = ?", pid, product.MasterVersion).First(&version).Error; err == nil {
		version.Details = product.Details
		version.Additional = product.Additional
		if err := tx.Save(&version).Error; err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	} else {
		tx.Rollback()
	}
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.AccountProductModel, auth.AccountModel().Id))
	ctx.JSON(iris.Map{
		"id": product.Id,
	})
}

// 删除产品
func DeleteProduct(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()
	logic := productLogic.NewProductLogic(auth, pid)
	product := logic.ProductModel()
	if !auth.IsAdmin() && product.AuthorId != auth.AccountModel().Id {
		panic(productException.NoPermission())
	}
	if err := db.Driver.Delete(product).Error; err == nil {
		// 删除所有版本
		db.Driver.Exec("delete from `product_version` where product_id = ?", pid)
	}
	// 删除缓存
	cache.Dijan.Del(paramsUtils.CacheBuildKey(constants.AccountProductModel, auth.AccountModel().Id))
	ctx.JSON(iris.Map{
		"id": pid,
	})
}

// 获取产品列表
func GetProductList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	var lists []dto.ProductList
	var count int
	table := db.Driver.Table("product as p")
	// 可看自己产品的所有状态产品列表
	if me, err := ctx.URLParamBool("me"); err == nil && me {
		table = table.Where("p.author_id = ?", auth.AccountModel().Id)
		// 非系统管理员则只能看发布状态的产品列表
	} else if !auth.IsAdmin() {
		table = table.Where("p.status = ?", productEnums.StatusReleased)
	}
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	// 条件过滤
	// TODO elasticsearch 分词搜索
	if name := ctx.URLParam("name"); len(name) > 0 {
		nameString := fmt.Sprintf("%%%s%%", name)
		table = table.Where("name like ?", nameString)
	}
	if author, err := ctx.URLParamInt("author_id"); err == nil {
		table = table.Where("author_id = ?", author)
	}
	if status, err := ctx.URLParamInt("status"); err == nil && auth.IsAdmin() {
		table = table.Where("status = ?", status)
	}
	if tag := ctx.URLParam("tag"); len(tag) > 0 {
		tagList := strings.Split(tag, ",")
		table = table.Table("product as p, product_tags as t")
		table = table.Where("t.product_id = p.id and t.tag_id in (?)", tagList)
	}

	table.Count(&count).Offset((page - 1) * limit).
		Limit(limit).
		Select("distinct p.id, p.update_time, p.cover, p.create_time, p.description, p.name, p.status, p.star").
		Order("update_time desc").Find(&lists)
	for i := 0; i < len(lists); i++ {
		lists[i].Cover = resourceLogic.GenerateToken(lists[i].Cover, -1, -1)
	}
	ctx.JSON(iris.Map{
		"products": lists,
		"total":    count,
		"limit":    limit,
		"page":     page,
	})
}

// 批量获取产品信息
func MgetProduct(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := productLogic.NewProductLogic(auth)

	ids := params.List("ids", "id列表")
	var data []interface{}
	table := db.Driver.Preload("Author").Table("product")
	products := db.Driver.GetMany("product", ids, db.Product{}, table)
	for _, productInterface := range products {
		// 跳过非发布产品
		product := productInterface.(db.Product)
		if !auth.IsAdmin() && product.Status != productEnums.StatusReleased {
			continue
		}
		db.Driver.Model(&product).Related(&product.Tag, "Tag")
		logic.SetProductModel(product)
		func(data *[]interface{}) {
			*data = append(*data, logic.GetProductInfo())
			defer func() {
				recover()
			}()
		}(&data)
	}
	if len(data) == 0 {
		data = make([]interface{}, 0)
	}
	ctx.JSON(data)
}

// 检测产品名是否存在
func CheckName(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, name string) {
	auth.CheckLogin()
	ctx.JSON(iris.Map{
		"status": productLogic.IskNameExists(name),
	})
}

// 修改产品信息
func putProductInfo(params paramsUtils.ParamsParser, product *db.Product, create bool, tx ...*gorm.DB) {
	defer func() {
		if err := recover(); err != nil {
			if len(tx) > 0 {
				tx[0].Rollback()
			}
			panic(err)
		}
	}()
	d := db.Driver.DB
	if len(tx) > 0 {
		d = tx[0]
	}
	if params.Has("tag") {
		tagIds := params.List("tag", "标签id列表")
		var tags []db.Tag
		d.Where("id in (?)", tagIds).Find(&tags)
		product.Tag = tags
	}

	if params.Has("status") {
		statusEnum := productEnums.NewStatusEnums()
		status := params.Int("status", "状态")
		// 该方法不提供修改状态为发布
		if statusEnum.Has(status, productEnums.StatusDraft, productEnums.StatusLowerShelf, productEnums.StatusExamine) {
			product.Status = int16(params.Int("status", "状态"))
		} else {
			product.Status = productEnums.StatusDraft
		}
	} else {
		product.Status = productEnums.StatusDraft
	}

	if params.Has("additional") {
		var additionalEntity entity.ProductAdditional
		additional := params.Map("additional", "附加页")
		if err := mapstructure.Decode(additional, &additionalEntity); err == nil {
			if payload, err := json.Marshal(additionalEntity); err == nil {
				product.Additional = string(payload)
			}
		}
	}
	if create {
		if err := d.Create(&product).Error; err != nil {
			panic(productException.ProductCreateFail())
		}
	}
	if params.Has("cover") {
		logic := resourceLogic.NewReousrcesLocalStorage("product_cover")
		product.Cover = logic.SaveFile(fmt.Sprintf("%d/%s", product.Id, "cover.jpg"), durl.DataUrlParser(params.Str("cover", "封面")), true)
		if create {
			d.Save(&product)
		}
	}

	if !create {
		d.Save(&product)
	}
}
