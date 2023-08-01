package main

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var host string

var hostFlag = flag.String("a", "localhost:8080", "IP address and port in 0.0.0.0:0000 format")

type Config struct {
	Host string `env:"ADDRESS"`
}

func main() {
	// создаём предустановленный регистратор zap
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()

	flag.Parse()

	var cfg Config

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Host != "" {
		host = cfg.Host
	} else {
		host = *hostFlag
	}

	router := chi.NewRouter()

	router.Get("/", withLogging(gzipHandle(getAllMetricsHandler)))
	router.Post("/value/", withLogging(gzipHandle(getMetricValueHandlerByPOST)))
	router.Get("/value/{metricType}/{metricName}", withLogging(gzipHandle(getMetricValueHandler)))

	router.Post("/update/", withLogging(gzipHandle(updateMetricsByJSONHandler)))
	router.Post("/update/{metricType}/{metricName}/{metricValue}", withLogging(gzipHandle(updateMetricsHandler)))

	s := &http.Server{
		Addr:           host,
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
