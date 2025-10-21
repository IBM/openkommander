#!/bin/bash

# OpenKommander Setup Script
# Interactive installation script for local, Docker, or Podman environments

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="openkommander"
BINARY_NAME="ok"
REPO_URL="https://github.com/IBM/openkommander"
MIN_GO_VERSION="1.24"
MIN_NODE_VERSION="18"

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to compare versions
version_gt() {
    test "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1"
}

# Function to check Go version
check_go_version() {
    if ! command_exists go; then
        return 1
    fi
    
    local go_version=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
    if version_gt "$MIN_GO_VERSION" "$go_version"; then
        return 1
    fi
    return 0
}

# Function to check Node.js version
check_node_version() {
    if ! command_exists node; then
        return 1
    fi
    
    local node_version=$(node --version | sed 's/v//')
    if version_gt "$MIN_NODE_VERSION" "$node_version"; then
        return 1
    fi
    return 0
}

# Function to check local prerequisites
check_local_prerequisites() {
    print_info "Checking local prerequisites..."
    
    local missing_deps=()
    
    # Check Go
    if ! check_go_version; then
        missing_deps+=("Go ${MIN_GO_VERSION}+ (https://golang.org/dl/)")
    else
        print_success "Go $(go version | grep -o 'go[0-9]\+\.[0-9]\+\.[0-9]\+') found"
    fi
    
    # Check Node.js
    if ! check_node_version; then
        missing_deps+=("Node.js ${MIN_NODE_VERSION}+ (https://nodejs.org/)")
    else
        print_success "Node.js $(node --version) found"
    fi
    
    # Check npm
    if ! command_exists npm; then
        missing_deps+=("npm (usually comes with Node.js)")
    else
        print_success "npm $(npm --version) found"
    fi
    
    # Check make
    if ! command_exists make; then
        missing_deps+=("make (build-essential on Ubuntu/Debian, Xcode tools on macOS)")
    else
        print_success "make found"
    fi
    
    # Check git
    if ! command_exists git; then
        missing_deps+=("git")
    else
        print_success "git found"
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_error "Missing prerequisites:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        return 1
    fi
    
    print_success "All prerequisites satisfied!"
    return 0
}

# Function to prompt user for installation type
prompt_installation_type() {
    echo >&2
    print_info "Choose installation method:" >&2
    echo "1) Local installation (installs directly on your system)" >&2
    echo "2) Docker (runs in Docker container)" >&2
    echo "3) Podman (runs in Podman container)" >&2
    echo >&2
    
    while true; do
        read -p "Enter your choice (1-3): " choice >&2
        case $choice in
            1)
                echo "local"
                return
                ;;
            2)
                echo "docker"
                return
                ;;
            3)
                echo "podman"
                return
                ;;
            *)
                print_warning "Invalid choice. Please enter 1, 2, or 3." >&2
                ;;
        esac
    done
}

# Function to install locally
install_locally() {
    print_info "Starting local installation..."
    
    # Check prerequisites
    if ! check_local_prerequisites; then
        print_error "Prerequisites not met. Please install missing dependencies and try again."
        exit 1
    fi
    
    # Confirm prerequisites with user
    echo
    read -p "All prerequisites are satisfied. Continue with local installation? (y/N): " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
        print_info "Installation cancelled."
        exit 0
    fi
    
    # Check if we're in the repo directory
    if [ ! -f "go.mod" ] || [ ! -f "main.go" ]; then
        print_error "This script must be run from the openkommander repository root directory."
        print_info "Please clone the repository first:"
        print_info "  git clone $REPO_URL"
        print_info "  cd openkommander"
        print_info "  ./scripts/setup.sh"
        exit 1
    fi
    
    # Build the application
    print_info "Building application..."
    if ! make build; then
        print_error "Build failed!"
        exit 1
    fi
    
    # Build frontend
    print_info "Building frontend..."
    if ! make frontend-build; then
        print_error "Frontend build failed!"
        exit 1
    fi
    
    # Install binary
    print_info "Installing binary..."
    echo "This step requires sudo access to install the binary to /usr/local/bin/"
    
    if ! make install-sudo; then
        print_error "Installation failed!"
        exit 1
    fi
    
    print_success "OpenKommander installed successfully!"
    print_info "You can now run: $BINARY_NAME --help"
}

