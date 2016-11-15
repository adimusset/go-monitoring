package main

import (
	"html/template"
	"os"
	"time"
)

const refreshingTimeInSeconds = 3

var layout, _ = template.New("console").Parse(`
{{.Time}}
Most visited section: {{.Section}}
------`)

func printOutput(t, section string) error {
	input := struct {
		Time    string
		Section string
	}{
		Time:    t,
		Section: section,
	}
	return layout.Execute(os.Stdout, input)
}

type Console struct {
	refresher   *time.Ticker
	statsPuller chan bool
	stats       chan Statistics
}

func NewConsole(puller chan bool, stats chan Statistics) *Console {
	t := time.NewTicker(time.Duration(refreshingTimeInSeconds) * time.Second)
	return &Console{
		refresher:   t,
		statsPuller: puller,
		stats:       stats,
	}
}

func (c *Console) Run() {
	for range c.refresher.C {
		c.statsPuller <- true
		stats := <-c.stats
		t := time.Now().Format("15:04:05")
		if len(stats.Sections) == 0 {
			printOutput(t, "none")
			continue
		}
		section := stats.Sections[0].value.String()
		printOutput(t, section)
	}
}
