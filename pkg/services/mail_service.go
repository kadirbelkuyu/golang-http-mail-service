package services

import (
	"fmt"
	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"net/smtp"
)

func SendEmail(cfg *config.Config, to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", cfg.SenderEmail, to, subject, body)

	auth := smtp.PlainAuth("", cfg.SenderEmail, cfg.SenderPassword, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	return smtp.SendMail(addr, auth, cfg.SenderEmail, []string{to}, []byte(msg))
}
