package system

import (
	"crypto/sha1"
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
	if ctx.URLParam("signature") == re {
		ctx.Text(ctx.URLParam("echostr"))
	}
}
