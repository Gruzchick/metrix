package main

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func sendMetrics() {
	client := &http.Client{}

	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)

		metrics := <-metricsChan

		var metricsArray = make([]Metrics, 0, len(metrics))

		for k, v := range metrics {

			var body Metrics

			switch {
			case v.metricType == gaugeTypeName:
				var value float64

				value, _ = strconv.ParseFloat(v.metricValue, 64)
				body = Metrics{
					ID:    k,
					MType: v.metricType,
					Value: &value,
				}

				metricsArray = append(metricsArray, body)
			case v.metricType == counterTypeName:
				var delta, _ = strconv.ParseInt(v.metricValue, 10, 64)

				body = Metrics{
					ID:    k,
					MType: v.metricType,
					Delta: &delta,
				}

				metricsArray = append(metricsArray, body)
			}
		}

		jsonBody, err := json.Marshal(metricsArray)
		if err != nil {
			fmt.Println(err)
		}

		var compressedBodyBuffer bytes.Buffer

		gz := gzip.NewWriter(&compressedBodyBuffer)

		_, err = gz.Write(jsonBody)
		if err != nil {
			fmt.Println(err)
			return
		}

		gz.Close()

		request, err := http.NewRequest(http.MethodPost, "http://"+host+"/updates/", bytes.NewBuffer(compressedBodyBuffer.Bytes()))
		if err != nil {
			panic(err)
		}

		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Content-Encoding", "gzip")

		if hashKey != "" {
			h := hmac.New(sha256.New, []byte(hashKey))
			h.Write(jsonBody)
			hash := h.Sum(nil)
			request.Header.Set("HashSHA256", hex.EncodeToString(hash))
		}

		resp, err := retryRequest(func() (*http.Response, error) { return client.Do(request) })
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

func retryRequest(cb func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 1; i <= 5; i += 2 {
		resp, err = cb()
		if err != nil {
			time.Sleep(time.Duration(i) * time.Second)
		} else {
			return resp, err
		}
	}

	return resp, err
}
