package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func storeValue(metrics *Metrics) {
	if db != nil {
		switch {
		case metrics.MType == gaugeTypeName:
			row := db.QueryRow(`select id, value FROM gauges WHERE id = $1`, metrics.ID)

			var (
				id    string
				value float64
			)

			err := row.Scan(&id, &value)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)
				return
			}

			if err == sql.ErrNoRows {
				_, err := db.Exec(`insert into gauges (id, value) values ($1, $2)`, metrics.ID, *metrics.Value)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				_, err := db.Exec(`update gauges set id = $1, value=$2 where id=$1`, metrics.ID, *metrics.Value)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		case metrics.MType == counterTypeName:
			row := db.QueryRow(`select id, value FROM counters WHERE id = $1`, metrics.ID)

			var (
				id    string
				value int64
			)

			err := row.Scan(&id, &value)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)
				return
			}

			if err == sql.ErrNoRows {
				_, err := db.Exec(`insert into counters (id, value) values ($1, $2)`, metrics.ID, *metrics.Delta)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				_, err := db.Exec(`update counters set id = $1, value=$2 where id=$1`, metrics.ID, *metrics.Delta+value)
				if err != nil {
					fmt.Println(err)
					return
				}
				delta := *metrics.Delta + value

				metrics.Delta = &delta
			}

		}
	} else {
		switch {
		case metrics.MType == gaugeTypeName:

			store.Gauges[metrics.ID] = *metrics.Value

			value := store.Gauges[metrics.ID]

			metrics.Value = &value
		case metrics.MType == counterTypeName:

			store.Counters[metrics.ID] = store.Counters[metrics.ID] + *metrics.Delta

			delta := store.Counters[metrics.ID]

			metrics.Delta = &delta
		}
	}
}

func getValueAsString(metricType string, name string) (string, error) {
	if db != nil {
		switch {
		case metricType == gaugeTypeName:
			row := db.QueryRow(`select id, value FROM gauges WHERE id = $1`, name)

			var (
				id    string
				value float64
			)

			err := row.Scan(&id, &value)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)
				return "0", err
			}

			return strconv.FormatFloat(value, 'f', -1, 64), nil
		case metricType == counterTypeName:
			row := db.QueryRow(`select id, value FROM counters WHERE id = $1`, name)

			var (
				id    string
				value int64
			)

			err := row.Scan(&id, &value)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err)
				return "0", err
			}

			return strconv.FormatInt(value, 10), nil
		}

	} else {
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
	}

	return "0", errors.New("")
}
