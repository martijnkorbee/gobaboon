baboonctl_build:
	echo "Building baboonctl in bin/ ..."; \
	go build -o ./bin/baboonctl ./tools/baboonctl; \
	echo "baboonctl has been built!"

baboonctl_install:
	echo "Installing baboonctl in ${HOME}/go/bin ..."; \
	go install ./tools/baboonctl; \
	echo "baboonctl has been installed!"

app_build:
	echo "Building app in bin/"; \
	go build -o ./app/bin/app .; \
	echo "Your app is ready!"

app_run: app_build
	echo "Starting application"; \
	cd app; \
	./bin/app

app_start:
	echo "Starting application"; \
	cd app; \
	./bin/app
