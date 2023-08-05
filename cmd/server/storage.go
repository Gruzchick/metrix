package main

import (
	"errors"
	"strconv"
)

type MemStorage struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]int64   `json:"counters"`
}

var store = MemStorage{
	Gauges:   map[string]float64{},
	Counters: map[string]int64{},
}

func getValueAsString(metricType string, name string) (string, error) {
	switch {
	case metricType == gaugeTypeName:
		val, ok := store.Gauges[name]
		if ok {
			return strconv.FormatFloat(val, 'f', -1, 64), nil
		} else {
			return "0", errors.New("")
		}
	case metricType == counterTypeName:
		val, ok := store.Counters[name]
		if ok {
			return strconv.FormatInt(val, 10), nil
		} else {
			return "0", errors.New("")
		}
	}

	return "0", errors.New("")
}
