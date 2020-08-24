BIN_NAME := sonarhawk
BIN_PATH := bin/$(BIN_NAME)
BUILD_FLAGS := -installsuffix "static"

.PHONY: all build clean deps

all: linux

linux:
	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=amd64 \
	go build \
	$(BUILD_FLAGS) \
	-o $(BIN_PATH)-survey \
	$$(find cmd/survey/*.go)

	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=amd64 \
	go build \
	$(BUILD_FLAGS) \
	-o $(BIN_PATH)-process \
	$$(find cmd/process/*.go)

windows:
	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=amd64 \
	go build \
	$(BUILD_FLAGS) \
	-o $(BIN_PATH)-survey.exe \
	cmd/survey/main.go

	CGO_ENABLED=1 \
	GO111MODULE=on \
	GOARCH=amd64 \
	go build \
	$(BUILD_FLAGS) \
	-o $(BIN_PATH)-process.exe \
	cmd/process/main.go

deps:
	@echo Downloading go.mod dependencies && \
		go mod download
