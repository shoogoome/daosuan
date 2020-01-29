package db

import (
	_ "github.com/go-sql-driver/mysql"
)

// 用户主账户表
type Account struct {

	Id        int `gorm:"primary_key" json:"id"`

	// 昵称
	Nickname string `json:"nickname" gorm:"not null"`

	// 角色
	Role int16 `json:"role" gorm:"not null"`

	// 电话
	Phone string `json:"phone"`

	// 电话验证与否
	PhoneValidated bool `json:"phone_validated" gorm:"default:false"`

	// 密码
	Password string `json:"password" gorm:"not null"`

	// 邮箱
	Email string `json:"email"`

	// 邮箱验证与否
	EmailValidated bool `json:"email_validated" gorm:"default:false"`

	// 头像
	Avator string `json:"avator"`

	// 一句话签名
	Motto string `json:"motto"`

	// 设置 保留字段
	Options string `json:"options" gorm:"default:''"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`

	// ----------------------------------------
}

// 用户点赞关联表
type AccountStar struct {
	Id        int `gorm:"primary_key" json:"id"`

	AccountId int `json:"account_id" `
	Account Account `json:"-" gorm:"ForeignKey:AccountId"`

	ProductId int `json:"product_id"`
	Product Product `json:"-" gorm:"ForeignKey:ProductId"`

	CreateTime int64 `json:"create_time"`
}

// 用户关注关联表
type AccountFollow struct {

	Id        int `gorm:"primary_key" json:"id"`

	SourceId int `json:"source_id"`
	Source Account `json:"source" gorm:"ForeignKey:SourceId"`

	TargetId int `json:"target_id"`
	Target Account `json:"target" gorm:"ForeignKey:TargetId"`

	CreateTime int64 `json:"create_time"`
}
