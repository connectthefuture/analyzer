package scenarios

import (
	"flag"
	"io/ioutil"

	"code.cloudfoundry.org/lager/chug"
	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/viz"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"

	"github.com/onsi/say"
)

func GeneratePWSSlowEvacuationCommand() say.Command {
	return say.Command{
		Name:        "pws-slow-evacuation",
		Description: "11/11/2015: with AUFS, container creation time gets real slow with lots of containers and evac takes longer.  Why?",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzePWSSlowEvacuation()
		},
	}
}

func analyzePWSSlowEvacuation() {
	runs := []string{"cell_z1.16", "cell_z1.17", "cell_z2.23", "cell_z2.24", "cell_z2.25"}
	for _, run := range runs {
		data, err := ioutil.ReadFile(config.DataDir("pws-slow-evacuation", run+".unified"))
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

		options := analyzers.SignificantEventsOptions{
			LineOverlays:    gardenStressTestContainerCountOverlays(entries),
			VerticalMarkers: pwsSlowEvacuationFailedProcessOverlays(entries),
			WidthStretch:    2,
			MaxY:            400,
		}

		analyzers.VisualizeSignificantEvents(
			filteredSignificantEvents,
			config.DataDir("pws-slow-evacuation", run+".png"),
			options,
		)
	}
}

func pwsSlowEvacuationFailedProcessOverlays(entries []chug.LogEntry) []analyzers.VerticalMarker {
	markers := []analyzers.VerticalMarker{}
	for _, entry := range entries {
		if entry.Message == "rep.depot-client.run-container.run.run-step-process.step-finished-with-error" {
			markers = append(markers, analyzers.VerticalMarker{T: entry.Timestamp, LineStyle: viz.LineStyle(viz.Red, 1, viz.Dot)})
		}
	}
	return markers
}
