package scenarios

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/say"
)

func GenerateEventDurationCommand() say.Command {
	var minT, maxT float64

	var fs = &flag.FlagSet{}
	fs.Float64Var(&minT, "tmin", 0, "Min time")
	fs.Float64Var(&maxT, "tmax", 0, "Max time")

	return say.Command{
		Name:        "event-duration",
		Description: "lager.log -- generate event duration plots",
		FlagSet:     fs,
		Run: func(args []string) {
			if len(args) != 1 {
				say.Println(0, say.Red("please provide a lager file to read"))
				os.Exit(1)
			}
			options := analyzers.SignificantEventsOptions{
				MinX: minT,
				MaxX: maxT,
			}
			analyzeEventDurations(args[0], options)
		},
	}
}

func analyzeEventDurations(path string, options analyzers.SignificantEventsOptions) {
	data, err := ioutil.ReadFile(path)
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	significantEvents := analyzers.ExtractSignificantEvents(entries)
	significantEvents.LogWithThreshold(0.2)

	analyzers.VisualizeSignificantEvents(
		significantEvents,
		"out.svg",
		options,
	)
}
