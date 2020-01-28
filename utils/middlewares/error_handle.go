package middlewares

import (
	"github.com/kataras/iris"
	"daosuan/models"
)

func AbnormalHandle(ctx iris.Context) {
	defer func() {
		re := recover()
		if re == nil {
			return
		}
		switch result := re.(type) {
		case models.RestfulAPIResult:
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(result)
		default:
			panic(re)
			//ctx.StatusCode(iris.StatusInternalServerError)
			//ctx.JSON(models.RestfulAPIResult{
			//	Status: false,
			//	ErrCode: 500,
			//	Message: fmt.Sprintf("%v", re),
			//})
		}
	}()
	ctx.Next()
}

