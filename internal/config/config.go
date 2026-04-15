package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	TelegramToken   string
	VTApiKey        string
	MaxFileSize     int64
	RequestInterval time.Duration
	VerboseOutput   bool
}

func Load() *Config {
	token := os.Getenv("TELEGRAM_TOKEN")
	apiKey := os.Getenv("VT_API_KEY")

	if token == "" || apiKey == "" {
		log.Fatal("Встановіть TELEGRAM_TOKEN та VT_API_KEY")
	}

	return &Config{
		TelegramToken:   token,
		VTApiKey:        apiKey,
		MaxFileSize:     20 * 1024 * 1024, // 20 MB
		RequestInterval: 16 * time.Second, // Free tier: 4 req/min
		VerboseOutput:   os.Getenv("VERBOSE_OUTPUT") == "true",
	}
}
