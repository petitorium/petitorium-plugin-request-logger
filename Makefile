.DEFAULT_GOAL := build

PLUGIN_NAME := request-logger
BUILD_DIR := bin
GO_VERSION := 1.24

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)-linux-amd64 .

darwin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)-darwin-amd64 .

darwin-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)-darwin-arm64 .

darwin: darwin-amd64 darwin-arm64

windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(PLUGIN_NAME)-windows-amd64.exe .

all: linux darwin windows

build:
	CGO_ENABLED=0 go build -o $(PLUGIN_NAME) .

deps:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(PLUGIN_NAME)
	rm -f $(PLUGIN_NAME).exe

.PHONY: all build clean deps fmt vet test linux darwin windows darwin-amd64 darwin-arm64
