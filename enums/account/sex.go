package accountEnums

import enumsbase "daosuan/enums"

const (
	SexUnknown = 1 // 未知
	SexMale    = 2 // 男性
	SexFemale  = 4 // 女性
)

type SexEnums struct {
	enumsbase.EnumBase
}

func NewSexEnums() enumsbase.EnumBaseInterface {
	return SexEnums{
		EnumBase: enumsbase.EnumBase{
			Enums: []int{1, 2, 4},
		},
	}
}
