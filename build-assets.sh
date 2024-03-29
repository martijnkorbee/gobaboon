##!/bin/bash
#
## Simple script to pre-build binaries to include in github releases.
#
#echo "Compiling windows-amd64 ..."
#GOOS=windows GOARCH=amd64 go build -o ./bin/bobo.windows-amd64 ./cmd/cli/bobo
#
#echo "Compiling darwin-amd64 ..."
#GOOS=darwin GOARCH=amd64 go build -o ./bin/bobo.darwin-amd64 ./cmd/cli/bobo
#
#echo "Compiling linux-amd64 ..."
#GOOS=linux GOARCH=amd64 go build -o ./bin/bobo.linux-amd64 ./cmd/cli/bobo
#
#echo "Running tar -czf for binaries in ./bin"
#mkdir -p ./assets
#for f in  $(cd ./bin && echo *);
#	do tar -czf ./assets/"$f".tar.gz -C ./bin ./"$f";
#done;
