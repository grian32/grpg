#!/bin/bash
rm -f editors/zed/grpgscriptlsp
go build -o grpgscriptlsp ./cmd
cp ./grpgscriptlsp editors/zed/
rm -f grpgscriptlsp