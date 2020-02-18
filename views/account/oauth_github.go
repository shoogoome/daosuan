package account

import (
	"context"
	"daosuan/core/auth"
	"daosuan/enums/account"
	resourceLogic "daosuan/logics/resource"
	"daosuan/models/db"
	"daosuan/utils"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 获取验证路由
func GitHubGetAuthUrl(ctx iris.Context) {
	referer := ctx.URLParam("referer")
	_type := ctx.URLParamIntDefault("type", accountEnums.GitHubLogging)
	if len(referer) > 0 {
		referer = url.QueryEscape(fmt.Sprintf("%s:%d", referer, _type))
	} else {
		referer = url.QueryEscape(fmt.Sprintf("%s:%d", utils.GlobalConfig.Oauth.GitHub.SuccessUrl, _type))
	}
	ctx.JSON(iris.Map{
		"url": utils.GlobalConfig.Oauth.GitHub.Oauth2Config.AuthCodeURL(referer),
	})
}

// github验证回调路由
func GitHubCallback(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	state := ctx.URLParam("state")
	code := ctx.URLParam("code")

	// 获取token
	token, err := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Exchange(context.Background(), code)

	if err != nil {
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}
	// 验证
	oauth2Client := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Client(context.Background(), token)
	client := github.NewClient(oauth2Client)
	userInfo, _, err := client.Users.Get(context.Background(), "")

	if err != nil || userInfo == nil {
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}

	stateSplit := strings.Split(state, ":")
	if len(stateSplit)  != 2 {
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}
	userinfo, _ := json.Marshal(userInfo)
	// 登录
	var accountOauth db.AccountOauth
	if err := db.Driver.
		Where("model = ? and open_id = ?", accountEnums.OauthGitHub, userInfo.ID).
		First(&accountOauth).Error; err != nil || accountOauth.Id == 0 {
		// 如果存在这个账号并且想绑定他则直接抛异常（提示去账号合并）
		if stateSplit[1] == strconv.Itoa(accountEnums.GitHubBinding){
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
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
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		// 尝试获取头像信息 (但github现阶段墙了头像)
		if !getAvator(tx, *userInfo.AvatarURL, &account) {
			tx.Callback()
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		// 绑定关联
		aid := createOauth(tx, account.Id, int(*userInfo.ID), string(userinfo))
		if aid == 0 {
			tx.Callback()
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		tx.Commit()
		accountOauth.AccountId = aid
	// 找不到这个账户并且想绑定则直接绑定
	} else if stateSplit[1] == strconv.Itoa(accountEnums.GitHubBinding) {
		auth.CheckLogin()
		aid := createOauth(db.Driver.DB, auth.AccountModel().Id, int(*userInfo.ID), string(userinfo))
		if aid == 0 {
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		ctx.Redirect(stateSplit[0], http.StatusFound)
		return
	}
	// 不管是第几次都直接给登录态
	auth.SetSession(accountOauth.AccountId)
	auth.SetCookie(accountOauth.AccountId)
	if len(state) > 0 {
		ctx.Redirect(state, http.StatusFound)
	} else {
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.SuccessUrl, http.StatusFound)
	}
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