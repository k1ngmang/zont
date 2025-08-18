.PHONY: build run clean test deps fmt vet


BINARY_NAME=zont


build:
	go build -o $(BINARY_NAME) ./cmd


run:
	go run ./cmd


clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./...

deps:
	go mod tidy


fmt:
	go fmt ./...

vet:
	go vet ./...
