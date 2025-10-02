CONTAINER_CMD := $(shell \
    if command -v docker &>/dev/null && docker compose version &>/dev/null; then \
        echo "docker compose"; \
    elif command -v podman &>/dev/null && podman compose version &>/dev/null; then \
        echo "podman compose"; \
    elif command -v docker-compose &>/dev/null; then \
        echo "docker-compose"; \
    elif command -v podman-compose &>/dev/null; then \
        echo "podman-compose"; \
    else \
        echo "docker compose"; \
    fi \
)
BINARY_NAME=ok

GO := go
CURRENT_DIR := $(shell pwd)
NPM := npm

PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

HOST_OS := $(shell go env GOOS)
HOST_ARCH := $(shell go env GOARCH)

container-start:
	$(CONTAINER_CMD) -f docker-compose.dev.yml up --build -d

container-stop:
	$(CONTAINER_CMD) -f docker-compose.dev.yml down
	make clean

container-restart:
	$(CONTAINER_CMD) -f docker-compose.dev.yml down
	$(CONTAINER_CMD) -f docker-compose.dev.yml up --build -d

container-logs:
	$(CONTAINER_CMD) -f docker-compose.dev.yml logs -f

container-exec:
	$(CONTAINER_CMD) -f docker-compose.dev.yml exec app bash

clean:
	@echo "Cleaning up..."
	rm -rf bin
	rm -rf go.sum
	rm -rf frontend/dist/
	rm -rf frontend/node_modules
	rm -rf ~/.ok

frontend-build:
	@echo "Building frontend..."
	@cd frontend && rm -rf dist
	@cd frontend && $(NPM) i && $(NPM) run build
	@echo "Frontend build complete."

	@echo "Copying frontend build to backend..."
	@mkdir -p ~/.ok/frontend
	@cp -r frontend/dist/* ~/.ok/frontend/
	@echo "Frontend build copied to ~/.ok/frontend"

	@echo "Frontend build complete."

setup:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

	@echo "Creating config directory..."
	@mkdir -p ~/.ok
	@echo "{}" > ~/.ok/.ok_config

build: setup
	@echo "Building..."
	go build -o bin/$(BINARY_NAME)

	@echo "Building for host platform ($(HOST_OS)/$(HOST_ARCH))..."
	$(GO) build -o bin/$(BINARY_NAME)_$(HOST_OS)_$(HOST_ARCH) .;

	echo "Binaries built at:"
	@echo "  ./bin/$(BINARY_NAME)_$(HOST_OS)_$(HOST_ARCH)"

install: build frontend-build
	@echo "Installing..."
	cp bin/$(BINARY_NAME)_$(HOST_OS)_$(HOST_ARCH) /usr/local/bin/${BINARY_NAME}

install-sudo: build frontend-build
	@echo "Installing with sudo..."
	sudo cp bin/$(BINARY_NAME)_$(HOST_OS)_$(HOST_ARCH) /usr/local/bin/${BINARY_NAME}

# Test targets
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-report:
	@echo "Generating coverage report..."
	$(GO) tool cover -func=coverage.out | grep total

test-clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html

test: test-clean test-coverage
	@echo "All tests completed successfully!"

dev-run: setup build install

.PHONY: dev clean setup build install test test-coverage test-coverage-report test-clean
