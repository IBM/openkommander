# Multi-Cluster Kafka Setup

This setup provides 3 separate Kafka clusters for testing OpenKommander's multi-cluster management capabilities.

## Cluster Configuration

### Cluster 1 (Primary)
- **External Port**: `9092` (for external connections)
- **Internal Port**: `9093` (for inter-container communication)
- **Controller Port**: `9094`
- **Container**: `kafka-cluster1`
- **Bootstrap Server**: `localhost:9092`

### Cluster 2
- **External Port**: `9095` (for external connections)
- **Internal Port**: `9096` (for inter-container communication)
- **Controller Port**: `9097`
- **Container**: `kafka-cluster2`
- **Bootstrap Server**: `localhost:9095`

### Cluster 3
- **External Port**: `9098` (for external connections)
- **Internal Port**: `9099` (for inter-container communication)
- **Controller Port**: `9100`
- **Container**: `kafka-cluster3`
- **Bootstrap Server**: `localhost:9098`

## Usage

### Starting All Services
```bash
# Start both Kafka clusters and application
make container-start
```

### Starting Only Kafka Clusters
```bash
# Start only the 3 Kafka clusters
make container-kafka-start
```

### Viewing Logs
```bash
# View application logs
make container-logs

# View only Kafka cluster logs
make container-kafka-logs
```

### Stopping Services
```bash
# Stop everything
make container-stop

# Stop only Kafka clusters
make container-kafka-stop
```

### Manual Docker Compose Commands
```bash
# Start Kafka clusters
docker compose -f docker-compose.kafka.yml up -d

# Start application (requires Kafka to be running)
docker compose -f docker-compose.dev.yml up -d

# Stop everything
docker compose -f docker-compose.dev.yml down
docker compose -f docker-compose.kafka.yml down
```

## Application Configuration

The application is configured with:
- **KAFKA_BROKERS**: `localhost:9092,localhost:9095,localhost:9098` (all clusters)
- **KAFKA_BROKER**: `localhost:9092` (primary cluster for backward compatibility)

## Testing Multi-Cluster Functionality

You can test OpenKommander against all three clusters:

```bash
# Connect to cluster 1
ok cluster add cluster1 --bootstrap-server localhost:9092

# Connect to cluster 2  
ok cluster add cluster2 --bootstrap-server localhost:9095

# Connect to cluster 3
ok cluster add cluster3 --bootstrap-server localhost:9098

# List all clusters
ok cluster list
```

## Development Notes

- Each cluster runs in KRaft mode (no Zookeeper required)
- All clusters are configured for development with relaxed settings
- Auto-topic creation is enabled
- 512MB heap per cluster
- 24-hour log retention for faster testing

## Port Reference

| Cluster | External | Internal | Controller |
|---------|----------|----------|------------|
| 1       | 9092     | 9093     | 9094       |
| 2       | 9095     | 9096     | 9097       |
| 3       | 9098     | 9099     | 9100       |

## Troubleshooting

### Port Conflicts
If you get port conflicts, check what's using the ports:
```bash
lsof -i :9092
lsof -i :9095  
lsof -i :9098
```

### Health Checks
Check if clusters are healthy:
```bash
docker compose -f docker-compose.kafka.yml ps
```

### Individual Cluster Management
```bash
# Start specific cluster
docker compose -f docker-compose.kafka.yml up kafka-cluster1 -d

# View specific cluster logs
docker compose -f docker-compose.kafka.yml logs kafka-cluster2 -f

# Stop specific cluster
docker compose -f docker-compose.kafka.yml stop kafka-cluster3
```