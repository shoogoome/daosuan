package authbase

import (
	"daosuan/models/db"
	"github.com/kataras/iris"
)



type daosuanClientAuthAuthorization struct {
	Account db.Account
	isLogin bool
	Context iris.Context
}

func NewClientAuthAuthorization(ctx *iris.Context) DaoSuanAuthAuthorization {
	authorization := daosuanClientAuthAuthorization{
		isLogin: false,
		Context: *ctx,
	}
	authorization.loadAuthenticationInfo()
	return &authorization
}


func (daosuanClientAuthAuthorization) loadAuthenticationInfo() {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) loadFromSession() bool {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) loadFromCookie() bool {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) SetSession(aid int) {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) SetCookie(aid int) {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) fetchAccount(aid int) bool {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) CheckLogin() {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) IsAdmin() bool {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) CheckAdmin() {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) AccountModel() *db.Account {
	panic("implement me")
}

func (daosuanClientAuthAuthorization) IsLogin() bool {
	panic("implement me")
}

