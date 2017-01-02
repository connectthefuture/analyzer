package analyzers

import (
	"fmt"
	"math"
	"sort"
	"time"

	. "github.com/onsi/analyzer/dsl"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg/draw"

	"code.cloudfoundry.org/lager/chug"
	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
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
	return ExtractSignificantEventsWithThreshold(entries, 1)
}

func ExtractSignificantEventsWithThreshold(entries []chug.LogEntry, n int) SignificantEvents {
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
		if len(events) > n {
			significantEvents[message] = events
			messages = append(messages, message)
		}
	}

	say.Println(0, "Significant Events: %d", len(significantEvents))

	return significantEvents
}

type LineOverlay struct {
	Events    Events
	LineStyle draw.LineStyle
}

type VerticalMarker struct {
	T         time.Time
	LineStyle draw.LineStyle
}

type SignificantEventsOptions struct {
	MarkedEvents    map[string]draw.LineStyle
	VerticalMarkers []VerticalMarker
	LineOverlays    []LineOverlay
	OverlayPlots    []plot.Plotter
	MinX            float64
	MaxX            float64
	MinT            time.Time
	MaxT            time.Time
	MaxY            float64
	WidthStretch    float64
}

func VisualizeSignificantEvents(events SignificantEvents, filename string, options SignificantEventsOptions) {
	firstTime := events.FirstTime()
	tr := func(t time.Time) float64 {
		return t.Sub(firstTime).Seconds()
	}

	minX := 0.0
	if options.MinX != 0 {
		minX = options.MinX
	}
	if !options.MinT.IsZero() {
		minX = tr(options.MinT)
	}
	maxX := 0.0
	maxY := 0.0
	minY := math.MaxFloat64

	scatters := map[string][]plot.Plotter{}
	verticalLines := []plot.Plotter{}
	lineOverlays := []plot.Plotter{}

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

		s.GlyphStyle = draw.GlyphStyle{
			Color:  viz.OrderedColor(colorCounter),
			Radius: 2,
			Shape:  draw.CircleGlyph{},
		}

		xErrsPlot, err := plotter.NewXErrorBars(struct {
			plotter.XYer
			plotter.XErrorer
		}{xys, xErrs})
		say.ExitIfError("Couldn't create x errors plot", err)
		xErrsPlot.LineStyle = viz.LineStyle(viz.OrderedColor(colorCounter), 1)

		scatters[name] = []plot.Plotter{s, xErrsPlot}

		colorCounter++
	}

	for _, marker := range options.VerticalMarkers {
		l := viz.NewVerticalLine(tr(marker.T))
		l.LineStyle = marker.LineStyle
		verticalLines = append(verticalLines, l)
	}

	for _, lineOverlay := range options.LineOverlays {
		xys := plotter.XYs{}
		for _, event := range lineOverlay.Events {
			if event.V > 0 {
				xys = append(xys, struct{ X, Y float64 }{tr(event.T), event.V})
			}
		}

		l, s, err := plotter.NewLinePoints(xys)
		say.ExitIfError("Couldn't create scatter plot", err)

		l.LineStyle = lineOverlay.LineStyle
		s.GlyphStyle = draw.GlyphStyle{
			Color:  lineOverlay.LineStyle.Color,
			Radius: lineOverlay.LineStyle.Width,
			Shape:  draw.CrossGlyph{},
		}
		lineOverlays = append(lineOverlays, l, s)
	}

	if options.MaxX != 0 {
		maxX = options.MaxX
	}
	if !options.MaxT.IsZero() {
		maxX = tr(options.MaxT)
	}

	maxY = math.Pow(10, math.Ceil(math.Log10(maxY)))

	if options.MaxY != 0 {
		maxY = options.MaxY
	}

	minY = math.Pow(10, math.Floor(math.Log10(minY)))

	n := len(scatters) + 1
	b := viz.NewUniformBoard(1, n, 0.01)

	allScatterPlot, _ := plot.New()
	allScatterPlot.Title.Text = "All Events"
	allScatterPlot.X.Label.Text = "Time (s)"
	allScatterPlot.Y.Label.Text = "Duration (s)"
	allScatterPlot.Y.Scale = plot.LogScale{}
	allScatterPlot.Y.Tick.Marker = plot.LogTicks{}

	for i, name := range events.OrderedNames() {
		scatter, ok := scatters[name]
		if !ok {
			continue
		}

		allScatterPlot.Add(scatter[0])
		allScatterPlot.Add(scatter[1])

		scatterPlot, _ := plot.New()
		scatterPlot.Title.Text = name
		scatterPlot.X.Label.Text = "Time (s)"
		scatterPlot.Y.Label.Text = "Duration (s)"
		scatterPlot.Y.Scale = plot.LogScale{}
		scatterPlot.Y.Tick.Marker = plot.LogTicks{}
		scatterPlot.Add(scatter...)
		scatterPlot.Add(verticalLines...)
		scatterPlot.Add(lineOverlays...)
		scatterPlot.Add(options.OverlayPlots...)
		scatterPlot.X.Min = minX
		scatterPlot.X.Max = maxX
		scatterPlot.Y.Min = 1e-5
		scatterPlot.Y.Max = maxY

		b.AddSubPlotAt(scatterPlot, 0, n-2-i)
	}

	allScatterPlot.Add(verticalLines...)
	allScatterPlot.Add(lineOverlays...)
	allScatterPlot.Add(options.OverlayPlots...)
	allScatterPlot.X.Min = minX
	allScatterPlot.X.Max = maxX
	allScatterPlot.Y.Min = 1e-5
	allScatterPlot.Y.Max = maxY
	fmt.Println("all", minX, maxX)

	b.AddSubPlotAt(allScatterPlot, 0, n-1)

	width := 12.0
	if options.WidthStretch > 0 {
		width = width * options.WidthStretch
	}
	b.Save(width, 5*float64(n), filename)
}
