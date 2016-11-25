package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

type LogReader struct {
	file      *os.File
	output    chan string
	refresher *time.Ticker
	state     int64
}

func NewLogReader(path string,
	output chan string) (*LogReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	ticker := time.NewTicker(time.Second)
	return &LogReader{
		file:      file,
		output:    output,
		refresher: ticker,
		state:     0,
	}, nil
}

func (p *LogReader) Read() {
	defer p.file.Close()
	for range p.refresher.C {
		newState, err := p.readFrom(p.state)
		if err != nil {
			fmt.Println("Error while reading file ", err.Error())
			return
		}
		p.state = newState
	}
}

func (p *LogReader) readFrom(start int64) (int64, error) {
	if _, err := p.file.Seek(start, 0); err != nil {
		return start, err
	}
	scanner := bufio.NewScanner(p.file)

	pos := start
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		pos += int64(advance)
		return
	}
	scanner.Split(scanLines)

	for scanner.Scan() {
		p.output <- scanner.Text()
	}
	return pos, scanner.Err()
}

var regexLog = regexp.MustCompile(`(.*?) (.*?) (.*?) \[(.*?)\] "(.*?) HTTP\/(.*?)" (.*?) (.*?)$`)

type LogParser struct {
	input  chan string
	output chan Object
}

func NewLogParser(input chan string,
	output chan Object) *LogParser {
	return &LogParser{
		input:  input,
		output: output,
	}
}

func (p *LogParser) Parse() {
	for in := range p.input {
		if !regexLog.MatchString(in) {
			fmt.Println("Could not parse line", in)
			continue
		}
		submatches := regexLog.FindStringSubmatch(in)
		date, err := time.Parse("02/Jan/2006:15:04:05 -0700", submatches[4])
		if err != nil {
			fmt.Println("Could not parse time", submatches[4], " - ", err.Error())
			continue
		}
		//to ignore old logs at start
		if time.Now().Sub(date) > 10*time.Second {
			continue
		}
		object := Object{
			Date:        date,
			RequestLine: submatches[5],
			IP:          submatches[1],
			StatusCode:  submatches[6],
			Size:        submatches[7],
		}
		p.output <- object
	}
}
