package main

import (
	"errors"
	"strconv"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

var store = MemStorage{
	gauges:   map[string]float64{},
	counters: map[string]int64{},
}

func getValueAsString(metricType string, name string) (string, error) {
	switch {
	case metricType == gaugeTypeName:
		val, ok := store.gauges[name]
		if ok {
			return strconv.FormatFloat(val, 'f', -1, 64), nil
		} else {
			return "0", errors.New("")
		}
	case metricType == counterTypeName:
		val, ok := store.counters[name]
		if ok {
			return strconv.FormatInt(val, 10), nil
		} else {
			return "0", errors.New("")
		}
	}

	return "0", errors.New("")
}
