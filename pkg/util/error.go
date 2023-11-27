package util

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (*HTTPError) Error() string {
	panic("unimplemented")
}

func NewHTTPError(statusCode int, message string) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func ErrorHandler(err error, w http.ResponseWriter, kp *KafkaProducer) {
	var httpErr *HTTPError
	if ok := errors.As(err, &httpErr); ok {
		respondWithJSON(w, httpErr.StatusCode, httpErr)
		kp.SendMessage(context.Background(), "http-error", httpErr.Message)
	} else {
		generalErr := HTTPError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Server Error",
		}
		respondWithJSON(w, http.StatusInternalServerError, generalErr)
		kp.SendMessage(context.Background(), "http-error", generalErr.Message)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
