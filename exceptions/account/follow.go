package accountException

import "daosuan/models"

func FollowingFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5307,
		Message: "关注失败",
	}
}

func CancelFollowingFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5308,
		Message: "取消关注失败",
	}
}

