package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"time"
)

const infoRequest = "request"
const infoSection = "section"

var layout, _ = template.New("console").Parse(`{{.Time}}
Most visited section: {{.Section}} - {{.TimesSection}}
Most done request: {{.Request}} - {{.TimesRequest}}{{if .Alerts}}
Recent alerts: {{.Alerts}}{{end}}
------`)

func printOutput(t time.Time, infos map[string]Counts, alerts []Alert) error {
	countsSection, _ := infos[infoSection]
	countsRequest, _ := infos[infoRequest]
	hour := t.Local().Format("15:04:05")
	input := struct {
		Alerts       string
		Time         string
		TimesSection string
		TimesRequest string
		Section      string
		Request      string
	}{
		Time:         hour,
		Section:      countsSection[0].value,
		Request:      countsRequest[0].value,
		TimesSection: fmt.Sprintf("%d", countsSection[0].n), //always 1 element after method sortTable
		TimesRequest: fmt.Sprintf("%d", countsRequest[0].n),
		Alerts:       alertsToString(alerts),
	}
	b := new(bytes.Buffer)
	err := layout.Execute(b, input)
	if err != nil {
		return err
	}
	fmt.Println(string(b.Bytes()))
	return nil
}

type Console struct {
	refresher *time.Ticker

	alerts       chan Alert
	recentAlerts []Alert
	requests     *Storage //could add more
	sections     *Storage
}

func NewConsole(requests *Storage, sections *Storage, alerts chan Alert,
	refreshingTimeInSeconds int) *Console {
	t := time.NewTicker(time.Duration(refreshingTimeInSeconds) * time.Second)
	return &Console{
		refresher:    t,
		alerts:       alerts,
		recentAlerts: []Alert{},
		requests:     requests,
		sections:     sections,
	}
}

func (c *Console) Run() {
	go c.WatchAlerts()
	for range c.refresher.C {
		infos := make(map[string]Counts)
		infos[infoRequest] = c.requests.GetCounts()
		infos[infoSection] = c.sections.GetCounts()
		sortTable(infos)

		t := time.Now()

		err := printOutput(t, infos, c.recentAlerts)
		if err != nil {
			fmt.Println("Display error, ", err.Error())
		}
	}
}

func sortTable(in map[string]Counts) {
	for key, counts := range in {
		if len(counts) == 0 {
			in[key] = Counts{Count{value: "none", n: 0}}
		}
		sort.Sort(counts)
	}
}

func (c *Console) WatchAlerts() {
	for alert := range c.alerts {
		hour := alert.Date.Local().Format("15:04:05")
		s := ""
		if alert.Up {
			s = fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", alert.Average, hour)
		} else {
			s = fmt.Sprintf("Traffic back to normal - hits = %d, triggered at %s", alert.Average, hour)
		}
		fmt.Println(s)

		c.recentAlerts = append(c.recentAlerts, alert) //could remove old ones
	}
}

func alertsToString(recentAlerts []Alert) string {
	result := ""
	for _, alert := range recentAlerts {
		result = result + alert.String() + "--"
	}
	return result
}
