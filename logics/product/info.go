package productLogic

import (
	"daosuan/constants"
	authbase "daosuan/core/auth"
	"daosuan/core/cache"
	"daosuan/entity"
	"daosuan/exceptions/product"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"encoding/json"
	"github.com/kataras/iris"
)

var fields = []string{
	"Name", "AuthorId", "Id", "Description", "Cover", "Details", "Additional",
	"Status", "Tag", "CreateTime", "UpdateTime", "Star", "MasterVersion",
}

type ProductLogic struct {
	auth    authbase.DaoSuanAuthAuthorization
	product db.Product
}

func NewProductLogic(auth authbase.DaoSuanAuthAuthorization, pid ...int) ProductLogic {
	var product db.Product
	if len(pid) > 0 {
		table := db.Driver.Preload("Author")
		if err := db.Driver.GetOne("product", pid[0], &product, table); err != nil || product.Id == 0 {
			panic(productException.ProductIsNotExists())
		}
		db.Driver.Model(&product).Related(&(product.Tag), "Tag")
	}
	return ProductLogic{
		auth:    auth,
		product: product,
	}
}

func (p *ProductLogic) ProductModel() *db.Product {
	return &p.product
}

func (p *ProductLogic) SetProductModel(product db.Product) {
	p.product = product
}

func (p *ProductLogic) LoadVersions() {
	db.Driver.Model(&p.product).Related(&p.product.Versions)
}

func (p *ProductLogic) GetProductInfo() interface{} {
	data := paramsUtils.ModelToDict(p.product, fields)
	if len(p.product.Cover) > 0 {
		p.product.Cover = resourceLogic.GenerateToken(p.product.Cover, -1, -1)
	}
	if len(p.product.Additional) > 0 {
		var additional entity.ProductAdditional
		if err := json.Unmarshal([]byte(p.product.Additional), &additional); err == nil {
			data["additional"] = additional
		}
	}
	data["author_name"] = p.product.Author.Nickname
	if p.auth.IsLogin() && p.IsStar() {
		data["is_star"] = true
	} else {
		data["is_star"] = false
	}
	return data
}

func (p *ProductLogic) VersionIsExists(versionName string) bool {
	var t db.ProductVersion
	if err := db.Driver.Where("product_id = ? and version_name = ?", p.product.Id, versionName).First(&t).Error; err != nil || t.Id == 0 {
		return false
	}
	return true
}

func (p *ProductLogic) CheckSelf() {
	if p.product.AuthorId != p.auth.AccountModel().Id {
		panic(productException.NoPermission())
	}
}

func (p *ProductLogic) IsStar() bool {

	if !p.auth.IsLogin() {
		return false
	}

	if s, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.StarModel, p.product.Id, p.auth.AccountModel().Id)); err == nil && len(s) > 0 {
		return true
	}

	var star db.AccountStar
	if err := db.Driver.Where("product_id = ? and account_id = ?", p.product.Id, p.auth.AccountModel().Id).First(&star).Error; err != nil || star.Id == 0 {
		return false
	}
	v := hash.RandInt64(240, 240 * 5)
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.StarModel, p.product.Id, p.auth.AccountModel().Id), []byte("star"), int(v) * 60 * 60)
	return true
}

// 获取产品审核信息列表
func (p *ProductLogic) GetExamineRecord() []db.ProductExamineRecord {
	var result []db.ProductExamineRecord
	if r, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.ProductExamineRecordModel, p.product.Id)); err == nil {
		if err := json.Unmarshal(r, &result); err == nil {
			return result
		}
	}

	db.Driver.Where("product_id = ?", p.product.Id).Order("-create_time").Find(&result)
	if r, err := json.Marshal(&result); err == nil {
		v := hash.RandInt64(240, 240 * 5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.ProductExamineRecordModel, p.product.Id), r, int(v) * 60 * 60)
	}
	return result
}

// 获取产品版本信息列表
func (p *ProductLogic) GetVersionInfo() []iris.Map {

	var result []iris.Map
	if r, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.ProductVersionInfoModel, p.product.Id)); err == nil && r != nil {
		if err := json.Unmarshal(r, &result); err == nil {
			return result
		}
	}
	p.LoadVersions()
	product := p.product

	result = make([]iris.Map, len(product.Versions))
	for i := 0; i < len(product.Versions); i++ {
		result[i] = iris.Map {
			"name": product.Versions[i].VersionName,
			"id": product.Versions[i].Id,
		}
	}
	if r, err := json.Marshal(&result); err == nil {
		v := hash.RandInt64(240, 240 * 5)
		cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.ProductVersionInfoModel, p.product.Id), r, int(v) * 60 * 60)
	}
	return result
}

// 检测产品名是否存在
func IskNameExists(name string) bool {
	if name, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.ProductNameModel, name)); err == nil && name != nil {
		return true
	}

	var t db.Product
	if err := db.Driver.Where("name = ?", name).First(&t).Error; err != nil || t.Id == 0 {
		return false
	}
	v := hash.RandInt64(240, 240*5)
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.ProductNameModel, name), []byte(name), int(v)*60*60)
	return true
}