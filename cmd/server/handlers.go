package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func getAllMetricsHandler(res http.ResponseWriter, req *http.Request) {
	var html = `<html>
    <head>
    <title></title>
    </head>
    <body>
        <h2>Metrics<h2>`

	if len(store.Gauges) != 0 {
		html += "<h3>Gauges</h3>"

		for k, v := range store.Gauges {
			html += "<div>" + "<span>" + k + ": " + "</span>" + "<span>" + strconv.FormatFloat(v, 'f', -1, 64) + "</span>" + "</div>"
		}

	}
	if len(store.Counters) != 0 {
		html += "<h3>Counters</h3>"

		for k, v := range store.Counters {
			html += "<div>" + "<span>" + k + ": " + "</span>" + "<span>" + strconv.FormatInt(v, 10) + "</span>" + "</div>"
		}

	}

	html += "</body></html>"

	res.Header().Set("content-type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(html))
}

type GetMetricValueHandlerRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

func getMetricValueHandlerByPOST(res http.ResponseWriter, req *http.Request) {
	var requestBody GetMetricValueHandlerRequest
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &requestBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	metricValue, err := getValueAsString(requestBody.MType, requestBody.ID)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	responseBody := Metrics{
		ID:    requestBody.ID,
		MType: requestBody.MType,
	}

	switch {
	case requestBody.MType == gaugeTypeName:
		v, _ := strconv.ParseFloat(metricValue, 64)

		responseBody.Value = &v
	case requestBody.MType == counterTypeName:
		d, _ := strconv.ParseInt(metricValue, 10, 64)

		responseBody.Delta = &d
	}

	resp, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func getMetricValueHandler(res http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	val, err := getValueAsString(metricType, metricName)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(val))
}

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func updateMetricsByJSONHandler(res http.ResponseWriter, req *http.Request) {
	var parsedBody Metrics

	var bodyBytes []byte

	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		defer gz.Close()

		bodyBytes, err = io.ReadAll(gz)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		var buf bytes.Buffer

		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		bodyBytes = buf.Bytes()
	}

	if err := json.Unmarshal(bodyBytes, &parsedBody); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if parsedBody.ID == "" {
		http.Error(res, "не указано имя метрики", http.StatusNotFound)
		return
	}

	switch {
	case parsedBody.MType == gaugeTypeName:

		store.Gauges[parsedBody.ID] = *parsedBody.Value

		value := store.Gauges[parsedBody.ID]

		parsedBody.Value = &value
	case parsedBody.MType == counterTypeName:

		store.Counters[parsedBody.ID] = store.Counters[parsedBody.ID] + *parsedBody.Delta

		delta := store.Counters[parsedBody.ID]

		parsedBody.Delta = &delta
	}

	resp, err := json.Marshal(parsedBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(resp)

	if storeInterval == 0 {
		writeStoreToFileByInterval(storeInterval)
	}
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

		store.Gauges[parsedURLPathParams.metricName] = metricValue
	case parsedURLPathParams.metricType == counterTypeName:
		metricValue, err := strconv.ParseInt(parsedURLPathParams.metricValue, 10, 64)

		if err != nil {
			http.Error(res, "Значение метрики для типа "+counterTypeName+" должно быть целой числовой строкой", http.StatusBadRequest)
			return
		}

		store.Counters[parsedURLPathParams.metricName] = store.Counters[parsedURLPathParams.metricName] + metricValue
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)

	if storeInterval == 0 {
		writeStoreToFileByInterval(storeInterval)
	}
}

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
