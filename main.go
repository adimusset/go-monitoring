package main

import (
	"fmt"
	"os"
	"strconv"
)

//conf
const (
	consoleRefreshingTime = 10
	alertDuration         = 120
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Please give 2 arguments, not ", len(args))
		return
	}
	t, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("Please give an integer threshold - ", err.Error())
		return
	}

	logLines := make(chan string)

	reader, err := NewLogReader(args[2], logLines)
	if err != nil {
		fmt.Println("Please a working log file - ", err.Error())
		return
	}

	logs := make(chan Object)
	parser := NewLogParser(logLines, logs)

	go reader.Read()
	go parser.Parse()

	start(logs, t)
}

func start(input chan Object, t int) {
	requests := NewStorage()
	sections := NewStorage()

	alerts := make(chan Alert)

	consumer := NewConsumer(input)
	console := NewConsole(requests, sections, alerts, consoleRefreshingTime)
	statsReporter := NewStatisticsReporter(requests, sections)
	averageAlerter := NewAverageAlerter(t, alertDuration, alerts)

	consumer.Subscribe(statsReporter)
	consumer.Subscribe(averageAlerter)

	go consumer.Consume()

	console.Run()
}
