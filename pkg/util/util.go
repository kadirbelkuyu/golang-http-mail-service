package util

import (
	"encoding/json"
	"log"
	"net/http"
)

// RespondWithError, HTTP yanıtı olarak bir hata mesajı gönderir
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON, HTTP yanıtı olarak JSON gönderir
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// LogError, hata mesajlarını loglar
func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}
