package main

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/email"
	"github.com/kadirbelkuyu/mail-service/pkg/util"
)

func main() {
	cfg := config.LoadConfig()

	kp := util.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)

	http.HandleFunc("/send-email", email.SendEmailHandler(cfg, kp))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
