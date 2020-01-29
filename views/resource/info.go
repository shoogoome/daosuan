package resource

import (
	"bytes"
	"daosuan/constants"
	"daosuan/core/auth"
	resourceLogic "daosuan/logics/resource"
	"github.com/kataras/iris"
	"path"
	"strings"
)

// 全局物理资源下载
func LocalDownload(ctx iris.Context, auth authbase.DaoSuanAuthAuthorization, token string) {

	fPath, fileName, tokenFileName := resourceLogic.DecodeToken(token, auth.AccountModel().Id)
	download := ctx.URLParamDefault("download", "0")

	disposition := bytes.Buffer{}
	if download == "1" {
		disposition.Write([]byte("attachment"))
	} else {
		disposition.Write([]byte("inline"))
	}
	nFileName := ""
	if tokenFileName == "" {
		nFileName = ctx.URLParamDefault("filename", fileName)
	} else {
		nFileName = ctx.URLParamDefault("filename", tokenFileName)
	}
	disposition.Write([]byte("; filename=" + nFileName))

	split := strings.Split(nFileName, ".")
	mime, ok := constants.MimeToExtMapping[split[len(split)-1]];
	if !ok {
		mime = "application/octet-stream"
	}
	ctx.ResponseWriter().Header().Set("Via", "ras")
	ctx.ResponseWriter().Header().Set("Content-type", mime)
	ctx.ResponseWriter().Header().Set("Content-Disposition", disposition.String())
	ctx.ResponseWriter().Header().Set(
		"X-Accel-Redirect",
		path.Join(strings.Replace(fPath, constants.DaoSuanDataRoot, constants.NginxResourcePath, -1), fileName))
	defer func() {
		if err := recover(); err != nil {
			ctx.Text("找无该资源")
		}
	}()
}
