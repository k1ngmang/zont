.PHONY: build run clean test


BINARY_NAME=zont


build:
	go build -o $(BINARY_NAME) .


run:
	go run .


clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go build -o $(BINARY_NAME) .
	@echo "Successfully!"

deps:
	go mod tidy


fmt:
	go fmt .

vet:
	go vet .
