package main

import (
	"fmt"
	"net/http"
	"time"
)

func sendMetrics() {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)

		metrics := <-metricsChan

		for k, v := range metrics {
			resp, err := http.Post("http://"+host+"/update/"+v.metricType+"/"+k+"/"+v.metricValue, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}

			err = resp.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
