package email

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/util"
)

var (
	logFile *os.File
	logger  *log.Logger
)

func init() {
	logFile, err := os.OpenFile("logs.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

// SendEmailHandler godoc
// @Summary E-posta gönder
// @Description E-posta göndermek için kullanılır
// @Tags email
// @Accept json
// @Produce json
// @Param email body EmailRequest true "E-posta İsteği"
// @Success 200 {object} map[string]string "Başarı Yanıtı"
// @Failure 400 {object} map[string]string "Hata Yanıtı"
// @Router /send-email [post]
func SendEmailHandler(cfg *config.Config, kp *util.KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ...
		var req util.EmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			util.ErrorHandler(util.NewHTTPError(http.StatusBadRequest, "Error parsing request body"), w, kp)
			return
		}

		go func(req util.EmailRequest) {
			err := kp.SendMessage(r.Context(), "message", req)
			if err != nil {
				logger.Printf("Error sending message to Kafka: %v, EmailRequest: %+v", err, req)
			}
		}(req)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Email sending initiated"})
	}
}
