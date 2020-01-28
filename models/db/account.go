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

	// 真实姓名
	Realname string `json:"realname"`

	// 身份证
	IdCode string `json:"id_code"`

	// 设置 保留字段
	Options string `json:"options" gorm:"default:''"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`

	// ----------------------------------------
}
