package main

import (
	"log"
	"vtbot/internal/bot"
	"vtbot/internal/config"
	"vtbot/internal/virustotal"
)

func main() {
	cfg := config.Load()

	vtClient := virustotal.NewClient(cfg.VTApiKey)

	b, err := bot.New(cfg, vtClient)
	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
