package account

import (
	accountEnums "daosuan/enums/account"
	accountException "daosuan/exceptions/account"
	accountLogic "daosuan/logics/account"
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

	email := params.Str("email", "用户名")
	if len(email) < 6 || len(email) > 50 {
		panic(accountException.UsernameLengthIsNotStandard())
	}
	db.Driver.Where("email = ?", email).First(&account)
	if account.Id != 0 {
		panic(accountException.EmailIsExists())
	}

	password := params.Str("password", "密码")
	if len(password) < 8 && len(password) > 32 {
		panic(accountException.PasswordLengthIsNotStandard())
	}
	nickname := params.Str("nickname", "昵称")
	if accountLogic.IsNicknameExists(nickname) {
		panic(accountException.NicknameIsExists())
	}
	account = db.Account{
		Email: email,
		Password: hash.PasswordSignature(password),
		Role: int16(accountEnums.RoleUser),
		EmailValidated: true, // 后续改
		Nickname: nickname,
	}
	db.Driver.Create(&account)
	ctx.JSON(iris.Map{
		"id": account.Id,
	})
}

