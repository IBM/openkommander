# OpenKommander Architecture

## System Architecture Overview

OpenKommander follows a three-tier architecture with a Go backend, REST API layer, and React frontend, all designed to interact with Apache Kafka clusters.

```
┌─────────────────────────────────────────────────────────────┐
│                         Users                                │
└────────────┬────────────────────────────┬───────────────────┘
             │                            │
             │ CLI                        │ Web Browser
             │                            │
┌────────────▼────────────┐  ┌───────────▼──────────────────┐
│   CLI Interface (ok)    │  │   React Frontend (Web UI)    │
│   - Cobra Commands      │  │   - Carbon Design System     │
│   - Interactive Prompts │  │   - Real-time Dashboard      │
│   - Formatted Output    │  │   - Resource Management      │
└────────────┬────────────┘  └───────────┬──────────────────┘
             │                            │
             │                            │ HTTP/REST
             │                            │
┌────────────▼────────────────────────────▼──────────────────┐
│              Go Backend Application                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           REST API Server (HTTP)                    │   │
│  │  - Health endpoints                                 │   │
│  │  - Topic management                                 │   │
│  │  - Broker information                               │   │
│  │  - Consumer group data                              │   │
│  │  - Message production                               │   │
│  └─────────────────┬───────────────────────────────────┘   │
│                    │                                        │
│  ┌─────────────────▼───────────────────────────────────┐   │
│  │         Core Business Logic Layer                   │   │
│  │  - commands/broker.go                               │   │
│  │  - commands/cluster.go                              │   │
│  │  - commands/topic.go                                │   │
│  │  - commands/produce.go                              │   │
│  │  - commands/status.go                               │   │
│  └─────────────────┬───────────────────────────────────┘   │
│                    │                                        │
│  ┌─────────────────▼───────────────────────────────────┐   │
│  │         Session Management Layer                    │   │
│  │  - Session persistence                              │   │
│  │  - Client connection pooling                        │   │
│  │  - Authentication state                             │   │
│  └─────────────────┬───────────────────────────────────┘   │
│                    │                                        │
│  ┌─────────────────▼───────────────────────────────────┐   │
│  │         Kafka Client Layer (Sarama)                 │   │
│  │  - ClusterAdmin interface                           │   │
│  │  - Client interface                                 │   │
│  │  - Producer interface                               │   │
│  └─────────────────┬───────────────────────────────────┘   │
└────────────────────┼────────────────────────────────────────┘
                     │
                     │ Kafka Protocol
                     │
┌────────────────────▼────────────────────────────────────────┐
│              Apache Kafka Cluster(s)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐                 │
│  │ Broker 1 │  │ Broker 2 │  │ Broker 3 │                 │
│  └──────────┘  └──────────┘  └──────────┘                 │
│                                                             │
│  Topics, Partitions, Consumer Groups, Messages              │
└─────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. Frontend Layer (React)

#### Technology Stack
- **React 18**: Modern React with hooks and concurrent features
- **Vite**: Fast build tool and dev server
- **Carbon Design System**: IBM's enterprise UI framework
- **React Router v6**: Client-side routing
- **Axios**: HTTP client for API calls

#### Key Components

##### Layout Components
- **MainLayout**: Main application shell with navigation
- **AppSideNav**: Sidebar navigation menu

##### Page Components
- **OverviewPage**: Dashboard with cluster metrics
- **TopicsPage**: Topic management interface
- **BrokersPage**: Broker information display
- **ConsumerGroupsPage**: Consumer group monitoring

##### Common Components
- **BaseTable**: Reusable data table component
- **BaseModal**: Modal dialog wrapper
- **ResourcePage**: Generic resource management page
- **StatusTag**: Status indicator component
- **ErrorNotification**: Error display component

##### Custom Hooks
- **useApi**: API call management with loading/error states
- **useResourceManager**: CRUD operations for resources

##### Services
- **api.js**: Centralized API client
- **localstorage.js**: Browser storage utilities

#### Data Flow
```
User Action → Component → Hook → API Service → Backend
                ↓
            State Update
                ↓
            Re-render
```

### 2. Backend Layer (Go)

#### Package Structure

##### Main Package (`main.go`)
- Application entry point
- Logger initialization
- CLI command execution

##### CLI Package (`pkg/cli/`)
- Cobra command definitions
- Command routing
- User input handling
- Output formatting

##### Core Commands (`internal/core/commands/`)

**broker.go**
- `GetBrokerClient()`: Retrieve Kafka client for broker operations

**cluster.go**
- `ListClusters()`: Get all cluster/broker information
- `GetClusterMetadata()`: Retrieve cluster metadata

**topic.go**
- `CreateTopic()`: Create new topics
- `DeleteTopic()`: Remove topics
- `ListTopics()`: Get all topics
- `DescribeTopic()`: Get topic details
- `DescribeTopicConfig()`: Get topic configuration
- `UpdateTopic()`: Modify topic partitions

**produce.go**
- `ProduceMessage()`: Send messages to topics

**status.go**
- `GetAdminClient()`: Get Sarama admin client
- `GetClient()`: Get Sarama client
- `NewFailure()`: Error handling utility

##### Session Package (`pkg/session/`)
- Session state management
- Broker connection persistence
- Client lifecycle management

##### Logger Package (`pkg/logger/`)
- Structured logging
- Configurable log levels
- Pretty formatting for development

#### Error Handling

Custom `Failure` type for consistent error responses:
```go
type Failure struct {
    Err      error
    HttpCode int
}
```

All command functions return `(result, *Failure)` for uniform error handling.

### 3. REST API Layer

#### API Design

**Base URL**: `/api/v1`

**Response Format**:
```json
{
  "status": "ok|error",
  "data": { ... },
  "message": "..."
}
```

#### Endpoint Categories

##### Health & Status
- `GET /health`: API health check
- `GET /status`: Broker status

##### Cluster Management
- `GET /clusters`: List all clusters
- `GET /clusters/{id}/metadata`: Get cluster metadata

##### Broker Management
- `GET /brokers`: List all brokers

##### Topic Management
- `GET /topics`: List all topics
- `POST /topics`: Create topic
- `GET /topics/{name}`: Get topic details
- `DELETE /topics/{name}`: Delete topic

##### Consumer Groups
- `GET /consumers`: List consumer groups
- `GET /consumers/{group}`: Get group details

##### Message Production
- `POST /messages/{topic}`: Produce message

### 4. Kafka Integration Layer

#### Sarama Client Usage

**ClusterAdmin Interface**:
- Topic CRUD operations
- Configuration management
- Partition management

**Client Interface**:
- Broker discovery
- Metadata retrieval
- Connection management

**Producer Interface**:
- Message production
- Partition selection
- Acknowledgment handling

#### Connection Management

```go
// Session maintains Kafka connections
type Session struct {
    brokers        []string
    client         sarama.Client
    adminClient    sarama.ClusterAdmin
    authenticated  bool
}
```

## Data Flow Patterns

### 1. CLI Command Flow

```
User Input → Cobra Command → Core Command Function → Sarama Client → Kafka
                                                                        ↓
