.PHONY: build run clean test deps fmt vet


BINARY_NAME=zont
OUT_DIR=.out


build:
	mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/$(BINARY_NAME) ./cmd

render:
	$(OUT_DIR)/$(BINARY_NAME) render

run:
	$(OUT_DIR)/$(BINARY_NAME)


clean:
	go clean
	rm -rf $(OUT_DIR)

test:
	go test ./...

deps:
	go mod tidy


fmt:
	go fmt ./...

vet:
	go vet ./...
