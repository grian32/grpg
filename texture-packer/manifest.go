package main

import (
	"log"
	"os"
	"strings"
)

type GRPGTexManifestEntry struct {
	InternalName string
	FilePath     string
}

func ParseManifestFile(path string) []GRPGTexManifestEntry {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	var lines = strings.Split(string(content), "\n")

	entries := make([]GRPGTexManifestEntry, len(lines))

	for idx, line := range lines {
		var contents = strings.Split(line, "=")

		entries[idx] = GRPGTexManifestEntry{
			InternalName: contents[0],
			FilePath:     contents[1],
		}
	}

	return entries
}
