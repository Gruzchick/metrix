package main

var store = MemStorage{
	gauges:   map[string]float64{},
	counters: map[string]int64{},
}
