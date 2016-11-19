package main

import (
	"time"
)

func main() {
	input := make(chan Object, 100)
	go logTest(input)

	start(input)
}

func start(input chan Object) {
	requests := NewStorage()
	sections := NewStorage()

	alerts := make(chan Alert)

	consumer := NewConsumer(input)
	console := NewConsole(requests, sections, alerts)
	statsReporter := NewStatisticsReporter(requests, sections)
	averageAlerter := NewAverageAlerter(10, 15, alerts)

	consumer.Subscribe(statsReporter)
	consumer.Subscribe(averageAlerter)

	go consumer.Consume()

	console.Run()
}

func logTest(output chan Object) {
	for k := 0; k < 10; k++ {
		time.Sleep(2 * time.Second)
		output <- Object{RequestLine: "GET /section/page", Date: time.Now()}
		output <- Object{RequestLine: "GET /section2/page", Date: time.Now()}
		output <- Object{RequestLine: "GET /section2/page", Date: time.Now()}
		output <- Object{RequestLine: "GET /section2/page2", Date: time.Now()}
	}
	time.Sleep(30 * time.Second)
	for {
		time.Sleep(time.Second)
		output <- Object{RequestLine: "POST /new/image", Date: time.Now()}
	}
}
