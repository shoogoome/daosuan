package db

import (
	_ "github.com/go-sql-driver/mysql"
)

// 产品表
type Product struct {
	Id int `gorm:"primary_key" json:"id"`

	// 名称
	Name string `json:"name" gorm:"not null"`

	// 用户id
	Author   Account `json:"author" gorm:"ForeignKey:AuthorId"`
	AuthorId int     `json:"author_id"`

	// 描述
	Description string `json:"description"`

	// 封面
	Cover string `json:"cover"`

	// 详情页(冗余字段)
	Details string `json:"details" gorm:"not null;type:longtext"`

	// 附加页(冗余字段)
	Additional string `json:"additional" gorm:"default:'{}';type:longtext"`

	// 状态enum
	Status int16 `json:"status" gorm:"not null"`

	// 赞
	Star int64 `json:"star"`

	// 课程类型
	Tag []Tag `json:"tag" gorm:"many2many:product_tags"`

	// 主版本名
	MasterVersion string `json:"master_version"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`

	// --------------------------------------------------

	// 版本关联
	Versions []ProductVersion `json:"versions" gorm:"ForeignKey:ProductId"`
}

// 产品版本表
type ProductVersion struct {
	Id int `gorm:"primary_key" json:"id"`

	// 产品关联
	ProductId int `json:"-"`

	// 版本名
	VersionName string `json:"version_name"`

	// 详情页(冗余字段)
	Details string `json:"details" gorm:"not null;type:longtext"`

	// 附加页(冗余字段)
	Additional string `json:"additional" gorm:"default:'{}';type:longtext"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

// 产品审核记录表
type ProductExamineRecord struct {
	Id int `gorm:"primary_key" json:"id"`

	// 产品关联
	ProductId int `json:"-"`

	// 通过与否
	Adopt bool `json:"adopt"`

	// 回复
	Reply string `json:"reply"`

	// 创建时间
	CreateTime int64 `json:"create_time"`
}
