package mail

import (
	"github.com/DaigangLi/godts/conf"
	"github.com/go-mail/mail"
)

func Send(mailConf *conf.Mail) {

	newMail := mail.NewMessage()
	newMail.SetHeader("From", mailConf.UserName)
	newMail.SetHeader("To", mailConf.UserName)
	newMail.SetHeader("Subject", "go dts alarm test from lidg")
	newMail.SetBody("text/html", "<b>Hello All</b> <i>go dts alarm test from lidg</i>!")

	dialer := mail.NewDialer(mailConf.Host, mailConf.Port, mailConf.UserName, mailConf.Password)

	// Send the email to Bob, Cora and Dan.
	if err := dialer.DialAndSend(newMail); err != nil {
		panic(err)
	}

	//dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
}
