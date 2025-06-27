package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

type GRPGTexManifestEntry struct {
	InternalName string
	FilePath     string
}

func ParseManifestFile(path string) ([]GRPGTexManifestEntry, error) {
	content, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	var lines = strings.Split(string(content), "\n")

	entries := make([]GRPGTexManifestEntry, len(lines))

	for idx, line := range lines {
		var contents = strings.Split(line, "=")

		// eh this is a bit shit but it's an "internal" tool anyway lol
		if !strings.HasSuffix(contents[1], ".png") {
			return nil, errors.New("only .png files are allowed as textures")
		}

		entries[idx] = GRPGTexManifestEntry{
			InternalName: contents[0],
			FilePath:     contents[1],
		}
	}

	return entries, nil
}
