package main

import (
	"fmt"
	"html/template"
	"os"
	"sort"
	"time"
)

const refreshingTimeInSeconds = 3

var layout, _ = template.New("console").Parse(`
{{.Time}}
Most visited section: {{.Section}} : {{.Times}} times
Last alerts: {{.Alerts}}
------`)

func printOutput(t, section string, times int, alerts string) error {
	input := struct {
		Time    string
		Section string
		Times   string
		Alerts  string
	}{
		Time:    t,
		Section: section,
		Times:   fmt.Sprintf("%d", times),
		Alerts:  alerts,
	}
	return layout.Execute(os.Stdout, input)
}

type Console struct {
	refresher    *time.Ticker
	requests     *Storage
	alerts       chan Alert
	recentAlerts []Alert
}

func NewConsole(requests *Storage, alerts chan Alert) *Console {
	t := time.NewTicker(time.Duration(refreshingTimeInSeconds) * time.Second)
	return &Console{
		refresher:    t,
		requests:     requests,
		alerts:       alerts,
		recentAlerts: []Alert{},
	}
}

func (c *Console) Run() {
	go c.WatchAlerts()
	for range c.refresher.C {
		counts := c.requests.GetCounts()
		sort.Sort(counts)
		// need to add intelligence

		t := time.Now().Format("15:04:05")
		if len(counts) == 0 {
			printOutput(t, "none", 0, "")
			continue
		}
		printOutput(t, counts[0].value, counts[0].n, "")
	}
}

func (c *Console) WatchAlerts() {
	for alert := range c.alerts {
		fmt.Println(alert)
		c.recentAlerts = append(c.recentAlerts, alert)
	}
}
