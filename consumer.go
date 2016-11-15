package main

import "sort"

type Consumer struct {
	input     chan Object
	reporters []Reporter
}

func NewConsumer(input chan Object) *Consumer {
	return &Consumer{input: input, reporters: []Reporter{}}
}

func (c *Consumer) Consume() {
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
}

type StatisticsReporter struct {
	input   chan Object
	pull    chan bool
	output  chan Statistics
	objects []Object
}

func NewStatisticsReporter(pull chan bool, output chan Statistics) *StatisticsReporter {
	return &StatisticsReporter{
		objects: []Object{},
		input:   make(chan Object),
		pull:    pull,
		output:  output,
	}
}

func (r *StatisticsReporter) Input() chan Object {
	return r.input
}

func (r *StatisticsReporter) Consume() {
	for object := range r.input {
		r.objects = append(r.objects, object)
	}
}

func computeStats(objects []Object) Statistics { //test
	stats := Statistics{Sections: make(Counts, 0)}
	for _, object := range objects {
		stats.Sections.add(object.Section())
	}
	sort.Sort(stats.Sections)
	return stats
}

func (r *StatisticsReporter) Serve() {
	for range r.pull {
		r.output <- computeStats(r.objects)
		r.objects = []Object{}
	}
}

type Statistics struct {
	Sections Counts
}
