package main

import "time"

type Object struct {
	Date        time.Time
	RequestLine string
}

func (o Object) Section() string {
	return "section"
}
