package product

import (
	"daosuan/constants"
	"daosuan/core/auth"
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
	"github.com/kataras/iris"
)

// 创建产品
func CreateProduct(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))

	var t db.Product
	if err := db.Driver.Where("name = ?", params.Str("name", "名称")).First(&t).Error; err != nil || t.Id != 0 {
		panic(productException.NameIsExist())
	}
	product := db.Product{
		Name:        params.Str("name", "名称"),
		AuthorId:    auth.AccountModel().Id,
		Description: params.Str("description", "简介", ""),
		Details:     params.Str("details", "详情页"),
	}

	putProductInfo(params, &product, true)
	version := db.ProductVersion{
		ProductId: product.Id,
		Details: product.Details,
		Additional: product.Additional,
		VersionName: "v1.0.0",
	}
	db.Driver.Create(&version)
	ctx.JSON(iris.Map{
		"id": product.Id,
	})
}

// 获取产品信息
func GetProductInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {

	logic := productLogic.NewProductLogic(auth, pid)
	if !auth.IsAdmin() && logic.ProductModel().AuthorId != auth.AccountModel().Id && logic.ProductModel().Status != productEnums.StatusReleased {
		panic(productException.NoPermission())
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
	params.Diff(product)
	product.Name = params.Str("name", "名称")
	product.Description = params.Str("description", "简介", "")
	product.Details = params.Str("details", "详情页")

	putProductInfo(params, &product, false)
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
	ctx.JSON(iris.Map{
		"id": pid,
	})
}

// 获取产品列表
func GetProductList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	var lists []dto.ProductList
	var count int
	table := db.Driver.Table("product as p")
	if me, err := ctx.URLParamBool("me"); err == nil && me {
		table = table.Where("p.author_id = ?", auth.AccountModel().Id)
	} else {
		table = table.Where("p.status = ?", productEnums.StatusReleased)
	}
	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	// 条件过滤
	// TODO 分词搜索
	if name := ctx.URLParam("name"); len(name) > 0 {
		nameString := fmt.Sprintf("%%%s%%", name)
		table = table.Where("name like ?", nameString)
	}

	table.Count(&count).Offset((page - 1) * limit).
		Limit(limit).
		Select("p.id, p.update_time, p.cover, p.create_time, p.description, p.name, p.status").
		Order("update_time desc").Find(&lists)
	for i := 0; i < len(lists); i++ {
		lists[i].Cover = resourceLogic.GenerateToken(lists[i].Cover, -1, constants.DaoSuanSessionExpires)
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
	var products []db.Product
	db.Driver.Preload("Author").Table("product").Where("id in (?)", ids).Find(&products)
	for _, product := range products {
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

// 修改产品信息
func putProductInfo(params paramsUtils.ParamsParser, product *db.Product, create bool) {
	if params.Has("tag") {
		tagIds := params.List("tag", "标签id列表")
		var tags []db.Tag
		db.Driver.Where("id in (?)", tagIds).Find(&tags)
		product.Tag = tags
	}

	if params.Has("Status") {
		statusEnum := productEnums.NewStatusEnums()
		if statusEnum.Has(params.Int("status", "状态")) {
			product.Status = int16(params.Int("status", "状态"))
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
		if err := db.Driver.Create(&product).Error; err != nil {
			panic(productException.ProductCreateFail())
		}
	}
	if params.Has("cover") {
		logic := resourceLogic.NewReousrcesLocalStorage("product_cover")
		product.Cover = logic.SaveFile(fmt.Sprintf("%d/%s", product.Id, "cover.jpg"), durl.DataUrlParser(params.Str("cover", "封面")), true)
		if create {
			db.Driver.Save(&product)
		}
	}
	if !create {
		db.Driver.Save(&product)
	}
}
