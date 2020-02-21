package mailUtils

import (
	"daosuan/utils/log"
	"strings"
)

const template = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
    <a href="https://www.baidu.com">冲</a>
    /*body*/
</body>
</html>
`


func generateBody(token string) string {
	logUtils.Println(strings.Replace(template, "/*body*/", token, -1))
	return strings.Replace(template, "/*body*/", token, -1)
}