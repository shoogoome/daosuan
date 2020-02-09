package middlewares

import (
	"daosuan/models"
	"fmt"
	"github.com/kataras/iris"
	"runtime"
)

// 异常控制器
func AbnormalHandle(ctx iris.Context) {
	defer func() {
		re := recover()
		if re == nil {
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		// 打印堆栈信息
		if debug, err := ctx.URLParamInt("debug"); err == nil && debug == 1 {
			ctx.Text(stack())
			return
		}
		// 输出api格式反馈
		switch result := re.(type) {
		case models.RestfulAPIResult:
			ctx.JSON(result)
		default:
			ctx.JSON(models.RestfulAPIResult{
				Status: false,
				ErrCode: 500,
				Message: fmt.Sprintf("系统错误: %v", result),
			})
		}
	}()
	ctx.Next()
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}

