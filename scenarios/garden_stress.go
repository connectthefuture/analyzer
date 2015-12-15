package scenarios

import (
	"flag"
	"io/ioutil"
	"reflect"

	"code.google.com/p/plotinum/plot"

	. "github.com/onsi/analyzer/dsl"

	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/viz"
	"github.com/pivotal-golang/lager/chug"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"

	"github.com/onsi/say"
)

func GenerateGardenStressTestsCommand() say.Command {
	return say.Command{
		Name:        "garden-stress-tests",
		Description: "10/30/2015: are we leaking containers?",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeGardenStressTests()
		},
	}
}

func analyzeGardenStressTests() {
	data, err := ioutil.ReadFile(config.DataDir("garden-stress-test", "all.unified"))
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	significantEvents := analyzers.ExtractSignificantEvents(entries)

	allow := map[string]bool{
		"garden-linux.garden-server.create.creating":        true,
		"rep.depot-client.delete-container.destroy.started": true,
	}

	filteredSignificantEvents := analyzers.SignificantEvents{}
	filteredSignificantEvents.LogWithThreshold(0.2)
	for name, events := range significantEvents {
		if allow[name] {
			filteredSignificantEvents[name] = events
		}
	}
	filteredSignificantEvents.LogWithThreshold(0.2)

	limit := viz.NewHorizontalLine(256)
	limit.LineStyle = viz.LineStyle(viz.Red, 1)

	options := analyzers.SignificantEventsOptions{
		LineOverlays:    gardenStressTestContainerCountOverlays(entries),
		VerticalMarkers: gardenStressTestFailedContainerCreates(entries),
		OverlayPlots:    []plot.Plotter{limit},
		WidthStretch:    3,
		MaxY:            400,
	}

	analyzers.VisualizeSignificantEvents(
		filteredSignificantEvents,
		config.DataDir("garden-stress-test", "out.png"),
		options,
	)
}

func gardenStressTestFailedContainerCreates(entries []chug.LogEntry) []analyzers.VerticalMarker {
	markers := []analyzers.VerticalMarker{}
	for _, entry := range entries {
		if entry.Message == "rep.depot-client.run-container.failed-creating-container-in-garden" {
			markers = append(markers, analyzers.VerticalMarker{T: entry.Timestamp, LineStyle: viz.LineStyle(viz.Red, 1, viz.Dot)})
		}
	}
	return markers
}

func gardenStressTestContainerCountOverlays(entries []chug.LogEntry) []analyzers.LineOverlay {
	lrps := Events{}
	tasks := Events{}
	for _, entry := range entries {
		if entry.Message == "garden-linux.garden-server.bulk_info.got-bulkinfo" {
			nTasks := 0.0
			nLRPS := 0.0

			handles := reflect.ValueOf(entry.Data["handles"])
			for i := 0; i < handles.Len(); i += 1 {
				handle := handles.Index(i).Interface().(string)
				if len(handle) == 110 {
					nLRPS += 1
				}
				if len(handle) == 69 {
					nTasks += 1
				}
			}
			lrps = append(lrps, Event{T: entry.Timestamp, V: nLRPS})
			tasks = append(tasks, Event{T: entry.Timestamp, V: nTasks})
		}
	}

	return []analyzers.LineOverlay{
		{lrps, viz.LineStyle(viz.Blue, 2)},
		{tasks, viz.LineStyle(viz.Red, 2)},
	}
}
