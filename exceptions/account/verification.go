package accountException

import "daosuan/models"

func EmailSendFrequently() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5310,
		Message: "两次发送邮件之间相隔不小于60秒",
	}
}

func EmailSendFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5311,
		Message: "邮箱发送失败， 检查目标邮箱是否存在或稍后再试",
	}
}

func EmailFormatError() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5312,
		Message: "邮箱格式错误",
	}
}
