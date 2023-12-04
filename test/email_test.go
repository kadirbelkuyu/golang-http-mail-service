package email

import (
	"net/smtp"
)

var smtpSendMail = smtp.SendMail

func sendMailTest(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return nil
}
