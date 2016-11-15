package main

type Consumer struct {
	Reporters []Reporter
}

func NewConsumer() *Consumer {
	return &Consumer{[]Reporter{}}
}

func (c *Consumer) Consume(input chan Object) {
	for httpObject := range input {
		for _, reporter := range c.Reporters {
			reporter.Input() <- httpObject
		}
	}
}

func (c *Consumer) Subscribe(reporter Reporter) {
	c.Reporters = append(c.Reporters, reporter)
}

type Reporter interface {
	Input() chan Object
	Consume()
}

type StatisticsReporter struct {
	input   chan Object
	objects []Object
}

func (r *StatisticsReporter) Input() chan Object {
	return r.input
}

func (r *StatisticsReporter) Consume() {
	for object := range r.input {
		r.objects = append(r.objects, object)
	}
}

func (r *StatisticsReporter) computeStats() Statistics {
	return Statistics{}
}

type Statistics struct {
}

type Bin struct {
	count int
	value interface{}
}
