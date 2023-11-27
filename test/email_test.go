package email

import (
	"errors"
	"net/smtp"
	"testing"

	"github.com/kadirbelkuyu/mail-service/pkg/config"
)

var smtpSendMail = smtp.SendMail

func sendMailTest(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return nil
}

func TestSendEmail_Error(t *testing.T) {
	smtpSendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		return errors.New("SMTP error")
	}
	defer func() { smtpSendMail = smtp.SendMail }()

	cfg := &config.Config{
		SMTPHost:       "smtp.example.com",
		SMTPPort:       "587",
		SenderEmail:    "test@example.com",
		SenderPassword: "password",
	}

	err := sendMail(cfg, "test@example.com", "Test subject", "Test body")
	if err == nil || err.Error() != "SMTP error" {
		t.Errorf("Expected SMTP error, got %v", err)
	}
}
