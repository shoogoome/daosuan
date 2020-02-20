package account

import (
	"daosuan/core/auth"
	"daosuan/utils"
	"daosuan/utils/log"
	"github.com/kataras/iris"
)

func WeChatGetAuthUrl(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	ctx.JSON(utils.GlobalConfig.Oauth.WeChat.OauthClient.AuthCodeUrl("123"))

}


func WeChatCallback(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization) {

	code := ctx.URLParam("code")
	state := ctx.URLParam("state")

	logUtils.Println(code, state)

	token, err := utils.GlobalConfig.Oauth.WeChat.OauthClient.Exchange(code)
	logUtils.Println(token, err)

	userInfo, err := utils.GlobalConfig.Oauth.WeChat.OauthClient.GetUserInfo(token.AccessToken, token.OpenId)

	logUtils.Println(userInfo, err)

	ctx.JSON(iris.Map {
		"token": token,
		"userinfo": userInfo,
	})

}