package main

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

type MetricsUpdatingURLPathParams struct {
	action      string
	metricType  string
	metricName  string
	metricValue string
}

type ParsingURLPathParamsError struct {
	error error
	code  int
}
