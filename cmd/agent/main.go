package main

type metric struct {
	metricType  string
	metricValue string
}

func main() {
	go collectMetrics()
	go sendMetrics(0)
	select {}
}
