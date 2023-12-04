package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/services"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaProducer struct {
	Writer  *kafka.Writer
	Context *context.Context
}

type KafkaConsumer struct {
	Reader  *kafka.Reader
	Context context.Context
}

func NewKafkaProducer(ctx *context.Context, brokers []string, topic string) *KafkaProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{
		Writer:  w,
		Context: ctx,
	}
}

func NewKafkaConsumer(ctx context.Context, brokers []string, topic string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return &KafkaConsumer{
		Reader:  r,
		Context: ctx,
	}
}

func (kc *KafkaConsumer) Consume(brokers []string, topic string, cfg *config.Config) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "mail-service",
	})
	defer r.Close()

	for {
		fmt.Printf("Çalıştı")
		select {
		case <-kc.Context.Done():
			return
		default:
			m, err := r.ReadMessage(kc.Context)
			if err != nil {
				fmt.Printf("Error reading message: %v", err)
				continue
			}
			var mess EmailRequest
			fmt.Printf("%v", m.Value)
			json.Unmarshal(m.Value, &mess)
			services.SendEmail(cfg, mess.To, mess.Subject, mess.Body)
			fmt.Printf("%+v", &mess)
			fmt.Printf("Message: %s\n", string(m.Value))
		}
	}
}

func (kp *KafkaProducer) SendMessage(key string, model EmailRequest) error {
	x, _ := json.Marshal(model)
	err := kp.Writer.WriteMessages(*kp.Context, kafka.Message{Key: []byte(key), Value: x})
	if err != nil {
		log.Printf("Error sending Kafka message: %v", err)
		return err
	}
	return nil
}

type MessageModel struct {
	To, Subject, Body string
}

// EmailRequest, gelen e-posta isteği için bir yapıdır
type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
