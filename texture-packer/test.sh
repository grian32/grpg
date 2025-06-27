#!/bin/bash
go build -o texpack
cd testdata || exit
../texpack -m test_manifest.txt --texv 2
