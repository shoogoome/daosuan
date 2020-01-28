package accountException

import (
	"daosuan/models"
)

func NicknameIsExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5301,
		Message: "昵称已存在",
	}
}

func PasswordLengthIsNotStandard() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5302,
		Message: "密码长度应在8到32个字符之间",
	}
}

func InsertAccountFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5302,
		Message: "创建账户失败",
	}
}

func AccountIsNotExists() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5302,
		Message: "账户不存在",
	}
}

func VerificationFail() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5303,
		Message: "用户名或密码错误，或账号（电话、邮箱）未验证。请重新输入",
	}
}

func UsernameLengthIsNotStandard() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5304,
		Message: "用户名应在8到50个字符之间",
	}
}

func OldPasswordIsNotTrue() models.RestfulAPIResult {
	return models.RestfulAPIResult{
		Status: false,
		ErrCode: 5305,
		Message: "旧密码错误",
	}
}