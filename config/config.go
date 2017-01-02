package config

import "path/filepath"

const DATA_DIR = "/Users/onsi/workspace/go/src/github.com/onsi/analyzer/data"

func DataDir(args ...string) string {
	paths := []string{DATA_DIR}
	paths = append(paths, args...)
	return filepath.Join(paths...)
}
