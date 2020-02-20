package system

import (
	"crypto/sha1"
	"daosuan/utils/log"
	"fmt"
	"github.com/kataras/iris"
	"sort"
	"strings"
)

func WeChatV(ctx iris.Context) {

	signatureList := []string{
		ctx.URLParam("timestamp"),
		ctx.URLParam("nonce"),
		ctx.URLParam("echostr")}

	sort.Strings(signatureList)

	h := sha1.New()
	h.Write([]byte(strings.Join(signatureList, "")))
	re := fmt.Sprintf("%x", h.Sum(nil))
	logUtils.Println(signatureList, re, ctx.URLParam("signature"))
	if ctx.URLParam("signature") == re {
		//ctx.Text(ctx.URLParam("echostr"))
	}
	ctx.Text(ctx.URLParam("echostr"))
}
