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
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// 获取验证路由
func GitHubGetAuthUrl(ctx iris.Context) {
	referer := ctx.URLParam("referer")
	if len(referer) > 0 {
		referer = url.QueryEscape(referer)
	}
	ctx.JSON(iris.Map{
		"url": utils.GlobalConfig.Oauth.GitHub.Oauth2Config.AuthCodeURL(referer),
	})
}

// github验证回调路由
func GitHubCallback(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	state := ctx.URLParam("state")
	code := ctx.URLParam("code")

	token, err := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Exchange(context.Background(), code)

	if err != nil {
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}
	oauth2Client := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Client(context.Background(), token)
	client := github.NewClient(oauth2Client)
	userInfo, _, err := client.Users.Get(context.Background(), "")

	if err != nil || userInfo == nil {
		ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
		return
	}

	var accountOauth db.AccountOauth
	// 第一次登录
	if err := db.Driver.
		Where("model = ? and open_id = ?", accountEnums.OauthGitHub, userInfo.ID).
		First(&accountOauth).Error; err != nil || accountOauth.Id == 0 {

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
		if response, err := utils.Requests("GET", *userInfo.AvatarURL, nil); err == nil && response.StatusCode == http.StatusOK {
			if body, err := ioutil.ReadAll(response.Body); err == nil {
				defer response.Body.Close()
				logic := resourceLogic.NewReousrcesLocalStorage("account_avator")
				account.Avator = logic.SaveFile(fmt.Sprintf("%d/%s", account.Id, "avator.jpg"), body, true)
			}
			if err := tx.Save(&account).Error; err != nil {
				tx.Callback()
				ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
				return
			}
		}
		// 绑定关联
		userinfo, _ := json.Marshal(userInfo)
		accountOauth = db.AccountOauth{
			AccountId: account.Id,
			Model: accountEnums.OauthGitHub,
			OpenId: strconv.Itoa(int(*userInfo.ID)),
			UserInfo: string(userinfo),
		}
		if err := tx.Create(&accountOauth).Error; err != nil {
			tx.Callback()
			ctx.Redirect(utils.GlobalConfig.Oauth.GitHub.ErrorUrl, http.StatusFound)
			return
		}
		tx.Commit()
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
