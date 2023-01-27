build_bobo:
	@echo "Building bobo cli at ./cmd/cli/bobo/bin ..."
	@go build -o ./cmd/cli/bobo/bin/bobo ./cmd/cli/bobo
	@echo "Bobo has been built!"

install_bobo:
	@echo "Installing bobo in ${HOME}/go/bin ..."
	@go build -o ${HOME}/go/bin/bobo ./cmd/cli/bobo
	@echo "Bobo has been installed!"
