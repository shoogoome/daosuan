package mailUtils

import (
	"strings"
)

const template = `
<!DOCTYPE html>
<html lang="en" style="margin: 0; background: #f5f5f5 !important">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮箱验证</title>
</head>
<body>
    <table style=" border-collapse: collapse; 
        width: 540px; 
        margin: 0 auto; 
        top: 100px; 
        margin-bottom: 100px;
        position: relative; 
        background: white; 
        border-radius: 4px; 
        box-shadow: 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 9px 28px 8px rgba(0, 0, 0, 0.05)
    ">
        <tbody>
            <tr>
                <td style="font-size: 32px; 
                    line-height: 36px; 
                    height: 40px; 
                    background: #40a9ff; 
                    color: white; 
                    border-radius: 4px 4px 0 0; 
                    padding: 15px 24px 8px;
                    margin-block-start: 0.67em;
                    margin-block-end: 0.67em;
                    margin-inline-start: 0px;
                    margin-inline-end: 0px;
                    font-weight: bold;
                ">
					<svg id="图层_1" data-name="图层 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40.61 40.77"><defs><style>.cls-1{fill:#13327d;}svg {width: 40px; height: 40px}</style></defs><title>daosuan</title><g id="图层_1-2" data-name="图层 1"><path class="cls-1" d="M43.75,33.15a16.14,16.14,0,0,1-3.1,9.38,17,17,0,0,1-8.18,6.07,17.57,17.57,0,0,1-5.14,1,19.37,19.37,0,0,1-2.63-.12,17.52,17.52,0,0,1-11.41-6.15,19.29,19.29,0,0,1-1.56-2.15,16.58,16.58,0,0,1-1.22-2.37,17.66,17.66,0,0,1,.75-15.39,17.45,17.45,0,0,1,7.43-7.18,16.55,16.55,0,0,1,4.91-1.59,15.48,15.48,0,0,1,2.55-.21,17.94,17.94,0,0,1,2.53.18l-.5.77c-.55-1.08-1.07-2.17-1.6-3.27-.27-.54-.53-1.09-.79-1.64L25,8.79l1,1.5c.34.5.69,1,1,1.51.68,1,1.36,2,2,3l.53.8-1,0A16.43,16.43,0,0,0,19.66,18a17.62,17.62,0,0,0-1.88,1.32,19,19,0,0,0-1.67,1.57,16.47,16.47,0,0,0-2.56,3.73,15.61,15.61,0,0,0,9.31,21.81,18.5,18.5,0,0,0,2.2.52,19,19,0,0,0,2.26.2,15.81,15.81,0,0,0,4.52-.57A16.21,16.21,0,0,0,36,44.69a18.09,18.09,0,0,0,1.89-1.36,20,20,0,0,0,1.69-1.63,16.28,16.28,0,0,0,1.45-1.86l.33-.5c.1-.17.19-.35.3-.52s.2-.34.29-.52.19-.35.27-.53a18.2,18.2,0,0,0,.9-2.25A21.59,21.59,0,0,0,43.75,33.15Z" transform="translate(-9.1 -8.79)"/><path class="cls-1" d="M21.91,39.33l1.64-1.95,1.66-1.94,3.31-3.89,3.34-3.87,3.36-3.84c1.12-1.28,2.21-2.59,3.32-3.89S40.79,17.4,42,16.19s2.41-2.4,3.67-3.55,2.54-2.29,3.91-3.34l.15.13c-.94,1.45-2,2.83-3,4.17s-2.13,2.67-3.24,4S41.2,20.1,40,21.31s-2.42,2.4-3.59,3.62l-3.56,3.66-3.59,3.64-3.59,3.63-1.79,1.81-1.81,1.79Z" transform="translate(-9.1 -8.79)"/></g></svg>
                </td>
            </tr>
            <tr>
                <td style="padding-top: 24px; 
                    text-align: center; 
                    line-height: 32px; 
                    font-size: 16px;
                    margin-block-start: 1em;
                    margin-block-end: 1em;
                    margin-inline-start: 0px;
                    margin-inline-end: 0px;
                ">
                    欢迎使用捣蒜，请验证邮箱。
					<br/>
					如果此操作非您本人发起的，请忽略此邮件。
                </td>
            </tr>
            <tr>
                <td style="text-align: center; 
                    font-size: 28px;
                    line-height: 30px; 
                    padding: 32px 0;
                    letter-spacing: 4px;
                    margin-block-start: 0.83em;
                    margin-block-end: 0.83em;
                    margin-inline-start: 0px;
                    margin-inline-end: 0px;
                    font-weight: bold;
                ">
                    /*token*/
                </td>
            </tr>
            <tr>
                <td style="
                    text-align: center;
                    font-size: 14px;
                    color: rgb(245, 34, 45);
                    margin: 0px 24px;
                    padding: 24px 0px 32px;
                    line-height: 28px;
                    margin-block-start: 1em;
                    margin-block-end: 1em;
                    margin-inline-start: 0px;
                    margin-inline-end: 0px;
                ">
                    验证码2小时内有效，请尽快完成验证操作。
                </td>
            </tr>
            <tr>
                <td style="
                    margin: 0px 24px;
                    border-top: 1px solid rgb(240, 240, 240);
                    padding: 12px 0px 16px;
                    text-align: center;
                    color: rgb(191, 191, 191);
                    font-size: 14px;
                ">
                    本邮件由系统自动发出，请勿回复！
                </td>
            </tr>
        </tbody>
    </table>
</body>
</html>
`


func generateBody(token string) string {
	return strings.Replace(template, "/*token*/", token, -1)
}