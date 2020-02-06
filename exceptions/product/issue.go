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


