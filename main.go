package main

import (
	"config"
	"log"
	"net"
	"net/http"
	"processing"
	"time"
)

func main() {
	cfg, err := config.Parse("config/config.json")
	if err != nil {
		log.Fatal("Error reading config", err)
	}
	client := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			DisableKeepAlives:     false,
			MaxIdleConnsPerHost:   64,
			ResponseHeaderTimeout: 10 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
		},
	}
	if cfg.Type == "Historical" {
		processing.RunHistorical(cfg, &client)
	} else if cfg.Type == "Daily" {
		processing.RunDaily(cfg, &client)
	} else {
		processing.Run("/Users/devendra/development/finance/FinanceDataFetcher/data")
	}
}
