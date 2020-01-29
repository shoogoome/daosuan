package productException

import "daosuan/models"

func ProductVersionNameIsExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5905,
		Message: "版本名称已存在",
	}
}

func ProductVersionIsNotExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5906,
		Message: "版本不存在",
	}
}

