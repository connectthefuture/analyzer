package scenarios

import "github.com/onsi/say"

func GenerateCommandGroups() []say.CommandGroup {
	return []say.CommandGroup{
		{
			Name:        "log-analysis",
			Description: "Reusable log analysis commands",
			Commands: []say.Command{
				GenerateEventDurationCommand(),
			},
		},
		{
			Name:        "diego",
			Description: "Diego analysis commands",
			Commands: []say.Command{
				GeneratePWSSlowEvacuationCommand(),
				GenerateGardenAUFSStressTestsCommand(),
				GenerateCPUWeightStressTestCommand(),
				GenerateGardenStressTestsCommand(),
				GenerateSlowPWSTasksCommand(),
				GenerateGardenDTCommand(),
				GenerateAuctioneerFetchStateDurationCommand(),
			},
		},
	}
}
