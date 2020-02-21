package account

import (
	"daosuan/constants"
	"daosuan/core/cache"
	accountEnums "daosuan/enums/account"
	accountException "daosuan/exceptions/account"
	accountLogic "daosuan/logics/account"
	"daosuan/models/db"
	"daosuan/utils/hash"
	"daosuan/utils/mail"
	paramsUtils "daosuan/utils/params"
	"github.com/kataras/iris"
)

// 创建账户
func EmailRegister(ctx iris.Context) {

	data := paramsUtils.RequestJsonInterface(ctx)
	params := paramsUtils.NewParamsParser(data)
	vToken :=  params.Str("token", "验证码")
	password := params.Str("password", "密码")
	nickname := params.Str("nickname", "昵称")
	email := params.Str("email", "邮箱")

	// 密码长度检测（先检查这个可以不用多余sql操作）
	if len(password) < 8 && len(password) > 32 {
		panic(accountException.PasswordLengthIsNotStandard())
	}

	var account db.Account
	// 检测邮箱格式
	if !mailUtils.CheckMailFormat(email) {
		panic(accountException.EmailFormatError())
	}
	// 检测重复
	db.Driver.Where("email = ?", email).First(&account)
	if account.Id != 0 {
		panic(accountException.EmailIsExists())
	}

	// token验证
	if token, err := cache.Dijan.Get(paramsUtils.CacheBuildKey(constants.AccountVerificationEmail, email)); err == nil && token != nil {
		if string(token) != vToken {
			panic(accountException.OauthVerificationFail())
		}
	}
	// 昵称检查
	if accountLogic.IsNicknameExists(nickname) {
		panic(accountException.NicknameIsExists())
	}
	account = db.Account{
		Email: email,
		Password: hash.PasswordSignature(password),
		Role: int16(accountEnums.RoleUser),
		EmailValidated: true,
		Nickname: nickname,
	}
	db.Driver.Create(&account)
	ctx.JSON(iris.Map{
		"id": account.Id,
	})
}

