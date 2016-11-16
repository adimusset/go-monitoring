package main

import (
	"sync"
	"time"
)

//the log
type Object struct {
	Date        time.Time
	RequestLine string
}

//thread safe sortable storage
type Storage struct {
	lock    sync.Mutex
	indexes map[string]int
	counts  Counts
}

func NewStorage() *Storage {
	return &Storage{
		lock:    sync.Mutex{},
		indexes: make(map[string]int),
		counts:  make(Counts, 0),
	}
}

func (s *Storage) Add(v string) {
	s.lock.Lock()
	i, ok := s.indexes[v]
	if !ok {
		s.counts = append(s.counts, Count{0, v})
		s.indexes[v] = len(s.counts) - 1
		i = len(s.counts) - 1
	}
	s.counts[i].n = s.counts[i].n + 1
	s.lock.Unlock()
}

func (s *Storage) GetCounts() Counts {
	s.lock.Lock()
	counts := s.counts
	s.indexes = make(map[string]int)
	s.counts = Counts{}
	s.lock.Unlock()
	return counts
}

type Count struct {
	n     int
	value string // could be more complicated
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
