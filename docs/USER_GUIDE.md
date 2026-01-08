# OpenKommander User Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Getting Started](#getting-started)
4. [CLI Usage](#cli-usage)
5. [Web UI Usage](#web-ui-usage)
6. [Common Tasks](#common-tasks)
7. [Troubleshooting](#troubleshooting)
8. [Best Practices](#best-practices)

---

## Introduction

OpenKommander is a comprehensive tool for managing Apache Kafka clusters. It provides both a command-line interface (CLI) and a web-based user interface for performing common Kafka administration tasks.

### What You Can Do

- **Manage Topics**: Create, list, describe, update, and delete topics
- **Monitor Brokers**: View broker status and health
- **Produce Messages**: Send test messages to topics
- **Monitor Consumer Groups**: Track consumer lag and status
- **Manage Clusters**: Connect to and manage multiple Kafka clusters

---

## Installation

### Quick Installation (Recommended)

Use the interactive setup script:

```bash
git clone https://github.com/IBM/openkommander.git
cd openkommander
./scripts/setup.sh
```

The script will guide you through:
- Local installation (direct on your system)
- Docker installation
- Podman installation

### Manual Installation

#### Prerequisites
- Go 1.24 or higher
- Make
- Docker or Podman (for containerized setup)

#### Steps

1. **Clone the repository**:
   ```bash
   git clone https://github.com/IBM/openkommander.git
   cd openkommander
   ```

2. **Build the CLI**:
   ```bash
   make setup
   make build
   ```

3. **Install the CLI**:
   ```bash
   sudo make install
   ```

4. **Verify installation**:
   ```bash
   ok --version
   ```

---

## Getting Started

### First-Time Setup

1. **Start a Kafka cluster** (if you don't have one):
   ```bash
   make container-start
   ```

2. **Connect to your cluster**:
   ```bash
   ok login
   ```
   
   You'll be prompted for broker addresses. For the development environment:
   ```
   Enter broker addresses (comma-separated): localhost:9092
   ```

3. **Verify connection**:
   ```bash
   ok session
   ```

4. **View cluster information**:
   ```bash
   ok metadata
   ```

### Starting the Web UI

1. **Start the REST server**:
   ```bash
   ok server start -p 8081
   ```

2. **Open your browser**:
   Navigate to `http://localhost:8081`

---

## CLI Usage

### Basic Commands

#### Session Management

**Login to a cluster**:
```bash
ok login
```

**View current session**:
```bash
ok session
```

**Logout**:
```bash
ok logout
```

#### Cluster Information

**View cluster metadata**:
```bash
ok metadata
```

**List clusters/brokers**:
```bash
ok cluster list
```

**View broker information**:
```bash
ok broker info
```

### Topic Management

#### Create a Topic

**Interactive mode**:
```bash
ok topic create my-topic
```
You'll be prompted for partitions and replication factor.

**With flags**:
```bash
ok topic create my-topic -p 3 -r 1
```

**Flags**:
- `-p, --partitions`: Number of partitions
- `-r, --replication-factor`: Replication factor

#### List Topics

```bash
ok topic list
```

#### Describe a Topic

```bash
ok topic describe my-topic
```

This shows:
- Topic metadata (name, replication factor, UUID)
- Partition details (leader, replicas, ISR)
- Configuration settings

#### Update a Topic

Increase partition count:
```bash
ok topic update my-topic -p 5
```

**Note**: You can only increase partitions, not decrease them.

#### Delete a Topic

```bash
ok topic delete my-topic
```

You'll be prompted to confirm the deletion.

### Message Production

#### Basic Message

```bash
ok produce my-topic -m "Hello, Kafka!"
```

#### Message with Key

```bash
ok produce my-topic -k "user123" -m "User login event"
```

#### Specific Partition

```bash
ok produce my-topic -m "Hello" -p 0
```

#### Custom Acknowledgments

```bash
ok produce my-topic -m "Important message" -a 1
```

**Acknowledgment Options**:
- `-1`: Wait for all in-sync replicas (default, most reliable)
- `0`: No acknowledgment (fastest, least reliable)
- `1`: Wait for leader only (balanced)

### Server Management

**Start the REST server**:
```bash
ok server start -p 8081
```

**With specific brokers**:
```bash
ok server start -p 8081 --brokers localhost:9092,localhost:9095
```

---

## Web UI Usage

### Overview Dashboard

The Overview page displays real-time cluster metrics:

- **Broker Count**: Number of active brokers
- **Topic Count**: Total topics in the cluster
- **Health Status**: Overall cluster health
- **Messages Per Minute**: Message throughput

**Features**:
- Auto-refresh toggle (updates every 3 seconds)
- Last refresh timestamp
- Visual metric cards

### Topics Page

Manage all your Kafka topics from one interface.

**View Topics**:
- Table shows: name, partitions, replication factor, status
- Status indicators: Active (green), Warning (yellow), Inactive (red)

**Create Topic**:
1. Click "Add Topic" button
2. Fill in the form:
   - Topic Name
   - Number of Partitions
   - Replication Factor
3. Click "Create"

**Delete Topic**:
1. Click the delete icon next to the topic
2. Confirm deletion in the dialog

**View Details**:
- Click on a topic row to see detailed information
- View partition distribution
- Check replication status

### Brokers Page

Monitor broker health and status.

**Information Displayed**:
- Broker ID
- Host and Port
- Status (Selected/Inactive)
- Partition Count
- Leader Partitions
- In-Sync Partitions

**Features**:
- Real-time status updates
- Connection status indicators
- Partition distribution view

### Consumer Groups Page

Monitor consumer group activity and lag.

**Information Displayed**:
- Group ID
- Number of Members
- Subscribed Topics
- Total Lag
- Coordinator Broker
- State (Stable/Rebalancing/Dead)

**Features**:
- View topic-level lag
- Monitor consumer group health
- Track rebalancing events

**Status Indicators**:
- **Stable** (green): Group is healthy and consuming
- **Rebalancing** (yellow): Group is rebalancing partitions
- **Dead** (red): Group has no active members

---

## Common Tasks

### Task 1: Create and Test a New Topic

```bash
# 1. Create the topic
ok topic create test-topic -p 3 -r 1

# 2. Verify creation
ok topic list

# 3. View details
ok topic describe test-topic

# 4. Send a test message
ok produce test-topic -m "Test message"

# 5. Verify in web UI
# Open http://localhost:8081 and check Topics page
```

### Task 2: Monitor Consumer Lag

**Via CLI**:
```bash
# View all consumer groups
ok cluster list

# Check specific group (via API)
curl http://localhost:8081/api/v1/consumers/my-consumer-group
```

**Via Web UI**:
1. Navigate to Consumer Groups page
2. Look for groups with high lag values
3. Check if state is "Stable"
4. Review topic-level lag breakdown

### Task 3: Increase Topic Partitions

```bash
# 1. Check current partition count
ok topic describe my-topic

# 2. Increase partitions
ok topic update my-topic -p 6

# 3. Verify the change
ok topic describe my-topic
```

### Task 4: Connect to Multiple Clusters

```bash
# Connect to first cluster
ok login
# Enter: localhost:9092

# Perform operations...

# Connect to second cluster
ok logout
ok login
# Enter: localhost:9095

# Now working with second cluster
```

### Task 5: Bulk Message Production

Create a script to send multiple messages:

```bash
#!/bin/bash
for i in {1..100}; do
  ok produce my-topic -m "Message $i"
done
```

### Task 6: Export Topic List

```bash
# Using CLI
ok topic list > topics.txt

# Using API
curl http://localhost:8081/api/v1/topics > topics.json
```

---

## Troubleshooting

### Connection Issues

**Problem**: Cannot connect to Kafka cluster

**Solutions**:
1. Verify broker addresses:
   ```bash
   # Test connectivity
   telnet localhost 9092
   ```

2. Check if Kafka is running:
   ```bash
   make container-logs
   ```

3. Verify network access:
   ```bash
   # Check firewall rules
   # Ensure ports are not blocked
   ```

### Topic Creation Fails

**Problem**: "Replication factor cannot be greater than the number of brokers"

**Solution**: Reduce replication factor to match or be less than broker count:
```bash
# If you have 1 broker, use replication factor 1
ok topic create my-topic -p 3 -r 1
```

### Session Not Found

**Problem**: "Error: no session found"

**Solution**: Login to establish a session:
```bash
ok login
```

### Web UI Not Loading

**Problem**: Cannot access http://localhost:8081

**Solutions**:
1. Verify server is running:
   ```bash
   # Check if port is in use
   lsof -i :8081
   ```

2. Start the server:
   ```bash
   ok server start -p 8081
   ```

3. Try a different port:
   ```bash
   ok server start -p 8082
   ```

### Permission Denied

**Problem**: "Permission denied" when installing

**Solution**: Use sudo for installation:
```bash
sudo make install
```

### Port Already in Use

**Problem**: "Address already in use"

**Solutions**:
1. Find what's using the port:
   ```bash
   lsof -i :8081
   ```

2. Kill the process or use a different port:
   ```bash
   ok server start -p 8082
   ```

---

## Best Practices

### Topic Management

1. **Plan Partition Count**: Consider your throughput needs
   - More partitions = higher parallelism
   - But more partitions = more overhead
   - Start with 3-6 partitions for most use cases

2. **Set Appropriate Replication Factor**:
   - Production: Use replication factor of 3
   - Development: Replication factor of 1 is acceptable
   - Never exceed broker count

3. **Use Meaningful Names**:
   ```bash
   # Good
   ok topic create user-events -p 3 -r 1
   
   # Avoid
   ok topic create topic1 -p 3 -r 1
   ```

4. **Test Before Production**:
   - Create topics in dev environment first
   - Test with sample data
   - Verify consumer behavior

### Message Production

1. **Use Message Keys for Ordering**:
   ```bash
   # Messages with same key go to same partition
   ok produce orders -k "user123" -m "Order placed"
   ```

2. **Choose Appropriate Acknowledgments**:
   - Critical data: Use `-a -1` (all ISR)
   - Logs/metrics: Use `-a 1` (leader only)
   - High throughput: Use `-a 0` (no ack)

3. **Monitor Production Success**:
   - Check for errors in output
   - Verify messages in consumer groups

### Monitoring

1. **Regular Health Checks**:
   ```bash
   # Daily health check script
   ok metadata
   ok broker info
   ok cluster list
   ```

2. **Watch Consumer Lag**:
   - Check lag regularly via Web UI
   - Alert on high lag values
   - Investigate rebalancing events

3. **Enable Auto-Refresh in Web UI**:
   - Keep dashboard open for monitoring
   - Set appropriate refresh interval
   - Watch for status changes

### Security

1. **Protect Session Files**:
   ```bash
   # Session data stored in ~/.ok/
   chmod 700 ~/.ok
   ```

2. **Use Secure Connections**:
   - Configure SSL/TLS in Kafka
   - Use SASL authentication
   - Restrict network access

3. **Limit Access**:
   - Use Kafka ACLs
   - Restrict topic operations
   - Monitor audit logs

### Performance

1. **Batch Operations**:
   - Create multiple topics in scripts
   - Use API for bulk operations
   - Minimize individual CLI calls

2. **Optimize Partition Count**:
   - Balance between parallelism and overhead
   - Consider consumer count
   - Monitor broker load

3. **Monitor Resource Usage**:
   - Watch broker CPU/memory
   - Check disk usage
   - Monitor network throughput

---

## Advanced Usage

### Scripting with CLI

Create automation scripts:

```bash
#!/bin/bash
# setup-topics.sh

TOPICS=("orders" "payments" "notifications")

for topic in "${TOPICS[@]}"; do
  ok topic create "$topic" -p 3 -r 1
  echo "Created topic: $topic"
done
```

### Using the API Programmatically

```python
import requests

API_BASE = "http://localhost:8081/api/v1"

# Create multiple topics
topics = ["topic1", "topic2", "topic3"]
for topic in topics:
    response = requests.post(
        f"{API_BASE}/topics",
        json={
            "name": topic,
            "partitions": 3,
            "replication_factor": 1
        }
    )
    print(f"Created {topic}: {response.status_code}")
```

### Integration with CI/CD

```yaml
# .github/workflows/kafka-setup.yml
name: Setup Kafka Topics

on:
  push:
    branches: [main]

jobs:
  setup:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Install OpenKommander
        run: |
          git clone https://github.com/IBM/openkommander.git
          cd openkommander
          make build
          sudo make install
      - name: Create Topics
        run: |
          ok login
          ok topic create production-events -p 6 -r 3
```

---

## Getting Help

### Documentation
- Project Overview: `docs/PROJECT_OVERVIEW.md`
- Architecture: `docs/ARCHITECTURE.md`
- API Reference: `docs/API_REFERENCE.md`
- Development Guide: `docs/DEVELOPMENT_GUIDE.md`

### Command Help
```bash
# General help
ok help

# Command-specific help
ok topic --help
ok produce --help
```

### Community Support
- GitHub Issues: https://github.com/IBM/openkommander/issues
- Discussions: GitHub Discussions
- Contributing: See `CONTRIBUTING.md`

### Reporting Issues
When reporting issues, include:
1. OpenKommander version
2. Operating system
3. Kafka version
4. Steps to reproduce
5. Error messages
6. Relevant logs

---

## Next Steps

- Explore the [API Reference](API_REFERENCE.md) for programmatic access
- Read the [Development Guide](DEVELOPMENT_GUIDE.md) to contribute
- Check out [Architecture](ARCHITECTURE.md) to understand internals
- Review [Best Practices](#best-practices) for production use