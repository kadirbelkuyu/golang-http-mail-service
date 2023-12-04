package main

import (
	"context"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/email"
	"github.com/kadirbelkuyu/mail-service/pkg/util"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kp := util.NewKafkaProducer(&ctx, cfg.KafkaBrokers, cfg.KafkaTopic)

	kc := util.NewKafkaConsumer(ctx, cfg.KafkaBrokers, cfg.KafkaTopic)

	go kc.Consume(cfg.KafkaBrokers, cfg.KafkaTopic, cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	http.HandleFunc("/send-email", email.SendEmailHandler(cfg, kp))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
