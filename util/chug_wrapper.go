package util

import (
	"bytes"

	"github.com/pivotal-golang/lager/chug"
)

func ChugLagerEntries(raw []byte) []chug.LogEntry {
	buf := bytes.NewBuffer(raw)
	out := make(chan chug.Entry)
	go chug.Chug(buf, out)
	entries := []chug.LogEntry{}
	for entry := range out {
		if entry.IsLager {
			entries = append(entries, entry.Log)
		}
	}
	return entries
}
