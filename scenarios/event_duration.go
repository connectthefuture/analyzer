package scenarios

import (
	"flag"
	"io/ioutil"
	"os"

	"code.google.com/p/plotinum/plot"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/say"
)

func GenerateEventDurationCommand() say.Command {
	return say.Command{
		Name:        "event-duration",
		Description: "lager.log -- generate event duration plots",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			if len(args) != 1 {
				say.Println(0, say.Red("please provide a lager file to read"))
				os.Exit(1)
			}
			analyzeEventDurations(args[0])
		},
	}
}

func analyzeEventDurations(path string) {
	data, err := ioutil.ReadFile(path)
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	significantEvents := analyzers.ExtractSignificantEvents(entries)
	significantEvents.LogWithThreshold(0.2)

	analyzers.VisualizeSignificantEvents(
		significantEvents,
		map[string]plot.LineStyle{},
		"out.svg",
	)
}