# Function to check container runtime
check_container_runtime() {
    local runtime=$1
    
    if ! command_exists "$runtime"; then
        print_error "$runtime is not installed or not in PATH"
        return 1
    fi
    
    # Check if Docker/Podman daemon is running
    if ! $runtime ps >/dev/null 2>&1; then
        print_error "$runtime daemon is not running or you don't have permission to access it"
        return 1
    fi
    
    # Check compose availability
    local compose_cmd=""
    if $runtime compose version >/dev/null 2>&1; then
        compose_cmd="$runtime compose"
    elif command_exists "${runtime}-compose"; then
        compose_cmd="${runtime}-compose"
    else
        print_error "No compose tool found for $runtime"
        return 1
    fi
    
    print_success "$runtime is available with compose support" >&2
    echo "$compose_cmd"
    return 0
}

# Function to install with container
install_with_container() {
    local runtime=$1
    print_info "Starting $runtime installation..."
    
    # Check if container runtime is available
    local compose_cmd
    if ! compose_cmd=$(check_container_runtime "$runtime"); then
        print_error "$runtime is not properly set up"
        exit 1
    fi
    
    # Check if we're in the repo directory
    if [ ! -f "docker-compose.dev.yml" ] || [ ! -f "docker-compose.kafka.yml" ]; then
        print_error "Required compose files not found. This script must be run from the repository root."
        print_info "Expected files: docker-compose.dev.yml, docker-compose.kafka.yml"
        exit 1
    fi
    
    print_info "Using compose command: $compose_cmd"
    
    # Start Kafka clusters first
    print_info "Starting Kafka clusters (3 clusters)..."
    if ! $compose_cmd -f docker-compose.kafka.yml up -d; then
        print_error "Failed to start Kafka clusters"
        print_info "Try checking if ports 9092, 9095, 9098 are already in use"
        exit 1
    fi
    
    print_success "Kafka clusters started successfully!"
    print_info "Kafka clusters available at:"
    print_info "  - Cluster 1: localhost:9092"
    print_info "  - Cluster 2: localhost:9095"
    print_info "  - Cluster 3: localhost:9098"
    
    # Wait for Kafka clusters to be ready
    print_info "Waiting for Kafka clusters to be ready..."
    local max_wait=120
    local wait_time=0
    
    while [ $wait_time -lt $max_wait ]; do
        if $compose_cmd -f docker-compose.kafka.yml exec -T kafka-cluster1 kafka-topics.sh --bootstrap-server localhost:9092 --list >/dev/null 2>&1; then
            print_success "Kafka cluster 1 is ready!"
            break
        fi
        sleep 5
        wait_time=$((wait_time + 5))
        print_info "Waiting for Kafka... ($wait_time/${max_wait}s)"
    done
    
    if [ $wait_time -ge $max_wait ]; then
        print_warning "Kafka took longer than expected to be ready, continuing anyway..."
    fi
    
    # Start application containers
    print_info "Starting application..."
    if ! $compose_cmd -f docker-compose.dev.yml up --build -d; then
        print_error "Failed to start application containers"
        print_info "Kafka clusters are still running. Check application logs for details."
        exit 1
    fi
    
    print_success "Containers started successfully!"
    
    # Wait for containers to be ready
    print_info "Waiting for containers to be ready..."
    sleep 10
    
    # Execute into container and install
    print_info "Installing application in container..."
    
    # Use the same approach as the Makefile: exec into the 'app' service
    if ! $compose_cmd -f docker-compose.dev.yml exec app make dev-run; then
        print_warning "Could not run dev-run in container. Trying build instead..."
        if ! $compose_cmd -f docker-compose.dev.yml exec app make build; then
            print_error "Failed to build application in container"
            exit 1
        fi
    fi
    
    print_success "Application installed in $runtime container!"
    print_success "Multi-cluster Kafka environment ready!"
    print_info ""
    print_info "🎯 Available Kafka Clusters:"
    print_info "  • Cluster 1: localhost:9092"
    print_info "  • Cluster 2: localhost:9095" 
    print_info "  • Cluster 3: localhost:9098"
    print_info ""
    print_info "📋 Management Commands:"
    print_info "  - View app logs: make container-logs"
    print_info "  - View Kafka logs: make container-kafka-logs"
    print_info "  - Execute into container: make container-exec"
    print_info "  - Stop all services: make container-stop"
    print_info "  - Stop only Kafka: make container-kafka-stop"
    print_info ""
    print_info "🧪 Test Multi-Cluster Setup:"
    print_info "  ok cluster add cluster1 --bootstrap-server localhost:9092"
    print_info "  ok cluster add cluster2 --bootstrap-server localhost:9095"
    print_info "  ok cluster add cluster3 --bootstrap-server localhost:9098"
    
    # # Ask if user wants to exec into container
    # echo
    # read -p "Would you like to exec into the container now? (y/N): " exec_choice
    # if [[ "$exec_choice" =~ ^[Yy]$ ]]; then
    #     print_info "Executing into container..."
    #     print_info "You are now inside the container. Type 'exit' to return to your host system."
    #     $compose_cmd -f docker-compose.dev.yml exec app bash
    # fi
}

