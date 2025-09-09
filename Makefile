BUILD_DIR := build
PKG 			:= ./cmd/hora

LINUX_ARCHS := amd64 arm64 386
DARWIN_ARCHS := amd64 arm64

build:
	go build -o build/hora ./cmd/hora

build-all: build-darwin build-linux

build-darwin:
	@for arch in $(DARWIN_ARCHS); do \
		echo "Building for darwin $$arch..."; \
		GOOS=darwin GOARCH=$$arch CGO_ENABLED=1 go build -o $(BUILD_DIR)/hora-darwin-$$arch $(PKG); \
	done

build-linux:
	@for arch in $(LINUX_ARCHS); do \
		echo "Building for linux $$arch..."; \
		GOOS=linux GOARCH=$$arch CGO_ENABLED=0 go build -o $(BUILD_DIR)/hora-linux-$$arch $(PKG); \
	done

clean:
	rm -f build/hora

test:
	go test ./...

install-dependencies:
	go install ./cmd/hora

run:
	go run ./cmd/hora/...

clean-all:
	rm -f build/hora build/hora-*

.PHONY: build build-all build-darwin build-linux clean test install-dependencies run
