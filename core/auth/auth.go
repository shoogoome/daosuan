package authbase

import (
	"daosuan/constants"
	"daosuan/enums/account"
	accountException "daosuan/exceptions/account"
	"daosuan/models/db"
	"daosuan/utils/hash"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"time"
)

type DaoSuanAuthAuthorization interface {
	loadAuthenticationInfo()
	loadFromSession() bool
	loadFromCookie() bool
	SetSession(aid int)
	SetCookie(aid int)
	fetchAccount(aid int) bool
	CheckLogin()
	IsAdmin() bool
	CheckAdmin()
	AccountModel() *db.Account
	IsLogin() bool
}

var sess = sessions.New(sessions.Config{
	Cookie: constants.DaoSuanSystemCookie,
})

type cookieInfo struct {
	AccountId  int   `json:"account_id"`
	ExpireTime int64 `json:"expire_time"`
}

type daosuanAuthAuthorization struct {
	Account db.Account
	isLogin bool
	Context iris.Context
}

func NewAuthAuthorization(ctx *iris.Context) DaoSuanAuthAuthorization {
	authorization := daosuanAuthAuthorization{
		isLogin: false,
		Context: *ctx,
	}
	authorization.loadAuthenticationInfo()
	return &authorization
}

func (r *daosuanAuthAuthorization) AccountModel() *db.Account {
	return &r.Account
}

func (r *daosuanAuthAuthorization) CheckLogin() {
	if !r.isLogin {
		panic(accountException.AuthIsNotLogin())
	}
}

func (r *daosuanAuthAuthorization) IsLogin() bool {
	return r.isLogin
}

func (r *daosuanAuthAuthorization) IsAdmin() bool {
	return r.Account.Role == accountEnums.RoleAdmin
}

func (r *daosuanAuthAuthorization) CheckAdmin() {
	r.CheckLogin()
	if !r.IsAdmin() {
		panic(accountException.NoPermission())
	}
}

func (r *daosuanAuthAuthorization) loadAuthenticationInfo() {
	succ := r.loadFromSession()
	if !succ {
		r.loadFromCookie()
	}
}

// 从session载入登录信息
func (r *daosuanAuthAuthorization) loadFromSession() bool {
	// 阻止解码方法异常传递
	defer func() {
		recover()
	}()

	session := sess.Start(r.Context)
	sestring := session.GetString(constants.DaoSuanSessionName)
	if sestring == "" {
		return false
	}
	var sessionStruct cookieInfo
	hash.DecodeToken(sestring, &sessionStruct)

	if sessionStruct.ExpireTime <= time.Now().Unix() {
		return false
	}
	succ := r.fetchAccount(sessionStruct.AccountId)
	if succ {
		r.isLogin = true
	}
	return true
}

// 从cookie载入登录信息
func (r *daosuanAuthAuthorization) loadFromCookie() bool {
	defer func() {
		recover()
	}()
	cookie := r.Context.GetCookie(constants.DaoSuanSessionName)
	if cookie == "" {
		return false
	}
	var cookieStruct cookieInfo
	hash.DecodeToken(cookie, &cookieStruct)
	if cookieStruct.ExpireTime <= time.Now().Unix() {
		return false
	}
	succ := r.fetchAccount(cookieStruct.AccountId)
	if succ {
		r.isLogin = true
	}
	return true
}

// 设置session
func (r *daosuanAuthAuthorization) SetSession(aid int) {
	session := sess.Start(r.Context)

	if aid == 0 {
		session.Set(constants.DaoSuanSessionName, "")
		return
	}
	payload := generateToken(aid, constants.DaoSuanSessionExpires)
	session.Set(constants.DaoSuanSessionName, payload)

}

// 设置cookie
func (r *daosuanAuthAuthorization) SetCookie(aid int) {
	if aid == 0 {
		r.Context.SetCookieKV(constants.DaoSuanSessionName, "")
		return
	}
	payload := generateToken(aid, constants.DaoSuanCookieExpires)
	r.Context.SetCookieKV(constants.DaoSuanSessionName, payload)
}

// 从数据库查找该用户
func (r *daosuanAuthAuthorization) fetchAccount(aid int) bool {
	err := db.Driver.GetOne("account", aid, &r.Account)
	if err != nil || r.Account.Id == 0 {
		return false
	}
	return true
}

// 生成token
func generateToken(aid int, expire int64) string {
	payload := cookieInfo{
		AccountId:  aid,
		ExpireTime: expire + time.Now().Unix(),
	}
	return hash.GenerateToken(payload, true)
}
