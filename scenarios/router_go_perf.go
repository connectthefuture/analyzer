package scenarios

import (
	"encoding/csv"
	"flag"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/gonum/plot/plotter"

	"github.com/gonum/plot"

	"github.com/onsi/analyzer/viz"

	. "github.com/onsi/analyzer/dsl"

	"github.com/onsi/analyzer/config"
	"github.com/onsi/say"
)

func GenerateRouterGoPerfCommand() say.Command {
	return say.Command{
		Name:        "router-go-perf",
		Description: "12/15/2016: compare go1.6 to go1.7 performance",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeRouterGoPerf()
		},
	}
}

func analyzeRouterGoPerf() {
	r, err := os.Open(config.DataDir("router-go-perf", "go16.csv"))
	say.ExitIfError("Couldn't read go16.csv", err)
	go16, err := csv.NewReader(r).ReadAll()
	say.ExitIfError("Couldn't parse go16.csv", err)

	r, err = os.Open(config.DataDir("router-go-perf", "go17.csv"))
	say.ExitIfError("Couldn't read go17.csv", err)
	go17, err := csv.NewReader(r).ReadAll()
	say.ExitIfError("Couldn't parse go17.csv", err)

	latencies16 := Data{}
	latencies17 := Data{}

	for i, entry := range go16 {
		if i == 0 {
			continue
		}

		latency, _ := strconv.ParseFloat(entry[1], 64)
		latencies16 = append(latencies16, latency)
	}

	for i, entry := range go17 {
		if i == 0 {
			continue
		}

		latency, _ := strconv.ParseFloat(entry[1], 64)
		latencies17 = append(latencies17, latency)
	}

	say.Println(0, latencies16.Stats().String())
	say.Println(0, latencies17.Stats().String())

	board := viz.NewUniformBoard(2, 1, 0.01)

	p, _ := plot.New()
	p.X.Label.Text = "Time (s)"
	p.Y.Label.Text = "N"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}
	p.X.Min = 0.007
	p.X.Max = 0.4

	style16 := viz.LineStyle(viz.Red, 1.0)
	style17 := viz.LineStyle(viz.Blue, 1.0)

	histogram16 := viz.NewLogHistogram(latencies16, 100, 0.0007, 0.4, 1.0)
	histogram16.LineStyle = style16
	p.Add(histogram16)

	histogram17 := viz.NewLogHistogram(latencies17, 100, 0.0007, 0.4, 1.0)
	histogram17.LineStyle = style17
	p.Add(histogram17)

	board.AddNextSubPlot(p)

	p, _ = plot.New()
	p.X.Label.Text = "100 - percentile"
	p.Y.Label.Text = "Time (s)"

	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}
	p.X.Min = 0.001
	p.X.Max = 100.0

	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}
	p.Y.Min = 0.0007
	p.Y.Max = 0.4

	sort.Float64s(latencies16)
	sort.Float64s(latencies17)

	n := 100
	m := math.Log10(0.001)
	M := math.Log10(90.0)
	delta := (M - m) / float64(n)
	percentiles := []float64{}
	for i := 0; i <= n; i++ {
		percentiles = append(percentiles, math.Pow(10, m+delta*float64(i)))
	}

	p16 := plotter.XYs{}
	p17 := plotter.XYs{}
	for _, x := range percentiles {
		i := int((100.0 - x) / 100.0 * float64(len(latencies16)))
		p16 = append(p16, struct{ X, Y float64 }{x, latencies16[i]})

		i = int((100.0 - x) / 100.0 * float64(len(latencies17)))
		p17 = append(p17, struct{ X, Y float64 }{x, latencies17[i]})
	}

	l16, _ := plotter.NewLine(p16)
	l16.LineStyle = style16
	p.Add(l16)

	l17, _ := plotter.NewLine(p17)
	l17.LineStyle = style17
	p.Add(l17)
	board.AddNextSubPlot(p)

	board.Save(24, 12, "./router-perf.png")
}
