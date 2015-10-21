package scenarios

import (
	"flag"
	"io/ioutil"

	"code.google.com/p/plotinum/plot"

	"github.com/onsi/analyzer/config"
	. "github.com/onsi/analyzer/dsl"
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
	fetchStateDurationEvents := ExtractEventsFromLagerData(entries, DurationExtractor("duration"))

	allDurations := fetchStateDurationEvents.Data()
	highDurations := allDurations.Filter(NumFilter(">", 0.2))
	lowDurations := allDurations.Filter(NumFilter("<=", 0.2))

	earlyAllDurations := fetchStateDurationEvents.FilterTimes(NumFilter("<", 1445274027070991039)).Data()
	earlyHighDurations := earlyAllDurations.Filter(NumFilter(">", 0.2))
	earlyLowDurations := earlyAllDurations.Filter(NumFilter("<=", 0.2))

	lateAllDurations := fetchStateDurationEvents.FilterTimes(NumFilter(">=", 1445274027070991039)).Data()
	lateHighDurations := lateAllDurations.Filter(NumFilter(">", 0.2))
	lateLowDurations := lateAllDurations.Filter(NumFilter("<=", 0.2))

	board := viz.NewUniformBoard(3, 1, 0)
	earlyScale := float64(len(allDurations)) / float64(len(earlyAllDurations))
	lateScale := float64(len(allDurations)) / float64(len(lateAllDurations))

	say.Println(0, "All Fetches:   %s", allDurations.Stats())
	say.Println(1, ">  0.2s:     %s", highDurations.Stats())
	say.Println(1, "<= 0.2s:     %s", lowDurations.Stats())
	say.Println(0, "Early Fetches: %s", earlyAllDurations.Stats())
	say.Println(1, ">  0.2s:     %s", earlyHighDurations.Stats())
	say.Println(1, "<= 0.2s:     %s", earlyLowDurations.Stats())
	say.Println(0, "Late Fetches:  %s", lateAllDurations.Stats())
	say.Println(1, ">  0.2s:     %s", lateHighDurations.Stats())
	say.Println(1, "<= 0.2s:     %s", lateLowDurations.Stats())

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

	board.Save(12, 4, config.DataDir("auctioneer-fetch-state-duration", "auctioneer-fetch-state-duration.svg"))
}
