#!/bin/bash

# build for windows
GOOS=windows GOARCH=amd64 go build -o supporter.exe
zip supporter-win.zip supporter.exe supporterCfg.json.sample
rm supporter.exe
echo "finish build for windows"

# build for mac
GOOS=darwin GOARCH=amd64 go build -o supporter
zip supporter-mac.zip supporter supporterCfg.json.sample
rm supporter
echo "finish build for mac"

