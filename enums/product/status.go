package productEnums

import "daosuan/enums"

const (
	StatusReleased   = 1 // 发布
	StatusDraft      = 2 // 草稿
	StatusLowerShelf = 4 // 下架
)

type StatusEnums struct {
	enumsbase.EnumBase
}

func NewStatusEnums() enumsbase.EnumBaseInterface {
	return StatusEnums{
		EnumBase: enumsbase.EnumBase{
			Enums: []int{1, 2, 4},
		},
	}
}
