package productException

import "daosuan/models"

func SearchFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5919,
		Message: "检索失败",
	}
}

