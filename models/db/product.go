package db

import (
	_ "github.com/go-sql-driver/mysql"
)

// 产品表
type Product struct {

	Id        int `gorm:"primary_key" json:"id"`

	// 名称
	Name string `json:"name" gorm:"not null"`

	// 用户id
	AuthorId int `json:"author_id"`

	// 描述
	Description string `json:"description"`

	// 封面
	Cover string `json:"cover"`

	// 详情页
	Details string `json:"details" gorm:"not null;type:longtext"`

	// 附加页
	Additional string `json:"additional" gorm:"default:'{}';type:longtext"`

	// 状态enum
	Status int16 `json:"status" gorm:"not null"`

	// 课程类型
	Tag []Tag `json:"tag" gorm:"many2many:product_tags"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}
