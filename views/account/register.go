package account

import (
	accountEnums "daosuan/enums/account"
	accountException "daosuan/exceptions/account"
	"daosuan/models/db"
	"daosuan/utils/hash"
	paramsUtils "daosuan/utils/params"
	"github.com/kataras/iris"
)

// 创建账户
func Register(ctx iris.Context) {

	data := paramsUtils.RequestJsonInterface(ctx)
	params := paramsUtils.NewParamsParser(data)

	var account db.Account
	username := params.Str("username", "用户名")
	if len(username) < 6 || len(username) > 50 {
		panic(accountException.UsernameLengthIsNotStandard())
	}
	db.Driver.Where("username = ?", username).First(&account)
	if account.Id != 0 {
		panic(accountException.NicknameIsExists())
	}

	password := params.Str("password", "密码")
	if len(password) < 8 && len(password) > 32 {
		panic(accountException.PasswordLengthIsNotStandard())
	}

	account = db.Account{
		Email: username,
		Password: hash.PasswordSignature(password),
		Nickname: username,
		Role: int16(accountEnums.RoleUser),
	}
	db.Driver.Create(&account)
	ctx.JSON(iris.Map{
		"id": account.Id,
	})
}

