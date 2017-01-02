package dsl

import (
	"time"

	"code.cloudfoundry.org/lager/chug"
)

func DurationExtractor(field string) func(chug.LogEntry) float64 {
	return func(e chug.LogEntry) float64 {
		duration, err := time.ParseDuration(e.Data[field].(string))
		if err != nil {
			panic(err)
		}
		return duration.Seconds()
	}
}

func ExtractEventsFromLagerData(entries []chug.LogEntry, f func(chug.LogEntry) float64) Events {
	events := Events{}
	for _, entry := range entries {
		v := f(entry)
		events = append(events, Event{
			T: entry.Timestamp,
			V: v,
		})
	}
	return events
}
