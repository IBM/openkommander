CONTAINER_CMD ?= podman compose
BINARY_NAME=ok

dev:
	$(CONTAINER_CMD) -f docker-compose.dev.yml up --build -d

dev-run: build install

clean:
	$(CONTAINER_CMD) -f docker-compose.dev.yml down
	rm -rf bin/*

setup:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME)

install: build
	@echo "Installing..."
	cp bin/$(BINARY_NAME) /usr/local/bin/

.PHONY: dev clean setup build install
