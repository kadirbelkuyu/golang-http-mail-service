package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kadirbelkuyu/mail-service/pkg/config"
	"github.com/kadirbelkuyu/mail-service/pkg/services"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaProducer struct {
	Writer  *kafka.Writer
	Channel *chan bool
}

type KafkaConsumer struct {
	Reader  *kafka.Reader
	Channel *chan bool
}

func NewKafkaProducer(brokers []string, topic string, channel *chan bool) *KafkaProducer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &KafkaProducer{
		Writer:  w,
		Channel: channel,
	}
}

//func NewKafkaConsumer(brokers []string, topic string, channel *chan bool) *KafkaConsumer {
//	r := kafka.NewReader(kafka.ReaderConfig{
//		Brokers: brokers,
//		Topic:   topic,
//	})
//
//	return &KafkaConsumer{
//		Reader:  r,
//		Channel: channel,
//	}
//}

func Consume(ctx context.Context, brokers []string, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "mail-service",
	})
	defer r.Close()

	for {
		fmt.Printf("Çalıştı")
		select {
		case <-ctx.Done():
			return
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// Log error and continue listening, or handle it as needed
				continue
			}
			fmt.Printf("Message: %s\n", string(m.Value))
		}
	}
}

func (kp *KafkaProducer) SendMessage(ctx context.Context, key string, message []byte) error {
	*kp.Channel <- true
	return kp.Writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: message,
		},
	)
}

func (kc *KafkaConsumer) ReadMessage(ctx context.Context, cfg *config.Config) {

	go func() {
		fmt.Printf("Çalıştı")
		reader, err := kc.Reader.ReadMessage(ctx)
		if err != nil {
			*kc.Channel <- false
		}
		if reader.Value == nil {
			*kc.Channel <- false
		}
		var m MessageModel
		json.Unmarshal(reader.Value, &m)
		services.SendEmail(cfg, m.to, m.subject, m.body)
		time.Sleep(time.Millisecond * 25)
	}()
	//<-*kc.Channel
}

type MessageModel struct {
	to, subject, body string
}