User Output ← Formatted Response ← Result/Error ← Response ← Kafka
```

### 2. Web UI Flow

```
User Action → React Component → API Call → REST Handler → Core Command → Sarama → Kafka
                                                                                    ↓
UI Update ← State Update ← JSON Response ← HTTP Response ← Result ← Response ← Kafka
```

### 3. Real-time Metrics Flow

```
Timer → API Poll → Multiple Endpoints → Aggregate Data → Update State → Re-render
  ↓
3s interval
```

## Deployment Architecture

### Development Environment

```
┌─────────────────────────────────────────────────────────┐
│  Docker/Podman Compose Environment                      │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │ Kafka        │  │ Kafka        │  │ Kafka        │ │
│  │ Cluster 1    │  │ Cluster 2    │  │ Cluster 3    │ │
│  │ :9092        │  │ :9095        │  │ :9098        │ │
│  └──────────────┘  └──────────────┘  └──────────────┘ │
│                                                          │
│  ┌──────────────────────────────────────────────────┐  │
│  │ OpenKommander Container                          │  │
│  │ - Go Backend                                     │  │
│  │ - React Frontend (built)                         │  │
│  │ - Port 8081                                      │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Production Deployment Options

#### Option 1: Standalone Binary
```
OpenKommander CLI → Kafka Cluster(s)
```

#### Option 2: Containerized Service
```
Container (Go + Frontend) → Load Balancer → Kafka Cluster(s)
```

#### Option 3: Kubernetes Deployment
```
K8s Pod (OpenKommander) → K8s Service → Kafka StatefulSet
```

## Security Architecture

### Authentication Flow
```
User → Login Command → Session Creation → Broker Connection → Authenticated Session
```

### Session Security
- Session data stored in `~/.ok/`
- File permissions restrict access
- No passwords stored (relies on Kafka security)

### API Security
- CORS configuration for web access
- Session validation for all operations
- No built-in authentication (delegates to Kafka)

## Scalability Considerations

### Backend
- Stateless REST API (horizontal scaling possible)
- Connection pooling for Kafka clients
- Efficient session management

### Frontend
- Code splitting for faster loads
- Lazy loading of components
- Optimized bundle size

### Kafka Integration
- Multiple broker support
- Partition-aware operations
- Efficient metadata caching

## Performance Optimizations

### Backend
- Connection reuse via session management
- Minimal memory footprint
- Fast JSON serialization
- Efficient error handling

### Frontend
- Virtual scrolling for large tables
- Debounced API calls
- Optimistic UI updates
- Configurable polling intervals

### Kafka Operations
- Batch operations where possible
- Efficient metadata queries
- Smart partition selection

## Monitoring & Observability

### Logging
- Structured logging with levels (Debug, Info, Warn, Error)
- Configurable output format (pretty/JSON)
- Request/response logging in API

### Metrics (Available via UI)
- Broker count
- Topic count
- Messages per minute
- Consumer lag
- Cluster health status

### Error Tracking
- Consistent error format
- HTTP status codes
- User-friendly error messages
- Stack traces in debug mode

## Technology Decisions

### Why Go?
- Excellent Kafka client library (Sarama)
- Fast compilation and execution
- Strong concurrency support
- Easy deployment (single binary)

### Why React?
- Component reusability
- Large ecosystem
- Excellent developer experience
- Strong community support

### Why Carbon Design System?
- Enterprise-grade components
- Accessibility built-in
- Consistent design language
- IBM backing and support

### Why Sarama?
- Mature and stable
- Full Kafka protocol support
- Active maintenance
- Production-proven

## Future Architecture Enhancements

### Planned Improvements
1. **Microservices Split**: Separate API and worker services
2. **Message Queue**: Add async job processing
3. **Caching Layer**: Redis for metadata caching
4. **Database**: Persist metrics and history
5. **Authentication**: OAuth2/OIDC integration
6. **Multi-tenancy**: Support multiple users
7. **WebSocket**: Real-time updates without polling
8. **Schema Registry**: Avro/Protobuf support

### Scalability Roadmap
1. Horizontal scaling of API servers
2. Load balancing for high availability
3. Distributed session management
4. Metrics aggregation service
5. Historical data storage and analysis