.PHONY: all build build-all clean frontend-build

BINARY_CLI := ok
GO := go
CURRENT_DIR := $(shell pwd)
NPM := npm

PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

HOST_OS := $(shell go env GOOS)
HOST_ARCH := $(shell go env GOARCH)

all: build-all

frontend-build:
	@echo "Building frontend..."
	@cd frontend && $(NPM) run build

build: 
	@echo "Building for host platform ($(HOST_OS)/$(HOST_ARCH))..."
	@TEMP_DIR=$$(mktemp -d); \
	cp -R backend/lib "$$TEMP_DIR/"; \
	cp -R backend/models "$$TEMP_DIR/"; \
	cp -R backend/components "$$TEMP_DIR/"; \
	cp backend/go.mod "$$TEMP_DIR/"; \
	cp backend/main.go "$$TEMP_DIR/"; \
	mkdir -p "$$TEMP_DIR/lib/server/static"; \
	cp -R frontend/dist/* "$$TEMP_DIR/lib/server/static"; \
	cd "$$TEMP_DIR" && $(GO) mod tidy && \
	$(GO) build -o $(BINARY_CLI)_$(HOST_OS)_$(HOST_ARCH) .; \
	BUILD_EXIT=$$?; \
	if [ $$BUILD_EXIT -ne 0 ]; then \
		echo "Build failed"; \
		rm -rf "$$TEMP_DIR"; \
		exit $$BUILD_EXIT; \
	fi; \
	mkdir -p "$$TEMP_DIR/bin"; \
	mv "$$TEMP_DIR/$(BINARY_CLI)_$(HOST_OS)_$(HOST_ARCH)" "$$TEMP_DIR/bin/"; \
	mkdir -p "$(CURRENT_DIR)/bin"; \
	cp -R "$$TEMP_DIR/bin/"* "$(CURRENT_DIR)/bin/"; \
	rm -rf "$$TEMP_DIR"; \
	echo "Binaries built at:"
	@echo "  ./bin/$(BINARY_CLI)_$(HOST_OS)_$(HOST_ARCH)"

build-all: frontend-build
	@echo "Building for all platforms..."
	@TEMP_DIR=$$(mktemp -d); \
	cp -R backend/lib "$$TEMP_DIR/"; \
	cp -R backend/models "$$TEMP_DIR/"; \
	cp -R backend/components "$$TEMP_DIR/"; \
	cp backend/go.mod "$$TEMP_DIR/"; \
	cp backend/main.go "$$TEMP_DIR/"; \
	mkdir -p "$$TEMP_DIR/lib/server/static"; \
	cp -R frontend/dist/* "$$TEMP_DIR/lib/server/static"; \
	cd "$$TEMP_DIR" && $(GO) mod tidy; \
	mkdir -p "$$TEMP_DIR/bin"; \
	mkdir -p "$(CURRENT_DIR)/bin"; \
	for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		echo "Building for $$os/$$arch..."; \
		cd "$$TEMP_DIR" && \
		GOOS=$$os GOARCH=$$arch $(GO) build -o bin/$(BINARY_CLI)_$${os}_$${arch} .; \
		if [ $$? -ne 0 ]; then \
			echo "⨯ Build for $$os/$$arch failed"; \
			continue; \
		fi; \
		echo "✓ Built for $$os/$$arch"; \
	done; \
	cp -R "$$TEMP_DIR/bin/"* "$(CURRENT_DIR)/bin/"; \
	rm -rf "$$TEMP_DIR"; \
	echo "All builds completed. Binaries are in ./bin/"
	@echo "Built binaries:"
	@ls -la bin/ | grep -v '^d' | tail -n +2

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f $(BINARY_CLI) 
	rm -f $(BINARY_CLI)_*
	rm -rf frontend/dist/