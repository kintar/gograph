BINARY=gograph
BUILD_DIR=bin
CMD=./cmd/gograph

.PHONY: build test run-build clean

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) $(CMD)
	@echo "Built $(BUILD_DIR)/$(BINARY)"

test:
	go test ./...

run-build:
	go run $(CMD) build .

clean:
	rm -rf $(BUILD_DIR)
