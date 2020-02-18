package account

import (
	"context"
	"daosuan/exceptions/account"
	"daosuan/utils"
	"daosuan/utils/log"
	"encoding/json"
	"github.com/google/go-github/github"
	"github.com/kataras/iris"
	"net/url"
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
func GitHubCallback(ctx iris.Context) {

	state := ctx.URLParam("state")
	code := ctx.URLParam("code")

	token, err := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Exchange(context.Background(), code)

	if err != nil {
		panic(accountException.OauthVerificationFail())
	}
	oauth2Client := utils.GlobalConfig.Oauth.GitHub.Oauth2Config.Client(context.Background(), token)
	client := github.NewClient(oauth2Client)
	userInfo, _, err := client.Users.Get(context.Background(), "")

	if err != nil || userInfo == nil {
		panic(accountException.OauthVerificationFail())
	}

	logUtils.Println(json.Marshal(userInfo))
	logUtils.Println(state)

}