# Function to display welcome message
show_welcome() {
    echo
    echo "========================================"
    echo "  OpenKommander Setup Script"
    echo "========================================"
    echo
    print_info "This script will help you install OpenKommander on your system."
    print_info "You can choose between local installation or containerized deployment."
    echo
}

# Function to display completion message
show_completion() {
    local install_type=$1
    
    echo
    echo "========================================"
    echo "  Installation Complete!"
    echo "========================================"
    echo
    
    case $install_type in
        "local")
            print_success "OpenKommander has been installed locally."
            print_info "Run '$BINARY_NAME --help' to get started."
            ;;
        "docker"|"podman")
            print_success "OpenKommander and Kafka clusters are running in $install_type containers."
            print_info ""
            print_info "Kafka Clusters Available (External Access):"
            print_info "  • Cluster 1: localhost:9092 (kafka-cluster1)"
            print_info "  • Cluster 2: localhost:9095 (kafka-cluster2)" 
            print_info "  • Cluster 3: localhost:9098 (kafka-cluster3)"
            print_info ""
            print_info "Container-to-Container Access (for apps in containers):"
            print_info "  • Cluster 1: kafka-cluster1:9093"
            print_info "  • Cluster 2: kafka-cluster2:9093"
            print_info "  • Cluster 3: kafka-cluster3:9093"
            print_info ""
            print_info "Container Management Commands:"
            print_info "  make container-logs    # View application logs"
            print_info "  make container-exec    # Execute into container"
            print_info "  make container-stop    # Stop all containers"
            print_info "  make container-restart # Restart containers"
            print_info ""
            print_info "Kafka Management Commands:"
            print_info "  make kafka-up          # Start Kafka clusters"
            print_info "  make kafka-down        # Stop Kafka clusters"
            print_info "  make kafka-logs        # View Kafka logs"
            ;;
    esac
    
    echo
    print_info "For more information, visit: $REPO_URL"
    echo
}

# Main function
main() {
    show_welcome
    
    # Get installation type from user
    local install_type
    install_type=$(prompt_installation_type)
    
    case $install_type in
        "local")
            install_locally
            ;;
        "docker")
            install_with_container "docker"
            ;;
        "podman")
            install_with_container "podman"
            ;;
    esac
    
    show_completion "$install_type"
}

# Run main function
main "$@"
