package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Config struct {
	SMTPHost       string
	SMTPPort       string
	SenderEmail    string
	SenderPassword string
	KafkaBrokers   []string
	KafkaTopic     string
}

func LoadConfig() *Config {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", ""),
		SenderEmail:    getEnv("SENDER_EMAIL", ""),
		SenderPassword: getEnv("SENDER_PASSWORD", ""),
		KafkaBrokers:   strings.Split(getEnv("KAFKA_BROKERS", ""), ","),
		KafkaTopic:     getEnv("KAFKA_TOPIC", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
