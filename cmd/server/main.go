package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var store = MemStorage{
	gauges:   map[string]float64{},
	counters: map[string]int64{},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, updateMetricsHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

const gaugeTypeName = "gauge"
const counterTypeName = "counter"

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

func updateMetricsHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.NotFound(res, req)
		return
	}

	//if err := validateHeaders(req.Header); err != nil {
	//	http.Error(res, err.Error(), http.StatusBadRequest)
	//	return
	//}

	parsedURLPathParams, err := getParsedMetricsUpdatingURLPathParams(req.URL.Path)

	if err != nil {
		http.Error(res, err.error.Error(), err.code)
		return
	}

	switch {
	case parsedURLPathParams.metricType == gaugeTypeName:
		metricValue, err := strconv.ParseFloat(parsedURLPathParams.metricValue, 64)

		if err != nil {
			http.Error(res, "Значение метрики для типа "+gaugeTypeName+" должно быть числовой строкой", http.StatusBadRequest)
			return
		}

		store.gauges[parsedURLPathParams.metricName] = metricValue
	case parsedURLPathParams.metricType == counterTypeName:
		metricValue, err := strconv.ParseInt(parsedURLPathParams.metricValue, 10, 64)

		if err != nil {
			http.Error(res, "Значение метрики для типа "+counterTypeName+" должно быть целой числовой строкой", http.StatusBadRequest)
			return
		}

		store.counters[parsedURLPathParams.metricName] = store.counters[parsedURLPathParams.metricName] + metricValue
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
}

//func validateHeaders(headers http.Header) error {
//
//	if headers["Content-Type"][0] != "text/plain" {
//		return errors.New("для заголовка \"Content-Type\" доступно только значение \"text/plain\"")
//	}
//
//	return nil
//}

func getParsedMetricsUpdatingURLPathParams(path string) (*MetricsUpdatingURLPathParams, *ParsingURLPathParamsError) {
	pathComponents := strings.Split(path, "/")

	var action string
	var metricType string
	var metricName string
	var metricValue string

	length := len(pathComponents)

	if length > 1 {
		action = pathComponents[1]
	}

	if length > 2 {
		metricType = pathComponents[2]
	}

	if length > 3 {
		metricName = pathComponents[3]
	}

	if length > 4 {
		metricValue = pathComponents[4]
	}

	if metricName == "" {
		return nil, &ParsingURLPathParamsError{
			error: errors.New("не указано имя метрики"),
			code:  http.StatusNotFound,
		}
	}

	metricsUpdatingURLPathParams := MetricsUpdatingURLPathParams{
		action:      action,
		metricType:  metricType,
		metricName:  metricName,
		metricValue: metricValue,
	}

	if metricsUpdatingURLPathParams.action != "update" {
		return nil, &ParsingURLPathParamsError{
			error: errors.New("не верный формат запроса, url должен быть в формате: /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>"),
			code:  http.StatusBadRequest,
		}
	}

	if metricsUpdatingURLPathParams.metricType != gaugeTypeName && metricsUpdatingURLPathParams.metricType != counterTypeName {
		return nil, &ParsingURLPathParamsError{
			error: errors.New("не корректный тип метрики, доступные значения: " + gaugeTypeName + ", " + counterTypeName),
			code:  http.StatusBadRequest,
		}
	}

	return &metricsUpdatingURLPathParams, nil
}
