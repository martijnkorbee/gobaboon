build_bobo:
	@echo "Building boboctl at ./bin/boboctl ..."
	@go build -o ./bin/boboctl ./tools/boboctl
	@echo "Bobo has been built!"

install_bobo:
	@echo "Installing bobo in ${HOME}/go/bin ..."
	@go install ./tools/boboctl
	@echo "Bobo has been installed!"

pre-build_bobo:
	@./build-assets.sh
