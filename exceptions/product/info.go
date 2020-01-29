package productException

import "daosuan/models"

func ProductIsNotExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5901,
		Message: "产品不存在",
	}
}

func ProductCreateFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5902,
		Message: "产品创建失败",
	}
}

func NoPermission() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5903,
		Message: "无权限执行此操作",
	}
}

func NameIsExist() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5904,
		Message: "名称已存在",
	}
}

func StarFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5907,
		Message: "Star失败",
	}
}

func ReStar() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5908,
		Message: "不得重复Star",
	}
}

func CancelStarFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5909,
		Message: "取消Star失败",
	}
}

