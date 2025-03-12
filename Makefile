CONTAINER_CMD ?= podman compose
BINARY_NAME=ok

dev:
	@echo "Stopping any running processes on ports 2181, 9092, and 8080..."
	@sudo lsof -ti :2181 | xargs -r sudo kill -9
	@sudo lsof -ti :9092 | xargs -r sudo kill -9
	@sudo lsof -ti :8080 | xargs -r sudo kill -9

	@echo "Removing any existing containers..."
	@podman ps -a -q | xargs -r podman rm -f
	
	@echo "Starting development environment..."
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
