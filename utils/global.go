package utils

import (
	"daosuan/models"
	"daosuan/utils/wechat"
	"fmt"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var GlobalConfig models.SystemConfiguration

func InitGlobal() {
	yamlFile, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %s", err.Error()))
	}
	err = yaml.Unmarshal(yamlFile, &GlobalConfig)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal configuration: %s", err.Error()))
	}

	// 初始化github oauth配置
	GlobalConfig.Oauth.GitHub.Oauth2Config = oauth2.Config {
		ClientID:     GlobalConfig.Oauth.GitHub.ClientId,
		ClientSecret: GlobalConfig.Oauth.GitHub.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: GlobalConfig.Oauth.GitHub.RedirectUrl,
		Scopes:      []string{"repo", "read:org", "read:user"},
	}

	// 初始化wechat oauth配置
	GlobalConfig.Oauth.WeChat.OauthClient = wecharUtils.WeCharClient{
		Appid: "wxbbf5d0d2fe30d53d",
		Secret: "6e93d09a9eccb23ca520b9fe16a8ff6d",
		RedirectUri: "http://api.v1.daosuan.net/accounts/oauth/wechat/callback",
		Scope: []string{"snsapi_login"},
	}
}