package scenarios

import (
	"flag"
	"io/ioutil"
	"reflect"

	. "github.com/onsi/analyzer/dsl"

	"code.cloudfoundry.org/lager/chug"
	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/viz"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"

	"github.com/onsi/say"
)

func GenerateCPUWeightStressTestCommand() say.Command {
	return say.Command{
		Name:        "cpu-weight-stress-test",
		Description: "10/30/2015: cpu-weight stress tests",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeCPUWeightStresstest()
		},
	}
}

func analyzeCPUWeightStresstest() {
	runs := []string{
		"unmodified-run",
		"aufs-run",
		"2-conc-run",
	}
	for _, run := range runs {
		say.Println(0, say.Green(run))
		data, err := ioutil.ReadFile(config.DataDir("cpu-wait-stress-test", run+".unified"))
		say.ExitIfError("couldn't read log file", err)

		entries := util.ChugLagerEntries(data)

		significantEvents := analyzers.ExtractSignificantEvents(entries)

		allow := map[string]bool{
			"rep.auction-fetch-state.handling":                                             true,
			"rep.container-metrics-reporter.tick.started":                                  true,
			"rep.depot-client.run-container.creating-container-in-garden":                  true,
			"rep.depot-client.delete-container.destroy.started":                            true,
			"rep.depot-client.run-container.run.action.download-step.fetch-starting":       true,
			"rep.depot-client.run-container.run.monitor-run.run-step.running":              true,
			"rep.depot-client.run-container.run.run-step-process.step-finished-with-error": true,
			"rep.depot-client.run-container.run.setup.download-step.fetch-starting":        true,
		}

		filteredSignificantEvents := analyzers.SignificantEvents{}
		for name, events := range significantEvents {
			if allow[name] {
				filteredSignificantEvents[name] = events
			}
		}
		filteredSignificantEvents.LogWithThreshold(0.2)

		options := analyzers.SignificantEventsOptions{
			LineOverlays: cpuWeightStressTestContainerCountOverlays(entries),
		}

		analyzers.VisualizeSignificantEvents(
			filteredSignificantEvents,
			config.DataDir("cpu-wait-stress-test", run+".png"),
			options,
		)

	}
}

func cpuWeightStressTestContainerCountOverlays(entries []chug.LogEntry) []analyzers.LineOverlay {
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
