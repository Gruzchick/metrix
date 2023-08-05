package main

import (
	"math/big"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func collectMetrics() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		pollCount.Add(&pollCount, big.NewInt(1))

		newMetrics := make(map[string]metric)

		newMetrics["Alloc"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.Alloc, 10),
		}
		newMetrics["BuckHashSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.BuckHashSys, 10),
		}
		newMetrics["Frees"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.Frees, 10),
		}
		newMetrics["GCCPUFraction"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatFloat(m.GCCPUFraction, 'f', -1, 64),
		}
		newMetrics["GCSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.GCSys, 10),
		}
		newMetrics["HeapAlloc"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.HeapAlloc, 10),
		}
		newMetrics["HeapIdle"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.HeapIdle, 10),
		}
		newMetrics["HeapInuse"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.HeapInuse, 10),
		}
		newMetrics["HeapObjects"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.HeapObjects, 10),
		}
		newMetrics["HeapReleased"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.HeapReleased, 10),
		}
		newMetrics["HeapSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.HeapSys, 10),
		}
		newMetrics["LastGC"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.LastGC, 10),
		}
		newMetrics["Lookups"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.Lookups, 10),
		}
		newMetrics["MCacheInuse"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.MCacheInuse, 10),
		}
		newMetrics["MCacheSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.MCacheSys, 10),
		}
		newMetrics["MSpanInuse"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.MSpanInuse, 10),
		}
		newMetrics["MSpanSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.MSpanSys, 10),
		}
		newMetrics["Mallocs"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.Mallocs, 10),
		}
		newMetrics["NextGC"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.NextGC, 10),
		}
		newMetrics["OtherSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.OtherSys, 10),
		}
		newMetrics["PauseTotalNs"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.PauseTotalNs, 10),
		}
		newMetrics["StackInuse"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.StackInuse, 10),
		}
		newMetrics["StackSys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.StackSys, 10),
		}

		newMetrics["Sys"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.Sys, 10),
		}
		newMetrics["TotalAlloc"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(m.TotalAlloc, 10),
		}

		newMetrics["PollCount"] = metric{
			metricType:  counterTypeName,
			metricValue: pollCount.String(),
		}
		newMetrics["RandomValue"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatFloat(rand.Float64(), 'f', -1, 64),
		}
		newMetrics["NumForcedGC"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(uint64(m.NumForcedGC), 10),
		}
		newMetrics["NumGC"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(uint64(m.NumGC), 10),
		}

		metricsChan <- newMetrics

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}

}
