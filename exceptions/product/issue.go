package productException

import "daosuan/models"

func CreateIssueFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5912,
		Message: "发起提问失败",
	}
}

func IssueIsNotExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5913,
		Message: "提问不存在",
	}
}

func AtAccountIsNotFound() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5914,
		Message: "目标回复人不存在",
	}
}

func ReplyIsNotExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5915,
		Message: "回复不存在",
	}
}

func ReplyFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5916,
		Message: "回复失败",
	}
}


