build_baboonctl:
	@echo "Building baboonctl at ./bin/baboonctl ..."
	@go build -o ./bin/baboonctl ./tools/baboonctl
	@echo "baboonctl has been built!"

install_baboonctl:
	@echo "Installing baboonctl in ${HOME}/go/bin ..."
	@go install ./tools/baboonctl
	@echo "baboonctl has been installed!"

pre-build_baboonctl:
	@./build-assets.sh
