package main

import (
	"fmt"
	"net/http"
)

func sendMetrics(counter int) {
	if counter < reportInterval/pollInterval-1 {
		<-metricsChan
		go sendMetrics(counter + 1)

		return
	} else {
		metrics := <-metricsChan

		for k, v := range metrics {
			resp, err := http.Post("http://localhost:8080/update/"+v.metricType+"/"+k+"/"+v.metricValue, "text/plain", nil)
			if err != nil {
				fmt.Println(err)
			}

			err = resp.Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}

		go sendMetrics(0)

		return
	}
}
