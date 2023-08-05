package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type SendMetricsCounterRequest struct {
	ID    string `json:"id"`
	MType string `json:"type"`
	Delta int64  `json:"delta"`
}

type SendMetricsGaugeRequest struct {
	ID    string  `json:"id"`
	MType string  `json:"type"`
	Value float64 `json:"value"`
}

func sendMetrics() {
	client := &http.Client{}

	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)

		metrics := <-metricsChan

		for k, v := range metrics {

			var body interface{}

			switch {
			case v.metricType == gaugeTypeName:
				var value float64

				value, _ = strconv.ParseFloat(v.metricValue, 64)
				body = SendMetricsGaugeRequest{
					ID:    k,
					MType: v.metricType,
					Value: value,
				}
			case v.metricType == counterTypeName:
				var delta, _ = strconv.ParseInt(v.metricValue, 10, 64)

				body = SendMetricsCounterRequest{
					ID:    k,
					MType: v.metricType,
					Delta: delta,
				}
			}

			jsonBody, err := json.Marshal(body)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("jsonBody.len ", len(jsonBody))

			var compressedBodyBuffer bytes.Buffer

			gz := gzip.NewWriter(&compressedBodyBuffer)

			_, err = gz.Write(jsonBody)
			if err != nil {
				fmt.Println(err)
				return
			}

			gz.Close()

			fmt.Println(len(compressedBodyBuffer.Bytes()))

			request, err := http.NewRequest(http.MethodPost, "http://"+host+"/update/", bytes.NewBuffer(compressedBodyBuffer.Bytes()))
			if err != nil {
				panic(err)
			}

			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Content-Encoding", "gzip")

			resp, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
			} else {
				err = resp.Body.Close()
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
