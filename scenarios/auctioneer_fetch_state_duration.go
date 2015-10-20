package scenarios

import (
	"flag"
	"io/ioutil"
	"time"

	"code.google.com/p/plotinum/plot"

	"github.com/onsi/analyzer/config"
	"github.com/onsi/analyzer/util"
	"github.com/onsi/analyzer/viz"
	"github.com/onsi/say"
)

func GenerateAuctioneerFetchStateDurationCommand() say.Command {
	return say.Command{
		Name:        "auctioneer-fetch-state-duration",
		Description: "10/19/2015: auctioneer showing spiky FetchState durations",
		FlagSet:     &flag.FlagSet{},
		Run: func(args []string) {
			generateAuctioneerFetchStateDuration()
		},
	}
}

func generateAuctioneerFetchStateDuration() {
	data, err := ioutil.ReadFile(config.DataDir("auctioneer-fetch-state-duration", "auctioneer-fetch-state-duration.logs"))
	say.ExitIfError("couldn't read log file", err)

	entries := util.ChugLagerEntries(data)

	allDurations := viz.Data{}
	highDurations := viz.Data{}
	lowDurations := viz.Data{}

	earlyAllDurations := viz.Data{}
	earlyHighDurations := viz.Data{}
	earlyLowDurations := viz.Data{}

	lateAllDurations := viz.Data{}
	lateHighDurations := viz.Data{}
	lateLowDurations := viz.Data{}

	for _, entry := range entries {
		duration, err := time.ParseDuration(entry.Data["duration"].(string))
		say.ExitIfError("couldn't parse duration", err)

		allDurations = append(allDurations, duration.Seconds())
		if duration.Seconds() > 0.2 {
			highDurations = append(highDurations, duration.Seconds())
		}
		if duration.Seconds() < 0.2 {
			lowDurations = append(lowDurations, duration.Seconds())
		}

		if entry.Timestamp.UnixNano() < 1445274027070991039 {
			earlyAllDurations = append(earlyAllDurations, duration.Seconds())
			if duration.Seconds() > 0.2 {
				earlyHighDurations = append(earlyHighDurations, duration.Seconds())
			}
			if duration.Seconds() < 0.2 {
				earlyLowDurations = append(earlyLowDurations, duration.Seconds())
			}
		} else {
			lateAllDurations = append(lateAllDurations, duration.Seconds())
			if duration.Seconds() > 0.2 {
				lateHighDurations = append(lateHighDurations, duration.Seconds())
			}
			if duration.Seconds() < 0.2 {
				lateLowDurations = append(lateLowDurations, duration.Seconds())
			}
		}
	}

	board := viz.NewUniformBoard(3, 1, 0)
	earlyScale := float64(len(allDurations)) / float64(len(earlyAllDurations))
	lateScale := float64(len(allDurations)) / float64(len(lateAllDurations))

	p, _ := plot.New()
	p.Title.Text = "All Auctioneer Fetch State Durations"
	p.Add(viz.NewHistogram(allDurations, 20, allDurations.Min(), allDurations.Max()))
	h := viz.NewScaledHistogram(earlyAllDurations, 20, allDurations.Min(), allDurations.Max(), earlyScale)
	h.LineStyle = viz.LineStyle(viz.Blue, 1)
	p.Add(h)
	h = viz.NewScaledHistogram(lateAllDurations, 20, allDurations.Min(), allDurations.Max(), lateScale)
	h.LineStyle = viz.LineStyle(viz.Red, 1)
	p.Add(h)
	board.AddNextSubPlot(p)

	say.Println(0, "All Fetches: %d (early: %d, late: %d)", len(allDurations), len(earlyAllDurations), len(lateAllDurations))

	p, _ = plot.New()
	p.Title.Text = "Auctioneer Fetch State Durations &lt; 0.2s"
	p.Add(viz.NewHistogram(lowDurations, 100, lowDurations.Min(), lowDurations.Max()))
	h = viz.NewScaledHistogram(earlyLowDurations, 100, lowDurations.Min(), lowDurations.Max(), earlyScale)
	h.LineStyle = viz.LineStyle(viz.Blue, 1)
	p.Add(h)
	h = viz.NewScaledHistogram(lateLowDurations, 100, lowDurations.Min(), lowDurations.Max(), lateScale)
	h.LineStyle = viz.LineStyle(viz.Red, 1)
	p.Add(h)
	board.AddNextSubPlot(p)

	say.Println(0, "Fetches that took < 0.2s: %d", len(lowDurations))

	p, _ = plot.New()
	p.Title.Text = "Auctioneer Fetch State Durations > 0.2s"
	p.Add(viz.NewHistogram(highDurations, 40, highDurations.Min(), highDurations.Max()))
	h = viz.NewScaledHistogram(earlyHighDurations, 40, highDurations.Min(), highDurations.Max(), earlyScale)
	h.LineStyle = viz.LineStyle(viz.Blue, 1)
	p.Add(h)
	h = viz.NewScaledHistogram(lateHighDurations, 40, highDurations.Min(), highDurations.Max(), lateScale)
	h.LineStyle = viz.LineStyle(viz.Red, 1)
	p.Add(h)
	board.AddNextSubPlot(p)

	say.Println(0, "Fetches that took > 0.2s: %d", len(highDurations))

	board.Save(12, 4, config.DataDir("auctioneer-fetch-state-duration", "auctioneer-fetch-state-duration.svg"))
}
