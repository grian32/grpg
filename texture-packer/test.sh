#!/bin/bash
go build -o texpack
./texpack -m test_manifest.txt --texv 1
