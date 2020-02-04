package productException

import "daosuan/models"

func ProductStatusIsNotExamine() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5910,
		Message: "产品未提交申请请求",
	}
}

func ExamineFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5911,
		Message: "操作失败",
	}
}

