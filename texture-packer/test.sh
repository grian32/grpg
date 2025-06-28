#!/bin/bash
go build -o texpack
./texpack -m testdata/test_manifest.txt -o testdata/textures.pak --texv 1
