package productException

import "daosuan/models"

func ProductStatusIsNotExamine() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5910,
		Message: "产品未提交申请请求",
	}
}

