# OpenKommander

OpenKommander is a command-line tool and API server for managing Apache Kafka clusters. It provides a unified interface for common Kafka operations through both CLI and REST API, powered by [IBM Sarama](https://github.com/IBM/sarama) client library.

## ⚠️ DISCLAIMER ⚠️

**⚠️ THIS SOFTWARE IS NOT INTENDED FOR PRODUCTION USE  ⚠️**

OpenKommander is currently in early development and is provided AS IS without warranty of any kind. Using this software in production environments may lead to:

- Data loss or corruption
- Security vulnerabilities
- Performance degradation
- Operational instability

This tool is designed for development, testing, and educational purposes only. For production Kafka management, please use officially supported tools from the Apache Kafka ecosystem or commercial alternatives with proper support and security assurances.

**⚠️ By using this software, you acknowledge the risks and agree that the authors and contributors cannot be held liable for any damages resulting from its use. ⚠️**

## Architecture

OpenKommander follows a modular architecture with these key components:

- **CLI Commands**: Organized by functionality (topics, brokers, consumers, messages)
- **API Server**: RESTful interface with corresponding endpoints to CLI commands
- **[IBM Sarama](https://github.com/IBM/sarama) client library**: Core abstraction for interacting with Kafka clusters
- **Configuration Management**: Flexible configuration for multiple clusters

## Key Features

- Topic management (create, list, delete, describe)
- Message producing and consuming
- Consumer group monitoring
- Broker information
- Multi-cluster support
- Flexible authentication (SASL, TLS)
- Interactive ReactJS-based web dashboard

## Components

- **Models**: Data structures for configuration and Kafka entities
- **Lib**: Core utilities and client implementation
- **Components**: Feature-specific implementations (CLI and API)
- **Server**: API server implementation with route registration

## CLI Commands

```bash
# Configuration and startup
ok connect                                # Initialize default configuration
ok server                                 # Start the API server

# Topic Management
ok topics list                            # List all topics
ok topics create [topic-name] -p [partitions] -r [replication-factor]
ok topics delete [topic-name]             # Delete a topic
ok topics describe [topic-name]           # Show topic details
ok topics consume [topic] [flags]         # Consume messages from a topic
  --group string                          # Consumer group ID (optional)
  --from-beginning                        # Consume messages from beginning of the topic

# Consumer Group Management
ok consumers list                         # List all consumer groups
ok consumers describe [group-id]          # Describe a consumer group

# Broker Management
ok brokers                                # List information about Kafka brokers

# Message Production
ok produce [topic]                        # Produce a message to a Kafka topic
  --key/-k string                         # Message key
  --value/-v string                       # Message value
  --file/-f string                        # Read message from file
  --json/-j                               # Treat input as JSON

# Cluster Management
ok clusters list                          # List configured clusters

# Global flags
  --config string                         # Path to config file (defaults to $HOME/.config/openkommander.json)
  --cluster string                        # Use named cluster from config
```

## API Endpoints

```
# Health Check
GET /api/v1/health                        # Server health check

# Topic Management
GET /api/v1/topics                        # List all topics
POST /api/v1/topics                       # Create a topic
GET /api/v1/topics/:name                  # Get topic details
DELETE /api/v1/topics/:name               # Delete a topic

# Message Operations
POST /api/v1/messages/:topic              # Produce a message to a topic

# Broker Information
GET /api/v1/brokers                       # List broker information

# Consumer Group Management
GET /api/v1/consumers                     # List all consumer groups
GET /api/v1/consumers/:group              # Get consumer group details
DELETE /api/v1/consumers/:id              # Stop a consumer

# Cluster Management (multi-cluster mode)
GET /api/v1/clusters                      # List all configured clusters
GET /api/v1/clusters/:name                # Get cluster details

# Cluster-specific endpoints
GET /api/v1/clusters/:name/topics         # List topics in specific cluster
POST /api/v1/clusters/:name/topics        # Create topic in specific cluster
GET /api/v1/clusters/:name/brokers        # List brokers in specific cluster
GET /api/v1/clusters/:name/consumers      # List consumer groups in specific cluster
```

## Configuration

Configuration is stored in `$HOME/.config/openkommander.json` with support for:
- Multiple Kafka clusters
- SASL authentication
- TLS encryption

## Limitations

- No schema registry support
- No ACL management support
- Limited message transformation options
- etc.

## TODO

- OpenAPI contract between API and frontend SPA
- Implement schema registry support
- Implement ACL support
- Ensure 1:1 correspondence in functionality between CLI and API
- Kafka Connect 
- mBeans metrics collection and visualization
- Complete dockerization