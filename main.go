package main

import (
	"time"
)

func main() {
	input := make(chan Object)
	go logTest(input)

	requests := NewStorage()

	alerts := make(chan Alert)

	consumer := NewConsumer(input)
	console := NewConsole(requests, alerts)
	statsReporter := NewStatisticsReporter(requests)
	averageAlerter := NewAverageAlerter(10, alerts)

	consumer.Subscribe(statsReporter)
	consumer.Subscribe(averageAlerter)

	go consumer.Consume()

	console.Run()
}

func logTest(output chan Object) {
	for k := 0; k < 10; k++ {
		time.Sleep(2 * time.Second)
		output <- Object{RequestLine: "/section"}
	}
	for {
		time.Sleep(time.Second)
		output <- Object{RequestLine: "/new"}
	}
}
