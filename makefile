# Variables
APP_NAME := goks
SRC_DIR := ./cmd/
BUILD_DIR := build
GO_FILES := $(shell find ./cmd/ -name '*.go')
VERSION := $(shell git describe --tags --always 2>/dev/null || echo "v0.1.0")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"


# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build: $(BUILD_DIR)/$(APP_NAME)

.PHONY: build-linux
build-linux:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux $(SRC_DIR)

$(BUILD_DIR)/$(APP_NAME): $(GO_FILES)
	mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

# Run the application
.PHONY: run
run:
	go run .

# Clean build files
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Test the application
.PHONY: test
test:
	go test ./...

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	go mod tidy

# Show the version
.PHONY: version
version:
	@echo $(VERSION)
