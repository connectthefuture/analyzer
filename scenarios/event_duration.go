package scenarios

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gonum/plot/vg/draw"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
)

func GenerateEventDurationCommand() say.Command {
	var minT, maxT float64
	var skipList string
	var blueEvent, redEvent string
	var outFile string
	var significantThreshold int

	var fs = &flag.FlagSet{}
	fs.Float64Var(&minT, "tmin", 0, "Min time")
	fs.Float64Var(&maxT, "tmax", 0, "Max time")
	fs.StringVar(&skipList, "skip", "", "Events to skip (comma delimited)")
	fs.StringVar(&blueEvent, "blue", "", "Events to use to generate blue markers")
	fs.StringVar(&redEvent, "red", "", "Events to use to generate red markers")
	fs.IntVar(&significantThreshold, "n", 2, "Minimum number of events required to make it onto the plot")
	fs.StringVar(&outFile, "o", "", "Output file")

	return say.Command{
		Name:        "event-duration",
		Description: "lager.log -- generate event duration plots",
		FlagSet:     fs,
		Run: func(args []string) {
			if len(args) != 1 {
				say.Println(0, say.Red("please provide a lager file to read"))
				os.Exit(1)
			}

			markedEvents := map[string]draw.LineStyle{}
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

			if outFile == "" {
				outFile = "out.png"
			}

			analyzeEventDurations(args[0], options, significantThreshold, skips, outFile)
		},
	}
}

func analyzeEventDurations(path string, options analyzers.SignificantEventsOptions, n int, skips []string, outFile string) {
	data, err := ioutil.ReadFile(path)
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	significantEvents := analyzers.ExtractSignificantEventsWithThreshold(entries, n)
	significantEvents.LogWithThreshold(0.2)

	for _, skip := range skips {
		say.Println(0, "Skipping %s", skip)
		delete(significantEvents, skip)
	}

	analyzers.VisualizeSignificantEvents(
		significantEvents,
		outFile,
		options,
	)
}
