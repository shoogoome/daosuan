package authbase

import (
	"daosuan/constants"
	"daosuan/enums/account"
	accountException "daosuan/exceptions/account"
	"daosuan/models/db"
	"daosuan/utils"
	"daosuan/utils/hash"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"net/http"
	"time"
)

type DaoSuanAuthAuthorization interface {
	SetSession(aid int)
	SetCookie(aid int)
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
		session.Delete(constants.DaoSuanSessionName)
		return
	}
	payload := GenerateToken(aid, constants.DaoSuanSessionExpires)
	session.Set(constants.DaoSuanSessionName, payload)

}

// 设置cookie
func (r *daosuanAuthAuthorization) SetCookie(aid int) {
	if aid == 0 {
		r.Context.RemoveCookie(constants.DaoSuanSessionName)
		// 删除oauth登录的cookit
		cookie := http.Cookie{
			Name: constants.DaoSuanSessionName,
			Domain: utils.GlobalConfig.Oauth.GitHub.CookieDomain,
			Path: "/",
			MaxAge: -1,
		}
		r.Context.SetCookie(&cookie)
		return
	}
	payload := GenerateToken(aid, constants.DaoSuanCookieExpires)
	r.Context.SetCookieKV(constants.DaoSuanSessionName, payload)
}

// 从数据库查找该用户
func (r *daosuanAuthAuthorization) fetchAccount(aid int) bool {
	err := db.Driver.GetOne("account", aid, &r.Account)
	if err != nil {
		return false
	}
	return true
}

// 生成token
func GenerateToken(aid int, expire int64) string {
	payload := cookieInfo{
		AccountId:  aid,
		ExpireTime: expire + time.Now().Unix(),
	}
	return hash.GenerateToken(payload, true)
}
