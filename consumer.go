package main

import (
	"fmt"
	"time"
)

type Consumer struct {
	input     chan Object
	reporters []Reporter
}

func NewConsumer(input chan Object) *Consumer {
	return &Consumer{input: input, reporters: []Reporter{}}
}

func (c *Consumer) Consume() {
	for _, reporter := range c.reporters {
		go reporter.Consume()
	}
	for obj := range c.input {
		for _, reporter := range c.reporters {
			reporter.Input() <- obj
		}
	}
}

func (c *Consumer) Subscribe(reporter Reporter) {
	c.reporters = append(c.reporters, reporter)
}

type Reporter interface {
	Input() chan Object
	Consume()
}

type StatisticsReporter struct {
	input    chan Object
	requests *Storage
	sections *Storage
}

func NewStatisticsReporter(requests *Storage,
	sections *Storage) *StatisticsReporter {
	return &StatisticsReporter{
		requests: requests,
		sections: sections,
		input:    make(chan Object),
	}
}

func (r *StatisticsReporter) Input() chan Object {
	return r.input
}

func (r *StatisticsReporter) Consume() { //this could be several routines
	for object := range r.input {
		r.requests.Add(object.RequestLine)
		section, err := getSection(object.RequestLine)
		if err != nil {
			fmt.Println("Error parsing section, ", object.RequestLine)
			continue
		}
		r.sections.Add(section)
	}
}

//not thread safe
type AverageAlerter struct {
	input             chan Object
	maxAverage        int
	durationInSeconds int
	output            chan Alert
	objects           []Object //only the time is important
	refresher         *time.Ticker
	overAverage       bool
}

type Alert struct {
	Date    time.Time
	Average int
	Up      bool
}

func (a Alert) String() string {
	date := a.Date.Local().Format("02.01.06 15:04:05")
	if a.Up {
		return fmt.Sprintf("%s above %d", date, a.Average)
	}
	return fmt.Sprintf("%s below %d", date, a.Average)
}

func NewAverageAlerter(max, durationInSeconds int,
	output chan Alert) *AverageAlerter {
	t := time.NewTicker(time.Second)
	return &AverageAlerter{
		maxAverage:  max,
		output:      output,
		input:       make(chan Object),
		overAverage: false,
		objects:     []Object{},
		refresher:   t,
	}
}

func (a *AverageAlerter) Input() chan Object {
	return a.input
}

func (a *AverageAlerter) Consume() {
	go a.Run()
	for o := range a.input {
		a.objects = append(a.objects, o)
	}
}

func (a *AverageAlerter) Run() {
	for range a.refresher.C {
		objects, overAverage, alert := nextState(time.Now(), a.objects, a.overAverage,
			a.maxAverage, a.durationInSeconds)
		a.objects = objects
		a.overAverage = overAverage
		if alert != nil {
			a.output <- *alert
		}
	}
}

func nextState(now time.Time, objects []Object, overAverage bool, maxAverage,
	durationInSeconds int) ([]Object, bool, *Alert) {
	i := 0
	//we suppose that the array is sorted for the date (we could also sort it)
	for k, o := range objects {
		if o.Date.After(now.Add(time.Duration(-durationInSeconds) * time.Second)) {
			i = k
			break
		}
	}
	objects = objects[i:]
	m := len(objects) //raw count, could be a moving average
	var alert *Alert
	if m > maxAverage && !overAverage {
		alert = &Alert{now, maxAverage, true}
		overAverage = true
	}
	if m < maxAverage && overAverage {
		alert = &Alert{now, maxAverage, false}
		overAverage = false
	}
	return objects, overAverage, alert
}
