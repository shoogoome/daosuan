package product

import (
	"daosuan/core/auth"
	"daosuan/entity"
	productException "daosuan/exceptions/product"
	"daosuan/logics/product"
	"daosuan/models/db"
	paramsUtils "daosuan/utils/params"
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"github.com/kataras/iris"
)

// 获取版本名称列表
func GetProductVersionList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {

	logic := productLogic.NewProductLogic(auth, pid)
	logic.LoadVersions()

	list := make([]struct{
		name string
		id int
	}, len(logic.ProductModel().Versions))
	for index, version := range logic.ProductModel().Versions {
		list[index] = struct {
			name string
			id   int
		}{name: version.VersionName, id: version.Id}
	}
	ctx.JSON(list)
}

// 创建版本
func CreateProductVersion(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := productLogic.NewProductLogic(auth, pid)
	product := logic.ProductModel()
	logic.CheckSelf()

	var version db.ProductVersion
	versionName := params.Str("version_name", "版本名称")
	var t db.ProductVersion
	if err := db.Driver.Where("product_id = ? and version_name = ?", pid, versionName).First(&t).Error; err != nil || t.Id != 0 {
		panic(productException.ProductVersionNameIsExists())
	}
	// 版本克隆
	if params.Has("version") {
		var srcVersion db.ProductVersion
		getVersion(pid, params.Int("version", "版本id"), &srcVersion)
		version = srcVersion
		version.VersionName = versionName
	// 自行写入或者同当前展示分支
	} else {
		params.Diff(product)
		version = db.ProductVersion{
			ProductId: product.Id,
			Details: params.Str("details", "详情页"),
		}
		if params.Has("additional") {
			var additionalEntity entity.ProductAdditional
			additional := params.Map("additional", "附加页")
			if err := mapstructure.Decode(additional, &additionalEntity); err == nil {
				if payload, err := json.Marshal(additionalEntity); err == nil {
					version.Additional = string(payload)
				}
			}
		} else {
			version.Additional = product.Additional
		}
	}
	db.Driver.Create(&version)
	ctx.JSON(iris.Map {
		"id": version.Id,
	})
}

// 修改版本信息
func PutProductVersionInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, vid int) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := productLogic.NewProductLogic(auth, pid)
	logic.CheckSelf()
	var version db.ProductVersion
	getVersion(pid, vid, &version)

	version.Details = params.Str("details", "详情页", version.Details)
	if params.Has("additional") {
		var additionalEntity entity.ProductAdditional
		additional := params.Map("additional", "附加页")
		if err := mapstructure.Decode(additional, &additionalEntity); err == nil {
			if payload, err := json.Marshal(additionalEntity); err == nil {
				version.Additional = string(payload)
			}
		}
	}
	db.Driver.Save(&version)
	ctx.JSON(iris.Map {
		"id": version.Id,
	})
}

// 删除版本
func DeleteVersion(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, vid int) {
	auth.CheckLogin()

	logic := productLogic.NewProductLogic(auth, pid)
	logic.CheckSelf()
	var version db.ProductVersion
	getVersion(pid, vid, &version)
	db.Driver.Delete(version)
	ctx.JSON(iris.Map {
		"id": vid,
	})
}

// 获取版本信息
func GetVersion(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, vid int) {
	logic := productLogic.NewProductLogic(auth, pid)
	var version db.ProductVersion
	getVersion(pid, vid, &version)
	logic.ProductModel().Details = version.Details
	logic.ProductModel().Additional = version.Additional
	logic.ProductModel().CreateTime = version.CreateTime
	logic.ProductModel().UpdateTime = version.UpdateTime
	data := logic.GetProductInfo()
	data.(map[string]interface{})["version_name"] = version.VersionName
	ctx.JSON(data)
}

// 设置master分支
func SetMaster(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int, vid int) {
	auth.CheckLogin()
	logic := productLogic.NewProductLogic(auth, pid)
	logic.CheckSelf()
	var version db.ProductVersion
	getVersion(pid, vid, &version)
	logic.ProductModel().Details = version.Details
	logic.ProductModel().Additional = version.Additional
	db.Driver.Save(logic.ProductModel())
	ctx.JSON(iris.Map {
		"id": pid,
	})
}

// 检查产品名是否存在
func CheckVersionName(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()
	nickname := ctx.URLParam("nickname")
	logic := productLogic.NewProductLogic(auth, pid)
	ctx.JSON(iris.Map {
		"exists": logic.VersionIsExists(nickname),
	})
}

// 获取版本实体
func getVersion(pid int, vid int, version *db.ProductVersion) {
	if err := db.Driver.Where("product_id = ? and id = ?", pid, vid).First(&version).Error; err != nil || version.Id == 0 {
		panic(productException.ProductVersionIsNotExists())
	}
}