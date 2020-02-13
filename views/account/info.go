package account

import (
	authbase "daosuan/core/auth"
	"daosuan/core/queue"
	"daosuan/enums/account"
	accountException "daosuan/exceptions/account"
	accountLogic "daosuan/logics/account"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils/durl"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"fmt"
	"github.com/kataras/iris"
)

// 获取用户信息
func GetAccount(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {

	if err := queue.Message.Publish("test", []byte("nihao")); err != nil {
		panic(err)
	}

	//logic := accountLogic.NewAccountLogic(auth, aid)
	//ctx.JSON(logic.GetAccountInfo())
}

// 获取用户列表
func GetAccountList(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	auth.CheckAdmin()

	var lists []struct{
		Id int `json:"id"`
		UpdateTime int64 `json:"update_time"`
		Role int16 `json:"role"`
		Nickname string `json:"nickname"`
	}
	var count int
	table := db.Driver.Table("account")

	limit := ctx.URLParamIntDefault("limit", 10)
	page := ctx.URLParamIntDefault("page", 1)

	// 条件过滤
	if key := ctx.URLParam("key");len(key) > 0 {
		keyString := fmt.Sprintf("%%%s%%", key)
		table = table.Where("nickname like ? or email like ?", keyString, keyString)
	}

	table.Count(&count).Offset((page - 1) * limit).Limit(limit).Select("id, nickname, update_time, role").Find(&lists)
	ctx.JSON(iris.Map {
		"accounts": lists,
		"total": count,
		"limit": limit,
		"page": page,
	})
}

// 批量获取信息
func MgetAccountInfo(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := accountLogic.NewAccountLogic(auth)

	ids := params.List("ids", "id列表")
	var data []interface{}
	accounts := db.Driver.GetMany("account", ids, db.Account{})
	for _, account := range accounts {
		logic.SetAccountModel(account.(db.Account))
		func(data *[]interface{}) {
			*data = append(*data, logic.GetAccountInfo())
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


// 修改用户信息
func PutAccount(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, aid int) {
	auth.CheckLogin()
	params := paramsUtils.NewParamsParser(paramsUtils.RequestJsonInterface(ctx))
	logic := accountLogic.NewAccountLogic(auth, aid)
	account := logic.AccountModel()

	if !auth.IsAdmin() && account.Id != auth.AccountModel().Id {
		panic(accountException.NoPermission())
	}

	params.Diff(*account)
	account.Nickname = params.Str("nickname", "昵称")
	// 检测昵称存在情况
	var a db.Account
	if err := db.Driver.Where("nickname = ? and id != ?", account.Nickname, account.Id).First(&a).Error; err == nil && account.Id != 0 {
		panic(accountException.NicknameExists())
	}
	account.Motto = params.Str("motto", "一句话签名")

	if params.Has("role") && auth.IsAdmin() {
		accountEnum := accountEnums.NewRoleEnums()
		if accountEnum.Has(params.Int("role", "角色")) {
			account.Role = int16(params.Int("role", "角色"))
		}
	}

	if params.Has("new_password") {
		newPassword := params.Str("new_password", "旧密码")
		if !auth.IsAdmin() {
			oldPassword := params.Str("old_password", "旧密码")
			if hash.PasswordSignature(oldPassword) != account.Password {
				panic(accountException.OldPasswordIsNotTrue())
			}
		}
		account.Password = hash.PasswordSignature(newPassword)
	}

	if params.Has("avator") {
		logic := resourceLogic.NewReousrcesLocalStorage("account_avator")
		account.Avator = logic.SaveFile(fmt.Sprintf("%d/%s", account.Id, "avator.jpg"), durl.DataUrlParser(params.Str("avator", "头像")), true)
	}

	db.Driver.Save(&account)
	ctx.JSON(iris.Map {
		"id": account.Id,
	})
}

// 检测昵称存在与否
func CheckNicknameExists(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, name string) {
	ctx.JSON(iris.Map {
		"exists": accountLogic.IsNicknameExists(name),
	})
}


