package main

import (
	"bytes"
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

			resp, err := http.Post("http://"+host+"/update/", "application/json", bytes.NewBuffer(jsonBody))
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
