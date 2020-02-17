package middlewares

import (
	"daosuan/utils/log"
	"github.com/kataras/iris"
)

// 请求日志
func RequestLogHandle(ctx iris.Context) {
	logUtils.Println(ctx.RemoteAddr(), ctx.Method(), ctx.Path())
	ctx.Next()
}
