# OpenKommander Project Overview

## Introduction

OpenKommander is a comprehensive command-line utility and web-based admin UI for managing Apache Kafka clusters. It provides both CLI and REST API interfaces for Kafka administration tasks, along with a modern React-based frontend for visual cluster management.

## Project Purpose

OpenKommander simplifies Kafka cluster management by providing:

- **Unified Interface**: Single tool for both command-line and web-based administration
- **Multi-Cluster Support**: Manage multiple Kafka clusters from one interface
- **Real-time Monitoring**: Live metrics and status updates for clusters, topics, and consumer groups
- **Developer-Friendly**: Easy setup with Docker/Podman support and comprehensive API
- **Production-Ready**: Built with IBM Sarama, a mature Kafka client library

## Key Features

### 1. Cluster Management
- Connect to multiple Kafka clusters
- View cluster metadata and broker information
- Monitor cluster health and status
- Support for KRaft mode (no Zookeeper required)

### 2. Topic Management
- Create, delete, and update topics
- List all topics with detailed information
- Describe topic configurations and partitions
- Update partition counts dynamically
- View replication status and ISR (In-Sync Replicas)

### 3. Message Production
- Produce messages to topics via CLI or API
- Support for message keys and partitioning
- Configurable acknowledgment settings
- Manual or automatic partition selection

### 4. Broker Management
- List all brokers in the cluster
- View broker status and connectivity
- Monitor partition distribution
- Check leader and replica assignments

### 5. Consumer Group Management
- List all consumer groups
- View consumer group status and lag
- Monitor topic assignments
- Track consumer group state (Stable, Rebalancing, Dead)

### 6. Web UI
- Modern React-based interface using Carbon Design System
- Real-time metrics dashboard
- Interactive tables for managing resources
- Auto-refresh capabilities with configurable polling
- Responsive design for desktop and mobile

### 7. REST API
- RESTful API for all operations
- OpenAPI 3.0 specification
- JSON request/response format
- Health check endpoints
- CORS support for web integration

## Architecture Components

### Backend (Go)
- **CLI Framework**: Cobra for command-line interface
- **Kafka Client**: IBM Sarama for Kafka operations
- **HTTP Server**: Built-in Go HTTP server for REST API
- **Session Management**: Persistent session handling
- **Logging**: Structured logging with configurable levels

### Frontend (React)
- **UI Framework**: React 18 with Vite
- **Design System**: IBM Carbon Design System
- **Routing**: React Router v6
- **State Management**: React Hooks
- **API Client**: Axios for HTTP requests
- **Styling**: SCSS with Carbon themes

### Infrastructure
- **Containerization**: Docker and Podman support
- **Development Environment**: Docker Compose for local development
- **Multi-Cluster Setup**: Support for 3+ Kafka clusters
- **Build System**: Make for automation

## Use Cases

### 1. Development Environment
- Quickly spin up local Kafka clusters for testing
- Experiment with topics and message production
- Test consumer applications with real-time monitoring

### 2. Operations & Monitoring
- Monitor cluster health and performance
- Track consumer lag and rebalancing
- Manage topic lifecycle (create, update, delete)
- View broker status and partition distribution

### 3. Troubleshooting
- Inspect topic configurations
- Check consumer group states
- Verify message production
- Analyze partition assignments

### 4. Automation & CI/CD
- Integrate with deployment pipelines
- Automate topic creation and configuration
- Script-based cluster management
- API-driven operations

### 5. Multi-Cluster Management
- Manage development, staging, and production clusters
- Compare configurations across environments
- Centralized monitoring dashboard

## Technology Stack

### Backend
- **Language**: Go 1.24+
- **Kafka Client**: IBM Sarama v1.46.3
- **CLI Framework**: Cobra v1.10.1
- **Table Rendering**: go-pretty v6.6.9

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **UI Library**: IBM Carbon Design System
- **HTTP Client**: Axios
- **Router**: React Router v6

### Development Tools
- **Container Runtime**: Docker/Podman
- **Build System**: Make
- **Version Control**: Git
- **Package Management**: Go modules, npm

## Project Structure

```
openkommander/
├── cmd/                    # CLI commands
├── internal/              # Internal packages
│   └── core/
│       └── commands/      # Core business logic
├── pkg/                   # Public packages
│   ├── cli/              # CLI implementation
│   ├── logger/           # Logging utilities
│   └── session/          # Session management
├── frontend/             # React web application
│   └── src/
│       ├── components/   # React components
│       ├── pages/        # Page components
│       ├── services/     # API services
│       └── hooks/        # Custom React hooks
├── docs/                 # Documentation
├── scripts/              # Setup and utility scripts
├── tests/                # Test files
└── docker/               # Docker configurations
```

## Getting Started

### Quick Start
```bash
# Clone the repository
git clone https://github.com/IBM/openkommander.git
cd openkommander

# Run interactive setup
./scripts/setup.sh
```

### Development Setup
```bash
# Start development environment
make setup
make container-start

# Build and install CLI
make dev-run

# Access the application
ok login
```

### Using the Web UI
```bash
# Start the REST server
ok server start -p 8081

# Access the UI at http://localhost:8081
```

## Configuration

### Environment Variables
- `KAFKA_BROKERS`: Comma-separated list of broker addresses
- `KAFKA_BROKER`: Primary broker address (backward compatibility)

### Session Storage
- Configuration stored in `~/.ok/`
- Session persistence across CLI invocations
- Frontend assets in `~/.ok/frontend/`

## Security Considerations

- No authentication required for local development
- Session-based access control
- CORS configuration for API access
- Secure broker connections supported

## Performance

- Efficient Kafka client connection pooling
- Minimal memory footprint
- Fast topic and broker operations
- Optimized frontend with code splitting
- Real-time updates with configurable polling intervals

## Limitations

- Currently focused on cluster administration (not data streaming)
- Consumer group management is read-only
- No built-in authentication/authorization (relies on Kafka security)
- Single-user session model

## Future Enhancements

- Schema registry integration
- Kafka Connect management
- Advanced consumer offset management
- Multi-user authentication
- Metrics persistence and historical data
- Alert and notification system
- Topic data browsing and search

## Contributing

OpenKommander is an open-source project. Contributions are welcome! See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0. See [LICENSE](../LICENSE) for details.

## Support

- **Issues**: GitHub Issues
- **Documentation**: `/docs` directory
- **Code of Conduct**: [CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md)
- **Security**: [SECURITY.md](../SECURITY.md)

## Acknowledgments

- Built with IBM Sarama Kafka client
- UI powered by IBM Carbon Design System
- Inspired by Kafka administration needs in enterprise environments