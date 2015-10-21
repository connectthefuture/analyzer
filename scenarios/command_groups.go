package scenarios

import "github.com/onsi/say"

func GenerateCommandGroups() []say.CommandGroup {
	return []say.CommandGroup{
		{
			Name:        "diego",
			Description: "Diego analysis commands",
			Commands: []say.Command{
				GenerateGardenDTCommand(),
				GenerateAuctioneerFetchStateDurationCommand(),
			},
		},
	}
}
