package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Get("/", getAllMetricsHandler)
	router.Get("/value/{metricType}/{metricName}", getMetricValueHandler)

	router.Post("/update/{metricType}/{metricName}/{metricValue}", updateMetricsHandler)

	err := http.ListenAndServe(`:8080`, router)
	if err != nil {
		panic(err)
	}
}
