package main

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/big"
	"runtime"
	"strconv"
	"time"
)

func collectPsutilMetrics(metricsChan chan<- map[string]metric) {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		pollCount.Add(&pollCount, big.NewInt(1))

		newMetrics := make(map[string]metric)

		memory, _ := mem.VirtualMemory()

		newMetrics["TotalMemory"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(memory.Total, 10),
		}
		newMetrics["FreeMemory"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatUint(memory.Free, 10),
		}

		cpuUtilization, _ := cpu.Percent(time.Duration(1)*time.Second, false)

		newMetrics["CPUutilization1"] = metric{
			metricType:  gaugeTypeName,
			metricValue: strconv.FormatFloat(cpuUtilization[0], 'f', -1, 64),
		}

		select {
		case metricsChan <- newMetrics:
		default:
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
