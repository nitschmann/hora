.PHONY: build clean test install-dependencies run

build:
	go build -o build/hora ./cmd/hora

clean:
	rm -f build/hora

test:
	go test ./...

install-dependencies:
	go install ./cmd/hora

run:
	go run ./cmd/hora/...

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o build/hora-linux-amd64 ./cmd/hora
	GOOS=darwin GOARCH=amd64 go build -o build/hora-darwin-amd64 ./cmd/hora
	GOOS=darwin GOARCH=arm64 go build -o build/hora-darwin-arm64 ./cmd/hora
	GOOS=windows GOARCH=amd64 go build -o build/hora-windows-amd64.exe ./cmd/hora

clean-all:
	rm -f build/hora build/hora-*
