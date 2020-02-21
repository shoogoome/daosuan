package mailUtils

import (
	"crypto/tls"
	"daosuan/utils"
	"daosuan/utils/log"
	"github.com/go-gomail/gomail"
)

func Send(token string, target ...string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", utils.GlobalConfig.Server.Mail.Username)
	m.SetHeader("To", target...)

	m.SetHeader("Subject", "捣蒜官方邮件")
	m.SetBody("text/html", generateBody(token))
	logUtils.Println(utils.GlobalConfig.Server.Mail.Username, utils.GlobalConfig.Server.Mail.SmtpHost, 25,
		utils.GlobalConfig.Server.Mail.Username, utils.GlobalConfig.Server.Mail.Password)
	d := gomail.NewDialer(
		utils.GlobalConfig.Server.Mail.SmtpHost, 25,
		utils.GlobalConfig.Server.Mail.Username, utils.GlobalConfig.Server.Mail.Password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
