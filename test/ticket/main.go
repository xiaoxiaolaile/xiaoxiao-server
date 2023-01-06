package main

import (
	logs "github.com/sirupsen/logrus"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Millisecond * 5000)
	go runTicker(ticker)
	go runTicker2(ticker)

	select {}
}

func runTicker(ticker *time.Ticker) {
	for range ticker.C {
		logs.Info("定时1")
	}
}

func runTicker2(ticker *time.Ticker) {
	for range ticker.C {
		logs.Info("定时2")
	}
}
