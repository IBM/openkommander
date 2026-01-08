# OpenKommander Development Guide

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure](#project-structure)
3. [Building the Project](#building-the-project)
4. [Running Tests](#running-tests)
5. [Development Workflow](#development-workflow)
6. [Code Style Guidelines](#code-style-guidelines)
7. [Contributing](#contributing)
8. [Debugging](#debugging)
9. [Release Process](#release-process)

---

## Development Environment Setup

### Prerequisites

- **Go**: Version 1.24 or higher
- **Node.js**: Version 18 or higher
- **npm**: Version 9 or higher
- **Make**: GNU Make
- **Podman** or **Docker**: For containerized development
- **Git**: Version control

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/IBM/openkommander.git
   cd openkommander
   ```

2. **Install Go dependencies**:
   ```bash
   make setup
   ```

3. **Install frontend dependencies**:
   ```bash
   cd frontend
   npm install
   cd ..
   ```

4. **Start development environment**:
   ```bash
   make container-start
   ```

5. **Verify setup**:
   ```bash
   make container-logs
   ```

### IDE Setup

#### VS Code

Recommended extensions:
- Go (golang.go)
- ESLint (dbaeumer.vscode-eslint)
- Prettier (esbenp.prettier-vscode)
- Docker (ms-azuretools.vscode-docker)

**Settings** (`.vscode/settings.json`):
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "editor.formatOnSave": true,
  "go.formatTool": "goimports",
  "[javascript]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  },
  "[javascriptreact]": {
    "editor.defaultFormatter": "esbenp.prettier-vscode"
  }
}
```

#### GoLand / IntelliJ IDEA

1. Open project directory
2. Enable Go modules support
3. Configure Go SDK (1.24+)
4. Enable ESLint for frontend code

---

## Project Structure

```
openkommander/
├── internal/                     # Internal packages (not importable)
│   └── core/
│       └── commands/            # Core business logic
│           ├── broker.go        # Broker operations
│           ├── cluster.go       # Cluster management
│           ├── topic.go         # Topic operations
│           ├── produce.go       # Message production
│           └── status.go        # Status and error handling
├── pkg/                         # Public packages (importable)
│   ├── cli/                    # CLI implementation
│   ├── cluster/                # Cluster utilities
│   ├── constants/              # Application constants
│   ├── logger/                 # Logging utilities
│   ├── rest/                   # REST API server
│   └── session/                # Session management
├── frontend/                    # React web application
│   ├── src/
│   │   ├── components/         # React components
│   │   │   ├── common/        # Reusable components
│   │   │   ├── layout/        # Layout components
│   │   │   └── overview/      # Overview-specific components
│   │   ├── pages/             # Page components
│   │   ├── services/          # API services
│   │   ├── hooks/             # Custom React hooks
│   │   └── config/            # Configuration
│   ├── public/                # Static assets
│   └── package.json           # Frontend dependencies
├── docs/                       # Documentation
├── scripts/                    # Utility scripts
├── tests/                      # Test files
├── docker/                     # Docker configurations
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── go.sum                      # Go dependency checksums
├── Makefile                    # Build automation
└── README.md                   # Project README
```

### Key Directories

- **`internal/core/commands/`**: Core Kafka operations
- **`pkg/cli/`**: CLI command definitions using Cobra
- **`pkg/session/`**: Session state management
- **`frontend/src/components/`**: React UI components
- **`frontend/src/services/`**: API client code

---

## Building the Project

### Backend (Go)

**Build CLI binary**:
```bash
make build
```

Output: `./openkommander` (or `./openkommander.exe` on Windows)

**Install CLI globally**:
```bash
sudo make install
```

Installs to: `/usr/local/bin/ok`

**Development build and install**:
```bash
make dev-run
```

This runs: `setup` → `build` → `install`

### Frontend (React)

**Development server**:
```bash
cd frontend
npm run dev
```

Access at: `http://localhost:5173`

**Production build**:
```bash
cd frontend
npm run build
```

Output: `frontend/dist/`

**Build and copy to config directory**:
```bash
make frontend-build
```

Copies to: `~/.ok/frontend/`

### Full Build

Build both backend and frontend:
```bash
make build
make frontend-build
```

---

## Running Tests

### Backend Tests

**Run all tests**:
```bash
go test ./...
```

**Run tests with coverage**:
```bash
go test -cover ./...
```

**Run tests with verbose output**:
```bash
go test -v ./...
```

**Run specific package tests**:
```bash
go test ./internal/core/commands/
```

**Generate coverage report**:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests

**Run tests**:
```bash
cd frontend
npm test
```

**Run tests with coverage**:
```bash
cd frontend
npm run test:coverage
```

**Run tests in watch mode**:
```bash
cd frontend
npm run test:watch
```

### Integration Tests

See `tests/TESTING_GUIDELINES.md` for integration testing procedures.

---

## Development Workflow

### 1. Feature Development

**Create a feature branch**:
```bash
git checkout -b feature/my-new-feature
```

**Make changes and test**:
```bash
# Backend changes
make build
./openkommander [command]

# Frontend changes
cd frontend
npm run dev
```

**Commit changes**:
```bash
git add .
git commit -m "feat: add new feature"
```

### 2. Code Review Process

1. Push branch to GitHub
2. Create Pull Request
3. Address review comments
4. Ensure CI passes
5. Merge when approved

### 3. Local Testing

**Start development environment**:
```bash
make container-start
```

**Execute into container**:
```bash
make container-exec
```

**View logs**:
```bash
make container-logs
```

**Stop environment**:
```bash
make container-stop
```

### 4. Hot Reload Development

**Backend** (requires manual restart):
```bash
# Terminal 1: Watch for changes
make build && ./openkommander server start -p 8081

# Terminal 2: Make changes and rebuild
make build
```

**Frontend** (automatic hot reload):
```bash
cd frontend
npm run dev
# Changes auto-reload in browser
```

---

## Code Style Guidelines

### Go Code Style

**Follow standard Go conventions**:
- Use `gofmt` for formatting
- Use `golint` for linting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)

**Naming conventions**:
```go
// Good
func CreateTopic(name string) error { }
var topicName string

// Avoid
func create_topic(name string) error { }
var topic_name string
```

**Error handling**:
```go
// Always check errors
result, err := someFunction()
if err != nil {
    return nil, NewFailure(fmt.Sprintf("operation failed: %v", err), http.StatusInternalServerError)
}
```

**Comments**:
```go
// CreateTopic creates a new Kafka topic with the specified configuration.
// Returns success message and nil on success, or empty string and Failure on error.
func CreateTopic(topicName string, numPartitions, replicationFactor int) (string, *Failure) {
    // Implementation
}
```

**Package organization**:
- Keep packages focused and cohesive
- Avoid circular dependencies
- Use internal packages for non-exported code

### JavaScript/React Code Style

**Use Prettier for formatting**:
```bash
cd frontend
npm run format
```

**ESLint for linting**:
```bash
cd frontend
npm run lint
```

**Component structure**:
```jsx
// Good: Functional component with hooks
import React, { useState, useEffect } from 'react';

const MyComponent = ({ prop1, prop2 }) => {
  const [state, setState] = useState(null);
  
  useEffect(() => {
    // Effect logic
  }, []);
  
  return (
    <div>
      {/* JSX */}
    </div>
  );
};

export default MyComponent;
```

**Naming conventions**:
- Components: PascalCase (`MyComponent.jsx`)
- Hooks: camelCase with 'use' prefix (`useMyHook.jsx`)
- Utilities: camelCase (`apiClient.js`)
- Constants: UPPER_SNAKE_CASE (`API_BASE_URL`)

**File organization**:
```
components/
  MyComponent/
    MyComponent.jsx
    MyComponent.scss
    MyComponent.test.jsx
```

---

## Contributing

### Before Contributing

1. Read [CONTRIBUTING.md](../CONTRIBUTING.md)
2. Check [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md)
3. Review existing issues and PRs
4. Discuss major changes in an issue first

### Contribution Workflow

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/my-feature
   ```
3. **Make your changes**
4. **Add tests** for new functionality
5. **Run tests**:
   ```bash
   go test ./...
   cd frontend && npm test
   ```
6. **Commit with conventional commits**:
   ```bash
   git commit -m "feat: add new feature"
   ```
7. **Push to your fork**:
   ```bash
   git push origin feature/my-feature
   ```
8. **Create Pull Request**

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples**:
```bash
feat(topics): add topic update functionality
fix(api): correct broker status endpoint
docs: update installation instructions
refactor(cli): simplify command structure
```

---

## Debugging

### Backend Debugging

**Enable debug logging**:
```go
// In main.go
config := &logger.Config{
    Level:     logger.LevelDebug,
    Format:    "pretty",
    AddColors: true,
}
```

**Use Delve debugger**:
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug CLI
dlv debug . -- topic list

# Debug with breakpoints
dlv debug .
(dlv) break main.main
(dlv) continue
```

**VS Code debugging** (`.vscode/launch.json`):
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug CLI",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}",
      "args": ["topic", "list"]
    }
  ]
}
```

### Frontend Debugging

**Browser DevTools**:
- Open Chrome/Firefox DevTools (F12)
- Use React DevTools extension
- Check Console for errors
- Use Network tab for API calls

**VS Code debugging**:
```json
{
  "type": "chrome",
  "request": "launch",
  "name": "Debug Frontend",
  "url": "http://localhost:5173",
  "webRoot": "${workspaceFolder}/frontend/src"
}
```

**Debug API calls**:
```javascript
// In api.js, add logging
console.log('API Request:', method, endpoint, data);
const response = await axios[method](url, data);
console.log('API Response:', response.data);
```

### Container Debugging

**View container logs**:
```bash
make container-logs
```

**Execute into container**:
```bash
make container-exec
```

**Inspect Kafka logs**:
```bash
docker compose -f docker-compose.kafka.yml logs kafka-cluster1
```

**Check container status**:
```bash
docker compose -f docker-compose.dev.yml ps
```

---

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- MAJOR.MINOR.PATCH (e.g., 1.2.3)
- MAJOR: Breaking changes
- MINOR: New features (backward compatible)
- PATCH: Bug fixes

### Creating a Release

1. **Update version**:
   ```bash
   # Update version in relevant files
   # - main.go
   # - package.json
   # - README.md
   ```

2. **Update CHANGELOG**:
   ```markdown
   ## [1.2.0] - 2026-01-08
   ### Added
   - New feature X
   ### Fixed
   - Bug Y
   ### Changed
   - Improvement Z
   ```

3. **Create release branch**:
   ```bash
   git checkout -b release/v1.2.0
   ```

4. **Build and test**:
   ```bash
   make build
   make frontend-build
   go test ./...
   cd frontend && npm test
   ```

5. **Tag release**:
   ```bash
   git tag -a v1.2.0 -m "Release version 1.2.0"
   git push origin v1.2.0
   ```

6. **Create GitHub Release**:
   - Go to GitHub Releases
   - Create new release from tag
   - Add release notes
   - Attach binaries

### Build Artifacts

**Build for multiple platforms**:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o openkommander-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o openkommander-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o openkommander-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o openkommander-windows-amd64.exe
```

---

## Useful Commands

### Makefile Commands

```bash
# Setup
make setup                  # Download dependencies

# Build
make build                  # Build CLI
make install               # Install CLI globally
make dev-run               # Setup + build + install
make frontend-build        # Build frontend

# Container Management
make container-start       # Start containers
make container-stop        # Stop containers
make container-restart     # Restart containers
make container-logs        # View logs
make container-exec        # Execute into container

# Kafka Clusters
make container-kafka-start # Start Kafka clusters only
make container-kafka-stop  # Stop Kafka clusters
make container-kafka-logs  # View Kafka logs

# Cleanup
make clean                 # Remove build artifacts
```

### Go Commands

```bash
go mod tidy               # Clean up dependencies
go mod vendor             # Vendor dependencies
go fmt ./...              # Format code
go vet ./...              # Vet code
go test -v ./...          # Run tests verbosely
go build -race            # Build with race detector
```

### Frontend Commands

```bash
npm install               # Install dependencies
npm run dev              # Start dev server
npm run build            # Production build
npm run preview          # Preview production build
npm run lint             # Run ESLint
npm run format           # Format with Prettier
npm test                 # Run tests
```

---

## Troubleshooting Development Issues

### Go Module Issues

**Problem**: Module not found
```bash
go mod download
go mod tidy
```

### Frontend Build Issues

**Problem**: Dependencies out of sync
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Container Issues

**Problem**: Containers won't start
```bash
make container-stop
docker system prune -a
make container-start
```

### Port Conflicts

**Problem**: Port already in use
```bash
# Find process using port
lsof -i :8081

# Kill process
kill -9 <PID>
```

---

## Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://react.dev/)
- [Carbon Design System](https://carbondesignsystem.com/)
- [Sarama Documentation](https://pkg.go.dev/github.com/IBM/sarama)

### Tools
- [Delve Debugger](https://github.com/go-delve/delve)
- [golangci-lint](https://golangci-lint.run/)
- [Prettier](https://prettier.io/)
- [ESLint](https://eslint.org/)

### Community
- GitHub Issues
- GitHub Discussions
- Contributing Guidelines

---

## Next Steps

- Review [Architecture Documentation](ARCHITECTURE.md)
- Read [API Reference](API_REFERENCE.md)
- Check [User Guide](USER_GUIDE.md)
- Explore [Testing Guidelines](../tests/TESTING_GUIDELINES.md)