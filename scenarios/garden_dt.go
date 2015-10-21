package scenarios

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
	"time"

	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"

	"github.com/onsi/analyzer/config"
	. "github.com/onsi/analyzer/dsl"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
	"github.com/pivotal-golang/lager/chug"
)

func GenerateGardenDTCommand() say.Command {
	return say.Command{
		Name:        "garden-dt",
		Description: "10/20/2015: garden sometimes takes forever - can we analyze?",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeGardenDT()
		},
	}
}

func analyzeGardenDT() {
	data, err := ioutil.ReadFile(config.DataDir("garden-dt", "garden-dt.logs"))
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	say.Println(0, "Log Lines: %d", len(entries))
	say.Println(0, "Spanning: %s", entries[len(entries)-1].Timestamp.Sub(entries[0].Timestamp))

	bySession := map[string][]chug.LogEntry{}
	for _, entry := range entries {
		bySession[entry.Session] = append(bySession[entry.Session], entry)
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

	significantEvents := map[string]Events{}
	messages := []string{}
	for message, events := range groupedEvents {
		if len(events) > 1 {
			significantEvents[message] = events
			messages = append(messages, message)
		}
	}

	sort.Strings(messages)

	say.Println(0, "Significant Events: %d", len(significantEvents))
	for _, message := range messages {
		events := significantEvents[message]
		s := fmt.Sprintf("%s %s", events.Data().Stats(), message)
		if events.Data().Max() > 0.2 {
			s = say.Red(s)
		}
		say.Println(0, s)
	}

	generateGardenDTHistograms(messages, significantEvents)
	generateGardenDTScatter(messages, significantEvents)
	generateManyGardenDTScatter(messages, significantEvents)
}

func generateGardenDTHistograms(messages []string, significantEvents map[string]Events) {
	nAcross := 4
	nHigh := int(math.Ceil(float64(len(significantEvents)) / float64(nAcross)))

	board := viz.NewUniformBoard(nAcross, nHigh, 0.01)
	for i, message := range messages {
		durations := significantEvents[message].Data()

		p, _ := plot.New()
		p.Title.Text = strings.Join(strings.Split(message, ".")[1:], ".")
		h := viz.NewHistogram(durations, 20, durations.Min(), durations.Max())
		h.LineStyle = viz.LineStyle(viz.OrderedColor(i), 1)
		p.Add(h)
		board.AddNextSubPlot(p)
	}
	board.Save(float64(nAcross)*4.0, float64(nHigh)*4.0, config.DataDir("garden-dt", "garden-dt.svg"))
}

func generateGardenDTScatter(messages []string, significantEvents map[string]Events) {
	tr := func(t time.Time) float64 {
		return float64(t.UnixNano())/1e9 - 1.445379815e9
	}

	nAcross := 1
	nHigh := 1

	board := viz.NewUniformBoard(nAcross, nHigh, 0)
	p, _ := plot.New()
	p.X.Label.Text = "Time (s)"
	p.Y.Label.Text = "Event Duration (s)"
	p.Y.Scale = plot.LogScale
	p.Y.Tick.Marker = plot.LogTicks
	for i, message := range messages {
		if message == "garden-linux.container.info-starting" {
			continue
		}

		xys := plotter.XYs{}
		xErrs := plotter.XErrors{}
		for _, event := range significantEvents[message] {
			if event.V > 0 {
				xys = append(xys, struct{ X, Y float64 }{tr(event.T), event.V})
				xErrs = append(xErrs, struct{ Low, High float64 }{-event.V / 2, event.V / 2})
			}
		}
		s, err := plotter.NewScatter(xys)
		say.ExitIfError("Couldn't create scatter plot", err)

		s.GlyphStyle = plot.GlyphStyle{
			Color:  viz.OrderedColor(i),
			Radius: 2,
			Shape:  plot.CircleGlyph{},
		}

		p.Add(s)

		xErrsPlot, err := plotter.NewXErrorBars(struct {
			plotter.XYer
			plotter.XErrorer
		}{xys, xErrs})
		say.ExitIfError("Couldn't create x errors plot", err)
		xErrsPlot.LineStyle = viz.LineStyle(viz.OrderedColor(i), 1)
		p.Add(xErrsPlot)

	}

	for _, event := range significantEvents["garden-linux.garden-server.create.creating"] {
		l := viz.NewVerticalLine(tr(event.T))
		l.LineStyle = viz.LineStyle(viz.Blue, 1)
		l.LineStyle.Dashes = viz.Dot
		p.Add(l)
	}

	for _, event := range significantEvents["garden-linux.garden-server.destroy.destroying"] {
		l := viz.NewVerticalLine(tr(event.T))
		l.LineStyle = viz.LineStyle(viz.Red, 1)
		l.LineStyle.Dashes = viz.Dot
		p.Add(l)
	}

	p.X.Min = 0
	p.X.Max = 1800
	p.Y.Min = 1e-5
	p.Y.Max = 1e2

	board.AddNextSubPlot(p)
	board.Save(float64(nAcross)*16.0, float64(nHigh)*4.0, config.DataDir("garden-dt", "garden-dt-scatter.svg"))
}

func generateManyGardenDTScatter(messages []string, significantEvents map[string]Events) {
	tr := func(t time.Time) float64 {
		return float64(t.UnixNano())/1e9 - 1.445379815e9
	}

	nAcross := 1
	nHigh := len(messages) - 3

	board := viz.NewUniformBoard(nAcross, nHigh, 0.01)
	for i, message := range messages {
		if message == "garden-linux.container.info-starting" {
			continue
		}

		p, _ := plot.New()
		p.Title.Text = message
		p.X.Label.Text = "Time (s)"
		p.Y.Label.Text = "Event Duration (s)"
		p.Y.Scale = plot.LogScale
		p.Y.Tick.Marker = plot.LogTicks

		xys := plotter.XYs{}
		xErrs := plotter.XErrors{}
		for _, event := range significantEvents[message] {
			if event.V > 0 {
				xys = append(xys, struct{ X, Y float64 }{tr(event.T), event.V})
				xErrs = append(xErrs, struct{ Low, High float64 }{-event.V / 2, event.V / 2})
			}
		}
		if len(xys) == 0 {
			say.Println(0, "No data for %s", message)
			continue
		}
		s, err := plotter.NewScatter(xys)
		say.ExitIfError("Couldn't create scatter plot", err)

		s.GlyphStyle = plot.GlyphStyle{
			Color:  viz.OrderedColor(i),
			Radius: 2,
			Shape:  plot.CircleGlyph{},
		}

		p.Add(s)

		xErrsPlot, err := plotter.NewXErrorBars(struct {
			plotter.XYer
			plotter.XErrorer
		}{xys, xErrs})
		say.ExitIfError("Couldn't create x errors plot", err)
		xErrsPlot.LineStyle = viz.LineStyle(viz.OrderedColor(i), 1)
		p.Add(xErrsPlot)

		for _, event := range significantEvents["garden-linux.garden-server.create.creating"] {
			l := viz.NewVerticalLine(tr(event.T))
			l.LineStyle = viz.LineStyle(viz.Blue, 1)
			l.LineStyle.Dashes = viz.Dot
			p.Add(l)
		}

		for _, event := range significantEvents["garden-linux.garden-server.destroy.destroying"] {
			l := viz.NewVerticalLine(tr(event.T))
			l.LineStyle = viz.LineStyle(viz.Red, 1)
			l.LineStyle.Dashes = viz.Dot
			p.Add(l)
		}

		p.X.Min = 0
		p.X.Max = 1800
		p.Y.Min = 1e-4
		p.Y.Max = 1e2

		board.AddNextSubPlot(p)
	}

	board.Save(float64(nAcross)*16.0, float64(nHigh)*4.0, config.DataDir("garden-dt", "many-garden-dt-scatter.svg"))
}
