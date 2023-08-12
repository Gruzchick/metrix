package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v9"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

var host string
var storeInterval int64
var storeFileName string
var restoreStoreFromFile bool
var databaseDsn string

var (
	hostFlag                 = flag.String("a", "localhost:8080", "IP address and port in 0.0.0.0:0000 format")
	storeFileNameFlag        = flag.String("f", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения")
	storeIntervalFlag        = flag.Int64("i", 300, "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	restoreStoreFromFileFlag = flag.Bool("r", true, "Определяет загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	databaseDsnFlag          = flag.String("d", "host=localhost user=yandex password=yandex dbname=video sslmode=disable", "Строка подключения к базе данных")
)

type Config struct {
	Host                 *string `env:"ADDRESS"`
	StoreInterval        *int64  `env:"STORE_INTERVAL"`
	StoreFileName        *string `env:"FILE_STORAGE_PATH"`
	RestoreStoreFromFile *bool   `env:"RESTORE"`
	DatabaseDsn          *string `env:"DATABASE_DSN"`
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

	cfg := Config{}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Host != nil {
		host = *cfg.Host
	} else {
		host = *hostFlag
	}

	if cfg.StoreInterval != nil {
		storeInterval = *cfg.StoreInterval
	} else {
		storeInterval = *storeIntervalFlag
	}

	if cfg.StoreFileName != nil {
		storeFileName = *cfg.StoreFileName
	} else {
		storeFileName = *storeFileNameFlag
	}

	if cfg.RestoreStoreFromFile != nil {
		restoreStoreFromFile = *cfg.RestoreStoreFromFile
	} else {
		restoreStoreFromFile = *restoreStoreFromFileFlag
	}

	if cfg.DatabaseDsn != nil {
		databaseDsn = *cfg.DatabaseDsn
	} else {
		databaseDsn = *databaseDsnFlag
	}

	if restoreStoreFromFile && storeFileName != "" {
		restoreFromFIle(storeFileName)
	}

	if storeInterval != 0 {
		go writeStoreToFileByInterval(storeInterval)
	}

	db, err = sql.Open("pgx", databaseDsn)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer db.Close()

	router := chi.NewRouter()

	router.Get("/", withLogging(gzipHandle(getAllMetricsHandler)))
	router.Get("/ping", withLogging(pingDBHandler))
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
		fmt.Println(err)
		panic(err)
	}
}
