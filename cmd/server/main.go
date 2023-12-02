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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kp := util.NewKafkaProducer(&ctx, cfg.KafkaBrokers, cfg.KafkaTopic)

	kc := util.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaTopic)

	go kc.Consume(ctx, cfg)

	//kc := util.NewKafkaConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, kp.Channel)
	//
	//kc.ReadMessage(ctx, cfg)

	//err := kp.SendMessage(context.Background(), "http-error", []byte("Test message"))
	//if err != nil {
	//	log.Printf("Failed to send message to Kafka: %v", err)
	//}

	//kc.ReadMessage(context.Background(), cfg)
	//r := kafka.NewReader(kafka.ReaderConfig{
	//	Brokers:   cfg.KafkaBrokers,
	//	Topic:     cfg.KafkaTopic,
	//	GroupID:   "mail-service",
	//	Partition: cap(cfg.KafkaBrokers),
	//	MinBytes:  10e3,
	//	MaxBytes:  10e6,
	//})
	//defer r.Close()
	//
	//for {
	//	m, err := r.ReadMessage(context.Background())
	//	if err != nil {
	//		break
	//	}
	//
	//	fmt.Printf("Mesaj: %s\n", string(m.Value))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	http.HandleFunc("/send-email", email.SendEmailHandler(kp))
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
