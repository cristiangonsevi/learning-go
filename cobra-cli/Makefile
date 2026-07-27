BIN_FILENAME ?= "cli"


build:
	@echo "Building app"
	@mkdir -p ./bin
	@go build -o ./bin/$(BIN_FILENAME) .

run: build
	@./bin/$(BIN_FILENAME)
