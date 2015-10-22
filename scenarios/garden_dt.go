package scenarios

import (
	"flag"
	"io/ioutil"

	"code.google.com/p/plotinum/plot"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
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

	significantEvents := analyzers.ExtractSignificantEvents(entries)
	significantEvents.LogWithThreshold(0.2)

	delete(significantEvents, "garden-linux.container.info-starting")

	markedEvents := map[string]plot.LineStyle{
		"garden-linux.garden-server.create.creating":    viz.LineStyle(viz.Blue, 1, viz.Dot),
		"garden-linux.garden-server.destroy.destroying": viz.LineStyle(viz.Red, 1, viz.Dot),
	}

	analyzers.VisualizeSignificantEvents(
		significantEvents,
		markedEvents,
		config.DataDir("garden-dt", "many-garden-dt.svg"),
	)
}
