# Scripts Directory

This directory contains utility scripts for OpenKommander development and deployment.

## Available Scripts

### `setup.sh` - Interactive Installation Script

An interactive script that guides users through different installation options for OpenKommander.

#### Features:

- **Interactive Prompts**: Guides users through installation choices
- **Multiple Installation Methods**: Local, Docker, or Podman
- **Prerequisite Checking**: Validates system requirements before installation
- **Error Handling**: Provides clear error messages and recovery instructions
- **Colored Output**: Easy-to-read colored terminal output

#### Usage:

```bash
# Make sure you're in the repository root
cd openkommander

# Run the setup script
./scripts/setup.sh
```

#### Installation Options:

1. **Local Installation**:

   - Checks for Go 1.24+, Node.js 18+, npm, make, and git
   - Builds application and frontend
   - Installs binary to `/usr/local/bin/` (requires sudo)
   - Best for development and direct system installation
2. **Docker Installation**:

   - Requires Docker with compose support
   - Starts application in Docker containers
   - Builds and installs within container environment
   - Best for containerized deployments
3. **Podman Installation**:

   - Requires Podman with compose support
   - Similar to Docker but uses Podman runtime
   - Alternative container solution for Docker

#### Prerequisites by Installation Type:

**Local Installation:**

- Go 1.24 or higher
- Node.js 18 or higher
- npm (usually comes with Node.js)
- make (build tools)
- git
- sudo access for final installation

**Docker Installation:**

- Docker installed and running
- Docker Compose (or `docker compose` command)
- Access to Docker daemon

**Podman Installation:**

- Podman installed and running
- Podman Compose (or `podman compose` command)
- Access to Podman daemon

#### Example Output:

```bash
$ ./scripts/setup.sh

========================================
  OpenKommander Setup Script
========================================

[INFO] This script will help you install OpenKommander on your system.
[INFO] You can choose between local installation or containerized deployment.

[INFO] Choose installation method:
1) Local installation (installs directly on your system)
2) Docker (runs in Docker container)
3) Podman (runs in Podman container)

Enter your choice (1-3): 1

[INFO] Checking local prerequisites...
[SUCCESS] Go go1.23.1 found
[SUCCESS] Node.js v18.17.0 found
[SUCCESS] npm 9.6.7 found
[SUCCESS] make found
[SUCCESS] git found
[SUCCESS] All prerequisites satisfied!

All prerequisites are satisfied. Continue with local installation? (y/N): y
...
```

#### Troubleshooting:

**Missing Prerequisites:**

- Follow the URLs provided in error messages to install missing dependencies
- On macOS: Install Xcode command line tools for make
- On Ubuntu/Debian: Install build-essential for make

**Container Issues:**

- Ensure Docker/Podman daemon is running
- Check user permissions for Docker/Podman access
- Verify compose tool availability

**Permission Issues:**

- Local installation requires sudo access for the final binary installation step
- Make sure you have administrative privileges on your system
