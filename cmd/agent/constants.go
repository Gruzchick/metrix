package main

import "math/big"

var pollCount big.Int

const gaugeTypeName = "gauge"
const counterTypeName = "counter"

type metric struct {
	metricType  string
	metricValue string
}
