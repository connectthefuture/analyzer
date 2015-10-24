package scenarios

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"code.google.com/p/plotinum/plot"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
)

func GenerateEventDurationCommand() say.Command {
	var minT, maxT float64
	var skipList string
	var blueEvent string
	var redEvent string

	var fs = &flag.FlagSet{}
	fs.Float64Var(&minT, "tmin", 0, "Min time")
	fs.Float64Var(&maxT, "tmax", 0, "Max time")
	fs.StringVar(&skipList, "skip", "", "Events to skip (comma delimited)")
	fs.StringVar(&blueEvent, "blue", "", "Events to use to generate blue markers")
	fs.StringVar(&redEvent, "red", "", "Events to use to generate red markers")

	return say.Command{
		Name:        "event-duration",
		Description: "lager.log -- generate event duration plots",
		FlagSet:     fs,
		Run: func(args []string) {
			if len(args) != 1 {
				say.Println(0, say.Red("please provide a lager file to read"))
				os.Exit(1)
			}

			markedEvents := map[string]plot.LineStyle{}
			if blueEvent != "" {
				markedEvents[blueEvent] = viz.LineStyle(viz.Blue, 1, viz.Dot)
			}
			if redEvent != "" {
				markedEvents[redEvent] = viz.LineStyle(viz.Red, 1, viz.Dot)
			}

			options := analyzers.SignificantEventsOptions{
				MinX:         minT,
				MaxX:         maxT,
				MarkedEvents: markedEvents,
			}

			skips := strings.Split(skipList, ",")

			analyzeEventDurations(args[0], options, skips)
		},
	}
}

func analyzeEventDurations(path string, options analyzers.SignificantEventsOptions, skips []string) {
	data, err := ioutil.ReadFile(path)
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	significantEvents := analyzers.ExtractSignificantEvents(entries)
	significantEvents.LogWithThreshold(0.2)

	for _, skip := range skips {
		say.Println(0, "Skipping %s", skip)
		delete(significantEvents, skip)
	}

	analyzers.VisualizeSignificantEvents(
		significantEvents,
		"out.svg",
		options,
	)
}
