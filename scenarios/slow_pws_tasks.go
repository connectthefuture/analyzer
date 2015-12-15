package scenarios

import (
	"flag"
	"io/ioutil"
	"reflect"
	"strings"
	"time"

	. "github.com/onsi/analyzer/dsl"
	"github.com/onsi/gomega/format"

	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/viz"
	"github.com/pivotal-golang/lager/chug"

	"github.com/onsi/analyzer/analyzers"
	"github.com/onsi/analyzer/util"

	"github.com/onsi/say"
)

type SlowPWSTaskRun struct {
	Name           string
	CliffTimestamp time.Time
	EndTimestamp   time.Time
}

var slowPWSTaskRuns = []SlowPWSTaskRun{
	// {
	// 	Name: "cell_aufs_z1.0",
	// },
	// {
	// 	Name: "cell_aufs_z1.1",
	// },
	// {
	// 	Name: "cell_aufs_z1.2",
	// },
	{
		Name:           "cell_z1.16.0",
		CliffTimestamp: time.Date(2015, time.October, 28, 17, 25, 20, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 28, 19, 20, 0, 0, time.Local),
	},
	{
		Name:           "cell_z1.19.0",
		CliffTimestamp: time.Date(2015, time.October, 28, 21, 17, 30, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 28, 23, 57, 0, 0, time.Local),
	},
	{
		Name:           "cell_z2.1.0",
		CliffTimestamp: time.Date(2015, time.October, 29, 0, 42, 50, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 29, 1, 04, 50, 0, time.Local),
	},
	{
		Name:           "cell_z2.10.0",
		CliffTimestamp: time.Date(2015, time.October, 29, 6, 39, 40, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 29, 7, 8, 40, 0, time.Local),
	},
	{
		Name:           "cell_z2.21.0",
		CliffTimestamp: time.Date(2015, time.October, 29, 7, 16, 20, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 29, 7, 41, 10, 0, time.Local),
	},
	{
		Name:           "cell_z2.19.0",
		CliffTimestamp: time.Date(2015, time.October, 29, 8, 11, 40, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 29, 8, 32, 10, 0, time.Local),
	},
	{
		Name:           "cell_z2.1.1",
		CliffTimestamp: time.Date(2015, time.October, 29, 8, 36, 10, 0, time.Local),
		EndTimestamp:   time.Date(2015, time.October, 29, 8, 59, 30, 0, time.Local),
	},
}

func GenerateSlowPWSTasksCommand() say.Command {
	return say.Command{
		Name:        "slow-pws-tasks",
		Description: "10/28/2015: slow tasks on PWS due to overwhelmed cell-z1/10 nailed by auction",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			analyzeSlowPWSTasks()
		},
	}
}

func analyzeSlowPWSTasks() {
	for _, run := range slowPWSTaskRuns {
		say.Println(0, say.Green(run.Name))
		data, err := ioutil.ReadFile(config.DataDir("pws-slow-tasks", run.Name+".unified"))
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
			LineOverlays: containerCountOverlays(entries),
		}
		if !run.CliffTimestamp.IsZero() {
			options.MaxT = run.EndTimestamp.Add(time.Minute * 30)
			options.VerticalMarkers = []analyzers.VerticalMarker{
				{T: run.CliffTimestamp, LineStyle: viz.LineStyle(viz.Red, 1, viz.Dash)},
				{T: run.EndTimestamp, LineStyle: viz.LineStyle(viz.Black, 1, viz.Dash)},
			}
		}

		analyzers.VisualizeSignificantEvents(
			filteredSignificantEvents,
			config.DataDir("pws-slow-tasks", run.Name+".png"),
			options,
		)

	}
}

func containerCountOverlays(entries []chug.LogEntry) []analyzers.LineOverlay {
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

func findTheApp() {
	allLrps := map[string][]string{}

	for _, run := range slowPWSTaskRuns {
		say.Println(0, say.Green(run.Name))
		lrps := map[string]bool{}

		data, err := ioutil.ReadFile(config.DataDir("pws-slow-tasks", run.Name+".unified"))
		say.ExitIfError("couldn't read log file", err)

		entries := util.ChugLagerEntries(data)
		tMin := run.CliffTimestamp.Add(-2 * time.Minute)
		tCliff := run.CliffTimestamp

		for _, entry := range entries {
			if entry.Timestamp.After(tMin) && entry.Timestamp.Before(tCliff) && entry.Message == "garden-linux.garden-server.bulk_info.got-bulkinfo" {
				handles := reflect.ValueOf(entry.Data["handles"])
				for i := 0; i < handles.Len(); i += 1 {
					handle := handles.Index(i).Interface().(string)
					if len(handle) == 110 {
						guid := handle[0:36]
						lrps[guid] = true
					}
				}
			}
		}
		say.Println(0, format.Object(lrps, 0))

		for lrp := range lrps {
			allLrps[lrp] = append(allLrps[lrp], run.Name)
		}
	}
	say.Println(0, say.Green("Counts"))
	for lrp, runs := range allLrps {
		if len(runs) > 1 {
			say.Println(0, "%s: %s", lrp, say.Green("%s", strings.Join(runs, ", ")))
		}
	}
}
