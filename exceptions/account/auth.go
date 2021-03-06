package accountException

import (
	"daosuan/models"
)

func AuthIsNotLogin() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5300,
		Message: "尚未登录",
	}
}

func NoPermission() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "无权限执行此操作",
	}
}

func OauthVerificationFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5302,
		Message: "验证码错误",
	}
}