package product

import (
	"daosuan/core/auth"
	"daosuan/entity"
	productException "daosuan/exceptions/product"
	"daosuan/logics/product"
	"daosuan/models/db"
	paramsUtils "daosuan/utils/params"
	"encoding/json"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/kataras/iris"
)

// 获取版本名称列表
func GetProductVersionList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {

	logic := productLogic.NewProductLogic(auth, pid)
	logic.LoadVersions()
	product := logic.ProductModel()

	l := make([]iris.Map, len(product.Versions))
	fmt.Println(l)
	for i := 0; i < len(product.Versions); i++ {
		l[i] = iris.Map {
			"name": product.Versions[i].VersionName,
			"id": product.Versions[i].Id,
		}
	}
	ctx.JSON(l)
}

// 创建版本
func CreateProductVersion(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()

	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := productLogic.NewProductLogic(auth, pid)
	product := logic.ProductModel()
	logic.CheckSelf()

	var version db.ProductVersion
	// 版本克隆
	if params.Has("version") {
		var srcVersion db.ProductVersion
		getVersion(pid, params.Int("version", "版本id"), &srcVersion)
		version.Details = srcVersion.Details
		version.Additional = srcVersion.Additional
		version.ProductId = srcVersion.ProductId
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
	versionName := params.Str("version_name", "版本名称")
	if logic.VersionIsExists(versionName) {
		panic(productException.ProductVersionNameIsExists())
	}
	version.VersionName = versionName
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
	product := logic.ProductModel()
	var version db.ProductVersion
	getVersion(pid, vid, &version)
	product.Details = version.Details
	product.Additional = version.Additional
	product.MasterVersion = version.VersionName
	db.Driver.Save(&product)
	ctx.JSON(iris.Map {
		"id": pid,
	})
}

// 检查产品名是否存在
func CheckVersionName(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, pid int) {
	auth.CheckLogin()
	name := ctx.URLParam("name")
	logic := productLogic.NewProductLogic(auth, pid)
	ctx.JSON(iris.Map {
		"exists": logic.VersionIsExists(name),
	})
}

// 获取版本实体
func getVersion(pid int, vid int, version *db.ProductVersion) {
	if err := db.Driver.Where("product_id = ? and id = ?", pid, vid).First(&version).Error; err != nil || version.Id == 0 {
		panic(productException.ProductVersionIsNotExists())
	}
}