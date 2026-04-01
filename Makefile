# Variables
BINARY_NAME=pptx-probe
GO_FILES=$(shell find . -name "*.go")
OUTPUT_DIR=output

GOFMT ?= gofmt "-s"

# Default target: build the binary
all: build

## build: Compile the Go binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) main.go

## run: Build and run with a default example
run: build
	./$(BINARY_NAME) example/example1.pptx

## test: Run go tests
test:
	go test -v ./...

## clean: Remove binary and the extracted output folder
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(OUTPUT_DIR)

## fmt: Format all Go files (important for clean code)
fmt:
	$(GOFMT) -w $(GO_FILES)


## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | sed -e 's/## //' | column -t -s ':'

.PHONY: all build run test clean fmt help