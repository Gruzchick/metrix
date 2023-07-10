package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	"net/http"
)

var host = flag.String("a", "localhost:8080", "IP address and port in 0.0.0.0:0000 format")

func main() {
	flag.Parse()

	router := chi.NewRouter()

	router.Get("/", getAllMetricsHandler)
	router.Get("/value/{metricType}/{metricName}", getMetricValueHandler)

	router.Post("/update/{metricType}/{metricName}/{metricValue}", updateMetricsHandler)

	err := http.ListenAndServe(*host, router)
	if err != nil {
		panic(err)
	}
}
