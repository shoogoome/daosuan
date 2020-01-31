package db

type Tag struct {

	// id
	Id int `json:"id"`

	// 名称
	Name string `json:"name" gorm:"not null"`

	// 创建时间
	CreateTime int64 `json:"-"`

	// 产品关联
	Products []Product `json:"-"`
}

