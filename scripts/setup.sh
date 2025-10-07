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
    if [ ! -f "docker-compose.dev.yml" ]; then
        print_error "docker-compose.dev.yml not found. This script must be run from the repository root."
        exit 1
    fi
    
    print_info "Using compose command: $compose_cmd"
    
    # Start containers
    print_info "Starting containers..."
    if ! $compose_cmd -f docker-compose.dev.yml up --build -d; then
        print_error "Failed to start containers"
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
    print_info "Container services are running. You can:"
    print_info "  - View logs: make container-logs"
    print_info "  - Execute into container: make container-exec"
    print_info "  - Stop containers: make container-stop"
    
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
            print_success "OpenKommander is running in $install_type containers."
            print_info "Available commands:"
            print_info "  make container-logs    # View application logs"
            print_info "  make container-exec    # Execute into container"
            print_info "  make container-stop    # Stop all containers"
            print_info "  make container-restart # Restart containers"
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
