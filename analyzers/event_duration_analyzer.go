package analyzers

import (
	"fmt"
	"math"
	"sort"
	"time"

	. "github.com/onsi/analyzer/dsl"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"

	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
	"github.com/pivotal-golang/lager/chug"
)

type SignificantEvents map[string]Events

func (s SignificantEvents) OrderedNames() []string {
	names := []string{}
	for name := range s {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (s SignificantEvents) FirstTime() time.Time {
	t := time.Unix(1e10, 0)
	for _, events := range s {
		for _, event := range events {
			if event.T.Before(t) {
				t = event.T
			}
		}
	}
	return t
}

func (s SignificantEvents) LogWithThreshold(threshold float64) {
	for _, message := range s.OrderedNames() {
		events := s[message]
		s := fmt.Sprintf("%s %s", events.Data().Stats(), message)
		if events.Data().Max() > threshold {
			s = say.Red(s)
		}
		say.Println(0, s)
	}
}

func ExtractSignificantEvents(entries []chug.LogEntry) SignificantEvents {
	say.Println(0, "Log Lines: %d", len(entries))
	say.Println(0, "Spanning: %s", entries[len(entries)-1].Timestamp.Sub(entries[0].Timestamp))

	bySession := map[string][]chug.LogEntry{}
	for _, entry := range entries {
		session := entry.Source + "-" + entry.Session
		bySession[session] = append(bySession[session], entry)
	}

	say.Println(0, "Sessions: %d", len(bySession))

	groupedEvents := map[string]Events{}
	for _, entries := range bySession {
		duration := entries[len(entries)-1].Timestamp.Sub(entries[0].Timestamp)
		timestamp := entries[0].Timestamp.Add(duration / 2)
		message := entries[0].Message
		groupedEvents[message] = append(groupedEvents[message], Event{T: timestamp, V: duration.Seconds()})
	}

	say.Println(0, "Grouped Events: %d", len(groupedEvents))

	significantEvents := SignificantEvents{}
	messages := []string{}
	for message, events := range groupedEvents {
		if len(events) > 1 {
			significantEvents[message] = events
			messages = append(messages, message)
		}
	}

	say.Println(0, "Significant Events: %d", len(significantEvents))

	return significantEvents
}

type SignificantEventsOptions struct {
	MarkedEvents map[string]plot.LineStyle
	MinX         float64
	MaxX         float64
}

func VisualizeSignificantEvents(events SignificantEvents, filename string, options SignificantEventsOptions) {
	firstTime := events.FirstTime()
	tr := func(t time.Time) float64 {
		return t.Sub(firstTime).Seconds()
	}

	minX := options.MinX
	maxX := 0.0
	maxY := 0.0
	minY := math.MaxFloat64

	histograms := map[string]plot.Plotter{}
	scatters := map[string][]plot.Plotter{}
	verticalLines := []plot.Plotter{}

	colorCounter := 0
	for _, name := range events.OrderedNames() {
		xys := plotter.XYs{}
		xErrs := plotter.XErrors{}
		for _, event := range events[name] {
			if event.V > 0 {
				xys = append(xys, struct{ X, Y float64 }{tr(event.T), event.V})
				xErrs = append(xErrs, struct{ Low, High float64 }{-event.V / 2, event.V / 2})

				if tr(event.T) > maxX {
					maxX = tr(event.T)
				}
				if event.V > maxY {
					maxY = event.V
				}
				if event.V < minY {
					minY = event.V
				}
			}
		}

		if len(xys) == 0 {
			say.Println(0, "No data for %s", name)
			continue
		}

		if options.MarkedEvents != nil {
			ls, ok := options.MarkedEvents[name]
			if ok {
				for _, event := range events[name] {
					l := viz.NewVerticalLine(tr(event.T))
					l.LineStyle = ls
					verticalLines = append(verticalLines, l)
				}
			}
		}

		s, err := plotter.NewScatter(xys)
		say.ExitIfError("Couldn't create scatter plot", err)

		s.GlyphStyle = plot.GlyphStyle{
			Color:  viz.OrderedColor(colorCounter),
			Radius: 2,
			Shape:  plot.CircleGlyph{},
		}

		xErrsPlot, err := plotter.NewXErrorBars(struct {
			plotter.XYer
			plotter.XErrorer
		}{xys, xErrs})
		say.ExitIfError("Couldn't create x errors plot", err)
		xErrsPlot.LineStyle = viz.LineStyle(viz.OrderedColor(colorCounter), 1)

		scatters[name] = []plot.Plotter{s, xErrsPlot}

		durations := events[name].Data()
		h := viz.NewHistogram(durations, 20, durations.Min(), durations.Max())
		h.LineStyle = viz.LineStyle(viz.OrderedColor(colorCounter), 1)
		histograms[name] = h
		colorCounter++
	}

	if options.MaxX != 0 {
		maxX = options.MaxX
	}

	maxY = math.Pow(10, math.Ceil(math.Log10(maxY)))
	minY = math.Pow(10, math.Floor(math.Log10(minY)))

	b := &viz.Board{}
	n := len(histograms) + 1
	padding := 0.1 / float64(n-1)
	height := (1.0 - padding*float64(n-1)) / float64(n)
	histWidth := 0.3
	scatterWidth := 0.7
	y := 1 - height - padding - height

	allScatterPlot, _ := plot.New()
	allScatterPlot.Title.Text = "All Events"
	allScatterPlot.X.Label.Text = "Time (s)"
	allScatterPlot.Y.Label.Text = "Duration (s)"
	allScatterPlot.Y.Scale = plot.LogScale
	allScatterPlot.Y.Tick.Marker = plot.LogTicks

	for _, name := range events.OrderedNames() {
		histogram, ok := histograms[name]
		if !ok {
			continue
		}
		scatter := scatters[name]

		allScatterPlot.Add(scatter[0])
		allScatterPlot.Add(scatter[1])

		histogramPlot, _ := plot.New()
		histogramPlot.Title.Text = name
		histogramPlot.X.Label.Text = "Duration (s)"
		histogramPlot.Y.Label.Text = "N"
		histogramPlot.Add(histogram)

		scatterPlot, _ := plot.New()
		scatterPlot.Title.Text = name
		scatterPlot.X.Label.Text = "Time (s)"
		scatterPlot.Y.Label.Text = "Duration (s)"
		scatterPlot.Y.Scale = plot.LogScale
		scatterPlot.Y.Tick.Marker = plot.LogTicks
		scatterPlot.Add(scatter...)
		scatterPlot.Add(verticalLines...)
		scatterPlot.X.Min = minX
		scatterPlot.X.Max = maxX
		scatterPlot.Y.Min = 1e-5
		scatterPlot.Y.Max = maxY

		b.AddSubPlot(histogramPlot, viz.Rect{0, y, histWidth, height})
		b.AddSubPlot(scatterPlot, viz.Rect{histWidth, y, scatterWidth, height})

		y -= height + padding
	}

	allScatterPlot.Add(verticalLines...)
	allScatterPlot.X.Min = minX
	allScatterPlot.X.Max = maxX
	allScatterPlot.Y.Min = 1e-5
	allScatterPlot.Y.Max = maxY
	fmt.Println("all", minX, maxX)

	b.AddSubPlot(allScatterPlot, viz.Rect{histWidth, 1 - height, scatterWidth, height})

	b.Save(16.0, 5*float64(n), filename)
}
