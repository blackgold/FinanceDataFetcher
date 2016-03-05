package main

import (
	"historical"
        "config"
	"log"
)

func main() {
        cfg,err := config.Parse("config/config.json")
	if err != nil {
		log.Fatal("Error reading config",err)
        }
	historical.Run(cfg)
}
