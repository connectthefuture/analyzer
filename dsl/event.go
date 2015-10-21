package dsl

import "time"

type Event struct {
	T time.Time
	V float64
}

type Events []Event

func (e Events) Data() Data {
	d := Data{}
	for _, event := range e {
		d = append(d, event.V)
	}
	return d
}

func (e Events) FilterValues(f Filter) Events {
	out := Events{}
	for _, event := range e {
		if f(event.V) {
			out = append(out, event)
		}
	}
	return out
}

func (e Events) FilterTimes(f Filter) Events {
	out := Events{}
	for _, event := range e {
		if f(float64(event.T.UnixNano())) {
			out = append(out, event)
		}
	}
	return out
}
