package main

import (
	"time"
)

func main() {
	input := make(chan Object, 100)
	go logTest(input)

	puller := make(chan bool)
	stats := make(chan Statistics)

	consumer := NewConsumer(input)
	console := NewConsole(puller, stats)
	statsReporter := NewStatisticsReporter(puller, stats)

	consumer.Subscribe(statsReporter)

	go consumer.Consume()
	go statsReporter.Consume()
	go statsReporter.Serve()

	console.Run()
}

func logTest(output chan Object) {
	for {
		time.Sleep(2 * time.Second)
		output <- Object{RequestLine: "/section"}
	}
}
