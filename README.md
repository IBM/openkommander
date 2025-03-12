# OpenKommander

OpenKommander is a command line utility and admin UI for Apache Kafka compatible brokers.

## Prerequisites

- [Podman](https://podman.io/getting-started/installation) (required for running the development environment)
- [Make](https://www.gnu.org/software/make/) (required for running development commands)

## Development Environment Setup

1. **Install Podman**
   Follow the installation instructions for your operating system on the [Podman website](https://podman.io/getting-started/installation).

2. **Clone the repository**
   ```bash
   git clone https://github.com/IBM/openkommander.git
   cd openkommander
   ```

3. **Start the development environment**
   ```bash
   make setup
   make dev
   ```

4. **Execute into the container**
    ```bash
    podman exec -it openkommander-app-1 bash
    ```

    Note: Replace `openkommander-app-1` in case your is different

5. **Build and install cli**
    ```bash
    make dev-run
    ```

## Configuration

The application uses a configuration file located at `config/config.yaml`. By default, it is configured for the development environment:

```yaml
kafka:
  broker: kafka:9093
```

### Custom Configuration

You can modify `config/config.yaml` to connect to different Kafka clusters:

```yaml
# Development environment (default)
kafka:
  broker: kafka:9093

# Custom environment example
kafka:
  broker: localhost:9092  # For local Kafka installation
  # broker: kafka-cluster.example.com:9093  # For remote cluster
```

The configuration file is loaded when the application starts. If you need to connect to a different broker after starting the application, you can use the `ok login` command with a custom broker address:

```bash
$ ok login
Enter broker address [kafka:9093]: localhost:9092
```

## CLI Usage

After running the application, you can use the following commands:


### Commands

All commands start with prefix `ok`

| Command | Description | Arguments |
|---------|-------------|-----------|
| `login` | Connect to a Kafka cluster | None |
| `logout` | End the current session | None |
| `session` | Display current session information | None |
| `metadata` | Display cluster information | None |
| `topics` | Topic management commands | Subcommands: `create`, `list`, `delete` |
| `help` | Display available commands | None |

### Topics Management

OpenKommander provides a set of commands to manage Kafka topics:

| Command | Description | Interactive Prompts |
|---------|-------------|-------------------|
| `topics create` | Create a new Kafka topic | Topic name, partitions, replication factor |
| `topics list` | List all available topics | None |
| `topics delete` | Delete an existing topic | Topic name |

| Endpoint | Method | Description | Request Body | Response |
|----------|--------|-------------|-------------|----------|
| `/topics` | GET | List all topics | None | JSON object with topic details |
| `/topics` | POST | Create a new topic | JSON with name, partitions, and replication_factor | Success message |
| `/topics/{topicName}` | DELETE | Delete a topic | None | Success message |

### Example Workflow

1. Build the CLI:
   ```bash
   make build
   ```

2. Install the CLI:
   ```bash
   make install
   ```

   Note: It may fail due to permission if needed add `sudo` for example
   ```bash
   sudo make install
   ```

3. Connect to the cluster:
   ```bash
   $ ok login
   Connected to Kafka cluster at kafka:9093
   ```

4. View session information:
   ```bash
   $ ok session
   Current session: Brokers: [kafka:9093], Authenticated: true
   ```

5. View cluster information:
   ```bash
   $ ok metadata
   Cluster Brokers:
    - kafka:9093 (ID: 1)
   ```

6. Create a new topic:
   ```bash
   $ ok topics create
   Enter topic name: my-new-topic
   Enter number of partitions (default 1): 3
   Enter replication factor (default 1): 2
   Successfully created topic 'my-new-topic' with 3 partitions and replication factor 2
   ```

7. List all topics:
   ```bash
   $ ok topics list
   Topics:
   --------
   Name: my-new-topic
   Partitions: 3
   Replication Factor: 2
   ```

8. Describe a topic:
   ```bash
   $ ok topics describe my-new-topic
   Topic Metadata:
      Topic Name: my-new-topic
      Replication Factor: 1
      Version: 10
      UUID: HHdnzvFrRpy1qIuLZuNO-w
      Is Internal: false
      Authorized Operations: -2147483648

   Topic Partitions:

   | PARTITION ID | LEADER | REPLICAS | IN-SYNC REPLICAS (ISR) |
   |--------------|--------|----------|------------------------|
   | 0            | 1      | [1 2]    | [1 2]                  |
   | 1            | 2      | [2 3]    | [2]                    |
   | 2            | 3      | [3 1]    | [3 1]                  |

   Topic Configurations:

   | CONFIG NAME                             | VALUE               |
   |-----------------------------------------|---------------------|
   | compression.type                        | producer            |
   | remote.log.delete.on.disable            | false               |
   | leader.replication.throttled.replicas   |                     |
   | remote.storage.enable                   | false               |
   | message.downconversion.enable           | true                |
   | min.insync.replicas                     | 1                   |
   | segment.jitter.ms                       | 0                   |
   | remote.log.copy.disable                 | false               |
   | local.retention.ms                      | -2                  |
   | cleanup.policy                          | delete              |
   | flush.ms                                | 9223372036854775807 |
   ......

   ```
9. Delete a topic:
   ```bash
   $ ok topics delete
   Enter topic name to delete: my-new-topic
   Successfully deleted topic 'my-new-topic'
   ```

9. End session and exit:
   ```bash
   $ ok logout
   Logged out successfully!
   ```