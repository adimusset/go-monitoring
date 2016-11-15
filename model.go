package main

import "time"

type Object struct {
	Date        time.Time
	RequestLine string
}

func (o Object) Section() Section {
	return Section(o.RequestLine)
}

type Section string

func (s Section) String() string {
	return string(s)
}

type Printable interface {
	String() string
}

type Count struct {
	n     int
	value Printable // or just a string ?
}

type Counts []Count

func (c Counts) Len() int {
	return len(c)
}

func (c Counts) Less(i, j int) bool {
	return c[i].n < c[j].n
}

func (c Counts) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Counts) add(value Printable) {
	for k, old := range c {
		if old.value.String() == value.String() {
			c[k] = Count{old.n + 1, value}
			return
		}
	}
	c = append(c, Count{1, value})
}
