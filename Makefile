build_bobo:
	@echo "Building bobo cli at ./cmd/cli/bobo/bin ..."
	@go build -o ./bin/bobo ./cmd/cli/bobo
	@echo "Bobo has been built!"

install_bobo:
	@echo "Installing bobo in ${HOME}/go/bin ..."
	@go install ./cmd/cli/bobo
	@echo "Bobo has been installed!"

pre-compile:
	@echo "Compiling windows-amd64 ..."
	@GOOS=windows GOARCH=amd64 go build -o ./bin/bobo.windows-amd64 ./cmd/cli/bobo
	@echo "Compiling darwin-amd64 ..."
	@GOOS=darwin GOARCH=amd64 go build -o ./bin/bobo.darwin-amd64 ./cmd/cli/bobo
	@echo "Compiling linux-amd64 ..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/bobo.linux-amd64 ./cmd/cli/bobo
	@echo "Running tar -czf on binaries"
	@cd bin && for f in $(ls bin); do tar -czf $f.tar.gz $f; done;
