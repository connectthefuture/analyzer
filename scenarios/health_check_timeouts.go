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

func GenerateHealthCheckTimeoutsCommand() say.Command {
	return say.Command{
		Name:        "health-check-timeouts",
		Description: "6/17/2016: Seeing health checks timeout.  Why?",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeHealthCheckTimeouts()
		},
	}
}

func analyzeHealthCheckTimeouts() {
	data, err := ioutil.ReadFile(config.DataDir("health-check-timeouts", "health-check-timeouts.log"))
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	significantEvents := analyzers.ExtractSignificantEvents(entries)
	significantEvents.LogWithThreshold(0.2)

	allow := map[string]bool{
		"rep.auction-fetch-state.handling":            true,
		"rep.container-metrics-reporter.tick.started": true,
		"rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-reserved-container.run-container.containerstore-run.node-run.monitor-run.run-step.running": true,
		"rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-reserved-container.run-container.containerstore-create.starting":                           true,
		"rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-completed-container.deleting-container":                                                    true,
		"rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-reserved-container.run-container.containerstore-run.node-run.action.run-step.running":      true,
		"rep.running-bulker.sync.starting":       true,
		"garden-linux.garden-server.run.spawned": true,
	}

	filteredSignificantEvents := analyzers.SignificantEvents{}

	for name, events := range significantEvents {
		if allow[name] {
			filteredSignificantEvents[name] = events
		}
	}
	filteredSignificantEvents.LogWithThreshold(0.2)

	options := analyzers.SignificantEventsOptions{
		LineOverlays:    gardenStressTestContainerCountOverlays(entries),
		VerticalMarkers: healthCheckTimeoutsFailedHealthMonitor(entries),
		WidthStretch:    2,
		MaxY:            400,
	}

	analyzers.VisualizeSignificantEvents(
		filteredSignificantEvents,
		config.DataDir("health-check-timeouts", "health-check-timeouts.png"),
		options,
	)
}

func healthCheckTimeoutsFailedHealthMonitor(entries []chug.LogEntry) []analyzers.VerticalMarker {
	markers := []analyzers.VerticalMarker{}
	for _, entry := range entries {
		if entry.Message == "rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-reserved-container.run-container.containerstore-create.starting" {
			markers = append(markers, analyzers.VerticalMarker{T: entry.Timestamp, LineStyle: viz.LineStyle(viz.Blue, 1, viz.Dot)})
		}
		if entry.Message == "rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-reserved-container.run-container.containerstore-run.node-run.monitor-run.timeout-step.timed-out" {
			markers = append(markers, analyzers.VerticalMarker{T: entry.Timestamp, LineStyle: viz.LineStyle(viz.Red, 2, viz.Dot)})
		}
		if entry.Message == "rep.event-consumer.operation-stream.executing-container-operation.ordinary-lrp-processor.process-reserved-container.run-container.containerstore-run.node-run.monitor-run.run-step.failed-creating-process" {
			markers = append(markers, analyzers.VerticalMarker{T: entry.Timestamp, LineStyle: viz.LineStyle(viz.Red, 2, viz.Dash)})
		}
	}
	return markers
}
