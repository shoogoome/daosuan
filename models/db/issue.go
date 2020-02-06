package db

// 问题表
type Issue struct {

	Id int `json:"id"`

	// 产品id
	ProductId int `json:"product_id" gorm:"not null"`

	// 发布人id
	AuthorId int `json:"author_id" gorm:"not null"`

	// 内容
	Content string `json:"content" gorm:"not null;type:text"`

	// 回复数量
	ReplyNumber int `json:"reply_number" gorm:"default:0"`

	// 创建时间
	CreateTime int64 `json:"create_time"`

	// 更新时间
	UpdateTime int64 `json:"update_time"`
}

// 问题回复
type IssueReply struct {

	Id int `json:"id"`

	// 回复者id
	AuthorId int `json:"author_id" gorm:"not null"`

	// 关联问题id
	IssueId int `json:"issue_id" gorm:"not null"`

	// @回复人id
	AtAccountId int `json:"at_account_id"`

	// 回复内容
	Reply string `json:"reply" grom:"type:text;not null"`

	// 创建时间
	CreateTime int64 `json:"create_time"`
}


