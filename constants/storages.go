package constants

const (

	DaoSuanDataRoot = "../data"

	// 账户
	DaoSuanAccount = DaoSuanDataRoot + "/account"

	// 产品
	DaoSuanProduct = DaoSuanDataRoot + "/product"

	// ----- 子模块 -----

	// 账户头像
	DaoSuanAccountAvator = DaoSuanAccount + "/avator"


	// 产品封面
	DaoSuanProductCover = DaoSuanProduct + "/cover"


	// nginx静态资源映射
	NginxResourcePath = "/resource_internal"
)

var StorageMapping = map[string]string {
	"account_avator": DaoSuanAccountAvator,
	"product_cover": DaoSuanProductCover,
}

var MimeToExtMapping = map[string]string {
	"jpg": "image/jpeg",
	"jpeg": "image/jpeg",
	"bmp": "image/bmp",
	"png": "image/png",
	"gif": "image/gif",
	"svg": "image/svg",
}
