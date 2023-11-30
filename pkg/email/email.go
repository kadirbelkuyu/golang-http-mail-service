package email

import (
	"encoding/json"
	"net/http"

	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/util"
)

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

		//kp := util.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)

		if r.Method != "POST" {
			util.ErrorHandler(util.NewHTTPError(http.StatusMethodNotAllowed, "Invalid request method"), w, kp)
			return
		}

		var req util.EmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		//kp.SendMessage(r.Context(), "message", []byte(fmt.Sprintf("%v", req)))
		kp.SendMessage(r.Context(), "message", req)
		if err != nil {
			util.ErrorHandler(util.NewHTTPError(http.StatusBadRequest, "Error parsing request body"), w, kp)
			return
		}

		//err = SendEmail(cfg, req.To, req.Subject, req.Body)
		if err != nil {
			util.ErrorHandler(util.NewHTTPError(http.StatusInternalServerError, "Error sending email"), w, kp)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
	}
}
