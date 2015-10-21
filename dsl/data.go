package dsl

import (
	"math"
)

type Data []float64

func (ds Data) Stats() Stats {
	return Stats{
		Min:   ds.Min(),
		Max:   ds.Max(),
		Mean:  ds.Mean(),
		Count: len(ds),
	}
}

func (ds Data) Sum() float64 {
	sum := 0.0
	for _, d := range ds {
		sum += d
	}
	return sum
}

func (ds Data) Mean() float64 {
	return ds.Sum() / float64(len(ds))
}

func (ds Data) Min() float64 {
	if len(ds) == 0 {
		return 0
	}

	min := math.MaxFloat64
	for _, d := range ds {
		if d < min {
			min = d
		}
	}
	return min
}

func (ds Data) Max() float64 {
	if len(ds) == 0 {
		return 0
	}

	max := -math.MaxFloat64
	for _, d := range ds {
		if d > max {
			max = d
		}
	}
	return max
}

func (ds Data) CountInRange(low float64, high float64) int {
	count := 0
	for _, d := range ds {
		if low <= d && d < high {
			count += 1
		}
	}
	return count
}

func (ds Data) Filter(f Filter) Data {
	out := Data{}
	for _, d := range ds {
		if f(d) {
			out = append(out, d)
		}
	}
	return out
}
