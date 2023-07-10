package main

import (
	"flag"
)

var (
	host           = flag.String("a", "localhost:8080", "IP address and port in 0.0.0.0:0000 format")
	pollInterval   = flag.Int64("r", 2, "Measure interval")
	reportInterval = flag.Int64("p", 10, "Report interval")
)

func main() {
	flag.Parse()

	go collectMetrics()
	go sendMetrics()
	select {}
}
