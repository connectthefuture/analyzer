package viz

import (
	"math"

	"code.google.com/p/plotinum/plotter"
	. "github.com/onsi/analyzer/dsl"
)

func NewHistogram(data Data, n int, min float64, max float64) *plotter.Histogram {
	return NewScaledHistogram(data, n, min, max, 1.0)
}

func NewScaledHistogram(data Data, n int, min float64, max float64, weight float64) *plotter.Histogram {
	bins := []plotter.HistogramBin{}

	delta := (max - min) / float64(n)
	for i := 0; i < n; i++ {
		low := min + delta*float64(i)
		high := min + delta*float64(i+1)
		if i == n-1 {
			high = max
		}
		bins = append(bins, plotter.HistogramBin{
			Min:    low,
			Max:    high,
			Weight: float64(data.CountInRange(low, high)) * weight,
		})
	}

	return &plotter.Histogram{
		Bins:      bins,
		LineStyle: plotter.DefaultLineStyle,
	}
}

func NewLogHistogram(data Data, n int, min float64, max float64, weight float64) *plotter.Histogram {
	bins := []plotter.HistogramBin{}
	m := math.Log10(min)
	M := math.Log10(max)
	delta := (M - m) / float64(n)
	for i := 0; i < n; i++ {
		low := math.Pow(10, m+delta*float64(i))
		high := math.Pow(10, m+delta*float64(i+1))
		if i == n-1 {
			high = max
		}
		bins = append(bins, plotter.HistogramBin{
			Min:    low,
			Max:    high,
			Weight: float64(data.CountInRange(low, high)) * weight,
		})
	}

	return &plotter.Histogram{
		Bins:      bins,
		LineStyle: plotter.DefaultLineStyle,
	}
}
