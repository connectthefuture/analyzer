package main

import (
	"github.com/onsi/analyzer/scenarios"
	"github.com/onsi/say"
)

func main() {
	executable := say.Executable{
		Name:          "analyzer",
		Description:   "Analyzer analyzes data",
		CommandGroups: scenarios.GenerateCommandGroups(),
	}
	say.Invoke(executable)
}
