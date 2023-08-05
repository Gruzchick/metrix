package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func restoreFromFIle(fileName string) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	data, _, err := bufio.NewReader(file).ReadLine()
	if err != nil {
		fmt.Println(err)
		return
	}

	newStore := MemStorage{
		Gauges:   map[string]float64{},
		Counters: map[string]int64{},
	}

	err = json.Unmarshal(data, &newStore)
	if err != nil {
		fmt.Println(err)
		return
	}

	store = newStore
}

func writeStoreToFile() {
	if storeFileName == "" {
		return
	}

	data, err := json.Marshal(&store)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.OpenFile(storeFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}

	file.Close()
}

func writeStoreToFileByInterval(storeInterval int64) {

	for {
		time.Sleep(time.Duration(storeInterval) * time.Second)

		writeStoreToFile()
	}
}
