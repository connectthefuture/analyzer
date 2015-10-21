package dsl

import "fmt"

type Stats struct {
	Min   float64
	Max   float64
	Mean  float64
	Count int
}

func (s Stats) String() string {
	return fmt.Sprintf("n=%d [%.4f, <%.4f>, %.4f]", s.Count, s.Min, s.Mean, s.Max)
}
