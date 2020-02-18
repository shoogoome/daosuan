package accountEnums

import "daosuan/enums"

const (
	OauthWeChat  = 1    // 微信
	OauthGitHub  = 2    // GitHub
)

func NewOauthEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase {
		Enums: []int{1, 2},
	}
}

