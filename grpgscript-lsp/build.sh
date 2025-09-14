#!/bin/bash
go build -o grpgscriptlsp ./cmd
cp ./grpgscriptlsp editors/zed/
rm -f grpgscriptlsp