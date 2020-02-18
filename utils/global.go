package utils

import (
	"daosuan/models"
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
	GlobalConfig.Oauth.GitHub.Oauth2Config = oauth2.Config{
		ClientID:     GlobalConfig.Oauth.GitHub.ClientId,
		ClientSecret: GlobalConfig.Oauth.GitHub.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		RedirectURL: GlobalConfig.Oauth.GitHub.RedirectUrl,
		Scopes:      []string{"repo", "read:org", "read:user"},
	}

}