package main

import (
	"config"
	"log"
	"processing"
)

func main() {
	cfg, err := config.Parse("config/config.json")
	if err != nil {
		log.Fatal("Error reading config", err)
	}
	if cfg.Type == "Historical" {
		processing.RunHistorical(cfg)
	} else {
		processing.RunDaily(cfg)
	}
}
