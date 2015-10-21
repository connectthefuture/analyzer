package dsl

import "fmt"

type Stats struct {
	Min   float64
	Max   float64
	Mean  float64
	Count int
}

func (s Stats) String() string {
	count := fmt.Sprintf("n=%d", s.Count)
	return fmt.Sprintf("%8s [%.4f, <%.4f>, %.4f]", count, s.Min, s.Mean, s.Max)
}
