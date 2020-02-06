package constants

// 缓存key
const (
	// 数据库缓存键
	DbModel = "model"

	// 昵称缓存键
	NicknameModel = "account:nickname"

	// 关注缓存键
	FollowingModel = "account:following"

	FollowerModel = "account:follower"

	// 用户点赞缓存键
	AccountProductModel = "account:product"

	// 点赞缓存键
	StarModel = "product:star"

	// 产品名缓存键
	ProductNameModel = "product:name"

	// 产品版本信息缓存键
	ProductVersionInfoModel = "product:version_info"

	// 产品审核信息缓存键
	ProductExamineRecordModel = "product:examine_record"

	// 产品issue回复缓存键
	ProductIssueReplyModel = "product:issue:reply"

)
