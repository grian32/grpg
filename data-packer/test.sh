#!/bin/bash
go build -o grpgpack
./grpgpack tex -m testdata/test_manifest.toml -o testdata/textures.grpgtex
