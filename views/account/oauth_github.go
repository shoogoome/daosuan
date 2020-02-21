package account

import (
	"context"
	"daosuan/constants"
	"daosuan/core/auth"
	"daosuan/enums/account"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils"
	"daosuan/utils/hash"
	"daosuan/utils/log"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"strconv"
)

type stateJWT struct {
	Referer string `json:"referer"`
	Type int `json:"type"`
	AccountId int `json:"account_id"`
}


// 获取验证路由
func GitHubGetAuthUrl(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {
	referer := ctx.URLParam("referer")
	_type := ctx.URLParamIntDefault("type", accountEnums.GitHubLogging)

	token := stateJWT {
		Type: _type,
	}
	if len(referer) > 0 {
		token.Referer = referer
	} else {
		token.Referer = utils.GlobalConfig.Oauth.GitHub.SuccessUrl
	}
	if auth.IsLogin() {
		token.AccountId = auth.AccountModel().Id
	}
	ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.Oauth2Config.AuthCodeURL(hash.GenerateToken(token, true)), http.StatusFound)
}

// github验证回调路由
func GitHubCallback(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
	}()
	state := ctx.URLParam("state")
	code := ctx.URLParam("code")

	// 获取token
	token, err := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Exchange(context.Background(), code)

	if err != nil {
		logUtils.Println("错误3")
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}
	// 验证
	oauth2Client := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Client(context.Background(), token)
	client := github.NewClient(oauth2Client)
	userInfo, _, err := client.Users.Get(context.Background(), "")

	if err != nil || userInfo == nil {
		logUtils.Println("错误2")
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}

	var jwt stateJWT
	hash.DecodeToken(state, &jwt)

	userinfo, _ := json.Marshal(userInfo)
	// 登录
	var accountOauth db.AccountOauth
	if err := db.Driver.
		Where("model = ? and open_id = ?", accountEnums.OauthGitHub, strconv.Itoa(int(*userInfo.ID))).
		First(&accountOauth).Error; err != nil {
		// 找不到这个账户并且想绑定则直接绑定
		if jwt.Type == accountEnums.GitHubBinding {
			if jwt.AccountId == 0 {
				logUtils.Println("没有用户id")
				ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
				return
			}
			aid := createOauth(db.Driver.DB, jwt.AccountId, int(*userInfo.ID), string(userinfo))
			if aid == 0 {
				logUtils.Println("错误8")
				ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
				return
			}
			logUtils.Println("错误9")
			ctx.Redirect(jwt.Referer, http.StatusFound)
			return

		}

		tx := db.Driver.Begin()
		// 创建用户
		account := db.Account{
			Nickname: *userInfo.Login,
			Role: accountEnums.RoleUser,
			Init: true,
		}

		if err := tx.Create(&account).Error; err != nil {
			tx.Callback()
			logUtils.Println("错误5")
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		// 尝试获取头像信息 (但github现阶段墙了头像)
		if !getAvator(tx, *userInfo.AvatarURL, &account) {
			tx.Callback()
			logUtils.Println("错误6")
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		// 绑定关联
		aid := createOauth(tx, account.Id, int(*userInfo.ID), string(userinfo))
		if aid == 0 {
			tx.Callback()
			logUtils.Println("错误7")
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		tx.Commit()
		accountOauth.AccountId = aid
	// 找到了这个账户并且他是想要绑定的话
	} else if jwt.Type == accountEnums.GitHubBinding {
		// 如果存在这个账号并且想绑定他则直接抛异常（提示去账号合并）
		logUtils.Println("错误4")
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}

	cookie := http.Cookie{
		Name: constants.DaoSuanSessionName,
		Value: authbase.GenerateToken(accountOauth.AccountId, constants.DaoSuanCookieExpires),
		Domain: utils.GlobalConfig.Oauth.GitHub.CookieDomain,
		Path: "/",
		MaxAge: constants.DaoSuanCookieExpires,
	}

	ctx.SetCookie(&cookie)
	logUtils.Println("奥利给")
	ctx.Redirect(jwt.Referer, http.StatusFound)
}

// 绑定
func createOauth (tx *gorm.DB, aid, openid int, userinfo string) int {
	accountOauth := db.AccountOauth{
		AccountId: aid,
		Model: accountEnums.OauthGitHub,
		OpenId: strconv.Itoa(openid),
		UserInfo: userinfo,
	}
	if err := tx.Create(&accountOauth).Error; err != nil {
		return 0
	}
	return accountOauth.Id
}

// 获取头像数据
func getAvator(tx *gorm.DB, url string, account *db.Account) bool {
	if response, err := utils.Requests("GET", url, nil); err == nil && response.StatusCode == http.StatusOK {
		if body, err := ioutil.ReadAll(response.Body); err == nil {
			defer response.Body.Close()
			logic := resourceLogic.NewReousrcesLocalStorage("account_avator")
			account.Avator = logic.SaveFile(fmt.Sprintf("%d/%s", account.Id, "avator.jpg"), body, true)
		}
		if err := tx.Save(&account).Error; err != nil {
			return false
		}
	}
	return true
}