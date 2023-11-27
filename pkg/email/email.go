package email

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/util"
)

// EmailRequest, gelen e-posta isteği için bir yapıdır
type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEmailHandler(cfg *config.Config, kp *util.KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			util.ErrorHandler(util.NewHTTPError(http.StatusMethodNotAllowed, "Invalid request method"), w, kp)
			return
		}

		var req EmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			util.ErrorHandler(util.NewHTTPError(http.StatusBadRequest, "Error parsing request body"), w, kp)
			return
		}

		err = SendEmail(cfg, req.To, req.Subject, req.Body)
		if err != nil {
			util.ErrorHandler(util.NewHTTPError(http.StatusInternalServerError, "Error sending email"), w, kp)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
	}
}

func SendEmail(cfg *config.Config, to, subject, body string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", cfg.SenderEmail, to, subject, body)

	auth := smtp.PlainAuth("", cfg.SenderEmail, cfg.SenderPassword, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	return smtp.SendMail(addr, auth, cfg.SenderEmail, []string{to}, []byte(msg))
}
