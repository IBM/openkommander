CONTAINER_CMD ?= podman compose

dev:
	$(CONTAINER_CMD) -f docker-compose.dev.yml up --build -d

clean:
	$(CONTAINER_CMD) -f docker-compose.dev.yml down
	rm -rf bin/*

setup:
	@echo "Downloading dependencies..."
	go mod download
