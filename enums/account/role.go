package accountEnums

import "daosuan/enums"

const (
	RoleUser  = 1  // 普通用户
	RoleAdmin = 99 // 系统管理员
)

type RoleEnums struct {
	enumsbase.EnumBase
}

func NewRoleEnums() enumsbase.EnumBaseInterface {
	return RoleEnums{
		EnumBase: enumsbase.EnumBase{
			Enums: []int{1, 99},
		},
	}
}
