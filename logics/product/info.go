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
	"fmt"
)

var fields = []string{
	"Name", "AuthorId", "Id", "Description", "Cover", "Details", "Additional",
	"Status", "Tag", "CreateTime", "UpdateTime", "Star", "MasterVersion",
}

type ProductLogic interface {
	GetProductInfo() interface{}
	ProductModel() *db.Product
	SetProductModel(product db.Product)
	LoadVersions()
	VersionIsExists(string) bool
	CheckSelf()
	IsStar() bool
}

type productStruct struct {
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
		fmt.Println(product.Author)
		db.Driver.Model(&product).Related(&(product.Tag), "Tag")
	}
	return &productStruct{
		auth:    auth,
		product: product,
	}
}

func (p *productStruct) ProductModel() *db.Product {
	return &p.product
}

func (p *productStruct) SetProductModel(product db.Product) {
	p.product = product
}

func (p *productStruct) LoadVersions() {
	db.Driver.Model(&p.product).Related(&p.product.Versions)
}

func (p *productStruct) GetProductInfo() interface{} {
	data := paramsUtils.ModelToDict(p.product, fields)
	if len(p.product.Cover) > 0 {
		p.product.Cover = resourceLogic.GenerateToken(p.product.Cover, -1, constants.DaoSuanSessionExpires)
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

func (p *productStruct) VersionIsExists(versionName string) bool {
	var t db.ProductVersion
	if err := db.Driver.Where("product_id = ? and version_name = ?", p.product.Id, versionName).First(&t).Error; err != nil || t.Id == 0 {
		return false
	}
	return true
}

func (p *productStruct) CheckSelf() {
	if p.product.AuthorId != p.auth.AccountModel().Id {
		panic(productException.NoPermission())
	}
}

func (p *productStruct) IsStar() bool {

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
	cache.Dijan.Set(paramsUtils.CacheBuildKey(constants.StarModel, p.product.Id, p.auth.AccountModel().Id), "star", int(v) * 60 * 60)
	return true
}