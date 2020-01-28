package tagException

import "daosuan/models"

func TagIsNotExsits() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5700,
		Message: "标签不存在",
	}
}



