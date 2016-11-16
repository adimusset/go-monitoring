package main

import (
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
}

func NewStatisticsReporter(requests *Storage) *StatisticsReporter {
	return &StatisticsReporter{
		requests: requests,
		input:    make(chan Object),
	}
}

func (r *StatisticsReporter) Input() chan Object {
	return r.input
}

func (r *StatisticsReporter) Consume() { //this could in several routines
	for object := range r.input {
		r.requests.Add(object.RequestLine)
	}
}

//only for 2 minutess atm
type AverageAlerter struct {
	input       chan Object
	maxAverage  int
	output      chan Alert
	objects     []Object //only the time is important
	refresher   *time.Ticker
	overAverage bool
}

type Alert struct {
	Date    time.Time
	Average int
	Up      bool
}

func NewAverageAlerter(max int, output chan Alert) *AverageAlerter {
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
		i := 0
		for k, o := range a.objects {
			if o.Date.After(time.Now().Add(-2 * time.Minute)) {
				i = k
				break
			}
		}
		if i != 0 {
			a.objects = a.objects[i-1:]
		}
		m := len(a.objects)
		if m > a.maxAverage && !a.overAverage {
			a.output <- Alert{time.Now(), a.maxAverage, true}
			a.overAverage = true
		}
		if m < a.maxAverage && a.overAverage {
			a.output <- Alert{time.Now(), a.maxAverage, false}
			a.overAverage = false
		}
	}
}
