package main

import "math/big"

const pollInterval = 2
const reportInterval = 10

var pollCount big.Int

const gaugeTypeName = "gauge"
const counterTypeName = "counter"

var metricsChan = make(chan map[string]metric)
