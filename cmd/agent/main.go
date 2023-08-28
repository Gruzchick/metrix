package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v9"
	"log"
)

var host string
var pollInterval int64
var reportInterval int64
var hashKey string
var rateLimit int64

var (
	hostFlag           = flag.String("a", "localhost:8080", "IP address and port in 0.0.0.0:0000 format")
	pollIntervalFlag   = flag.Int64("r", 2, "Measure interval")
	reportIntervalFlag = flag.Int64("p", 10, "Report interval")
	hashKeyFlag        = flag.String("k", "", "Ключ для вычисления хеша")
	rateLimitFlag      = flag.Int64("l", 1, "Количество одновременно исходящих запросов на сервер")
)

type Config struct {
	Host           *string `env:"ADDRESS"`
	PollInterval   *int64  `env:"REPORT_INTERVAL"`
	ReportInterval *int64  `env:"POLL_INTERVAL"`
	HashKey        *string `env:"KEY"`
	RateLimit      *int64  `env:"RATE_LIMIT"`
}

func main() {
	flag.Parse()

	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Host != nil {
		host = *cfg.Host
	} else {
		host = *hostFlag
	}

	if cfg.PollInterval != nil {
		pollInterval = *cfg.PollInterval
	} else {
		pollInterval = *pollIntervalFlag
	}

	if cfg.ReportInterval != nil {
		reportInterval = *cfg.ReportInterval
	} else {
		reportInterval = *reportIntervalFlag
	}

	if cfg.HashKey != nil {
		hashKey = *cfg.HashKey
	} else {
		hashKey = *hashKeyFlag
	}

	if cfg.RateLimit != nil {
		rateLimit = *cfg.RateLimit
	} else {
		rateLimit = *rateLimitFlag
	}

	fmt.Println("hashKey", hashKey)

	metricsChan := make(chan map[string]metric, rateLimit)

	go collectMetrics(metricsChan)
	go collectPsutilMetrics(metricsChan)

	for i := int64(1); i <= rateLimit; i++ {
		go sendMetrics(metricsChan)
	}

	select {}
}
