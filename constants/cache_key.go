package constants

// 缓存key
const (
	DbModel                      = "model"                           // 数据库缓存键
	NicknameModel                = "account:nickname"                // 昵称缓存键
	FollowingModel               = "account:following"               // 关注缓存键
	FollowerModel                = "account:follower"                // 被关注缓存键
	AccountProductModel          = "account:product"                 // 用户产品缓存键
	AccountVerificationEmail     = "account:verification:email"      // 用户发送邮件token缓存键
	AccountVerificationEmailTime = "account:verification:email:time" // 用户发送邮件延迟缓存键
	StarModel                    = "product:star"                    // 点赞缓存键
	ProductNameModel             = "product:name"                    // 产品名缓存键
	ProductVersionInfoModel      = "product:version_info"            // 产品版本信息缓存键
	ProductExamineRecordModel    = "product:examine_record"          // 产品审核信息缓存键
	ProductIssueReplyModel       = "product:issue:reply"             // 产品issue回复缓存键

)
