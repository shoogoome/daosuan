package accountLogic

import (
	"daosuan/constants"
	authbase "daosuan/core/auth"
	"daosuan/exceptions/account"
	"daosuan/logics/resource"
	"daosuan/models/db"
	paramsUtils "daosuan/utils/params"
)

var field = []string{
	"Nickname", "Email", "Id", "Role", "Phone", "PhoneValidated", "UpdateTime",
	"EmailValidated", "Avator", "Motto", "CreateTime", "Realname",
}

type AccountLogic interface {
	GetAccountInfo() interface{}
	AccountModel() db.Account
	SetAccountModel(account db.Account)
}

type accountStruct struct {
	auth    authbase.DaoSuanAuthAuthorization
	account db.Account
}

func NewAccountLogic(auth authbase.DaoSuanAuthAuthorization, aid ...int) AccountLogic {
	var account db.Account

	if len(aid) > 0 {
		if err := db.Driver.GetOne("account", aid[0], &account); err != nil || account.Id == 0 {
			panic(accountException.AccountIsNotExists())
		}
	} else {
		account = *auth.AccountModel()
	}
	return &accountStruct{
		account: account,
		auth:    auth,
	}
}

func (a *accountStruct) SetAccountModel(account db.Account) {
	a.account = account
}

func (a *accountStruct) AccountModel() db.Account {
	return a.account
}

func (a *accountStruct) GetAccountInfo() interface{} {
	// 卡权限
	if !a.auth.IsAdmin() && a.auth.AccountModel().Id != a.account.Id {
		panic(accountException.NoPermission())
	}

	if len(a.account.Avator) > 0 {
		a.account.Avator = resourceLogic.GenerateToken(a.account.Avator, -1, constants.DaoSuanSessionExpires)
	}

	return paramsUtils.ModelToDict(a.account, field)
}
