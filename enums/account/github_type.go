package accountEnums

import "daosuan/enums"

const (
	GitHubBinding  = 1    // github绑定
	GitHubLogging  = 2    // github登录
)

func NewGitHubTypeEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase {
		Enums: []int{1, 2},
	}
}

