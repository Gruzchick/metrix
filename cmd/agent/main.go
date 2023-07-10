package main

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"log"
)

var host string
var pollInterval int64
var reportInterval int64

var (
	hostFlag           = flag.String("a", "localhost:8080", "IP address and port in 0.0.0.0:0000 format")
	pollIntervalFlag   = flag.Int64("r", 2, "Measure interval")
	reportIntervalFlag = flag.Int64("p", 10, "Report interval")
)

type Config struct {
	Host           string `env:"ADDRESS"`
	PollInterval   int64  `env:"REPORT_INTERVAL"`
	ReportInterval int64  `env:"POLL_INTERVAL"`
}

func main() {
	flag.Parse()

	cfg := Config{}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Host != "" {
		host = cfg.Host
	} else {
		host = *hostFlag
	}

	if cfg.PollInterval != 0 {
		pollInterval = cfg.PollInterval
	} else {
		pollInterval = *pollIntervalFlag
	}

	if cfg.ReportInterval != 0 {
		reportInterval = cfg.ReportInterval
	} else {
		reportInterval = *reportIntervalFlag
	}

	go collectMetrics()
	go sendMetrics()
	select {}
}
