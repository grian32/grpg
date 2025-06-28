package main

import "testing"

func TestParseManifestFile(t *testing.T) {
	expected := []GRPGTexManifestEntry{
		{
			InternalName: "grass",
			FilePath:     "grass_texture.png",
		},
		{
			InternalName: "stone",
			FilePath:     "stone_texture.png",
		},
	}

	filePath := "./testdata/test_manifest.txt"

	output, err := ParseManifestFile(filePath)

	// ehh @ comparison
	if len(output) < 2 || output[0] != expected[0] || output[1] != expected[1] || err != nil {
		t.Errorf("ParseManifestFile = %q, %v, want match for %#q", output, err, expected)
	}
}
