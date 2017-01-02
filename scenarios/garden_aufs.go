package scenarios

import (
	"flag"
	"io/ioutil"

	"github.com/gonum/plot"

	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/viz"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"

	"github.com/onsi/say"
)

func GenerateGardenAUFSStressTestsCommand() say.Command {
	return say.Command{
		Name:        "garden-aufs-tests",
		Description: "11/1/2015: how is the new AUFS work holding up?",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeGardenAUFSStressTests()
		},
	}
}

func analyzeGardenAUFSStressTests() {
	runs := []string{"diego-2.cell-z1.0"} //, "diego-1.cell-z1.0"}
	for _, run := range runs {
		data, err := ioutil.ReadFile(config.DataDir("garden-aufs-test", run+".unified"))
		say.ExitIfError("couldn't read log file", err)

		entries := util.ChugLagerEntries(data)

		significantEvents := analyzers.ExtractSignificantEvents(entries)
		// significantEvents.LogWithThreshold(0.2)

		allow := map[string]bool{
			"garden-linux.loop-mounter.unmount.failed-to-unmount":                                                        true,
			"garden-linux.garden-server.create.creating":                                                                 true,
			"rep.auction-fetch-state.handling":                                                                           true,
			"rep.container-metrics-reporter.tick.started":                                                                true,
			"rep.depot-client.run-container.creating-container-in-garden":                                                true,
			"rep.depot-client.delete-container.destroy.started":                                                          true,
			"rep.depot-client.run-container.run.action.download-step.fetch-starting":                                     true,
			"rep.depot-client.run-container.run.monitor-run.run-step.running":                                            true,
			"rep.depot-client.run-container.run.run-step-process.step-finished-with-error":                               true,
			"rep.depot-client.run-container.run.setup.download-step.fetch-starting":                                      true,
			"rep.auction-delegate.auction-work.lrp-allocate-instances.requesting-container-allocation":                   true,
			"rep.event-consumer.operation-stream.executing-container-operation.task-processor.fetching-container-result": true,
			"rep.depot-client.run-container.run.action.run-step.running":                                                 true,
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
			WidthStretch:    2,
			MaxY:            400,
		}

		analyzers.VisualizeSignificantEvents(
			filteredSignificantEvents,
			config.DataDir("garden-aufs-test", run+".png"),
			options,
		)
	}
}
