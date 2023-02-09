build_bobo:
	@echo "Building bobo cli at ./cmd/cli/bobo/bin ..."
	@go build -o ./bin/bobo ./cmd/cli/bobo
	@echo "Bobo has been built!"

install_bobo:
	@echo "Installing bobo in ${HOME}/go/bin ..."
	@go install ./cmd/cli/bobo
	@echo "Bobo has been installed!"
