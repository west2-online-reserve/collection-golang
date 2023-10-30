package myutils

import (
	"four/config"
	"four/consts"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SendEmail(userEmail string, imgPath string) error {
	conf := config.Config.Email["qqmail"]
	em := email.NewEmail()
	em.From = conf.Sender            // 发送人
	em.To = []string{userEmail}      // 收件人
	em.Subject = consts.EmailSubject // subject

	_, err := em.AttachFile(imgPath) // 将存储的照片作为附件发送
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", em.From, conf.Password, conf.Host) // 身份验证

	err = em.Send(conf.Address, auth) // 发送
	if err != nil {
		return err
	}
	return err
}
