package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/services"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
	ctx    *context.Context
}

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func NewKafkaProducer(ctx *context.Context, brokers []string, topic string) *KafkaProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{
		Writer: w,
		ctx:    ctx,
	}
}

func NewKafkaConsumer(brokers []string, topic string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "mail-service",
	})

	return &KafkaConsumer{
		Reader: r,
	}
}

func (kc *KafkaConsumer) Consume(ctx context.Context, cfg *config.Config) {
	for {
		fmt.Printf("Çalıştı")
		select {
		case <-ctx.Done():
			return
		default:
			m, err := kc.Reader.ReadMessage(ctx)
			if err != nil {
				// Log error and continue listening, or handle it as needed
				continue
			}
			var mess EmailRequest
			json.Unmarshal(m.Value, &mess)
			go services.SendEmail(cfg, mess.To, mess.Subject, mess.Body)
		}
	}
}

func (kp *KafkaProducer) SendMessage(key string, model EmailRequest) {
	x, _ := json.Marshal(model)
	go kp.Writer.WriteMessages(*kp.ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: x,
		},
	)
}

//func (kc *KafkaConsumer) ReadMessage(ctx context.Context, cfg *config.Config) {
//
//	go func() {
//		fmt.Printf("Çalıştı")
//		reader, err := kc.Reader.ReadMessage(ctx)
//		if err != nil {
//			*kc.Channel <- false
//		}
//		if reader.Value == nil {
//			*kc.Channel <- false
//		}
//		var m MessageModel
//		json.Unmarshal(reader.Value, &m)
//		services.SendEmail(cfg, m.To, m.Subject, m.Body)
//		time.Sleep(time.Millisecond * 25)
//	}()
//	//<-*kc.Channel
//}

// EmailRequest, gelen e-posta isteği için bir yapıdır
type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
