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
   make container-start
   ```
4. **Execute into the container**

   ```bash
   make container-exec
   ```
5. **Build and install cli**

   ```bash
   make dev-run
   ```

## Makefile Commands

The project includes several Make commands to help with development and deployment:

### Setup and Dependencies

| Command        | Description                                          |
| -------------- | ---------------------------------------------------- |
| `make setup` | Download Go dependencies and create config directory |

### Build Commands

| Command                 | Description                                         |
| ----------------------- | --------------------------------------------------- |
| `make build`          | Build the CLI binary for the host platform          |
| `make install`        | Build and install the CLI to `/usr/local/bin/ok`  |
| `make dev-run`        | Complete development setup: setup + build + install |
| `make frontend-build` | Build the frontend and copy to `~/.ok/frontend`   |

### Container Management

| Command                    | Description                                  |
| -------------------------- | -------------------------------------------- |
| `make container-start`   | Start the development environment containers |
| `make container-stop`    | Stop containers and clean up                 |
| `make container-restart` | Restart the development environment          |
| `make container-logs`    | View container logs                          |
| `make container-exec`    | Execute into the app container               |

### Cleanup

| Command        | Description                                            |
| -------------- | ------------------------------------------------------ |
| `make clean` | Remove build artifacts, dependencies, and config files |

## CLI Usage

After running the application, you can use the following commands:

### Commands

All commands start with prefix `ok`

| Command      | Description                         | Usage                                                                |
| ------------ | ----------------------------------- | -------------------------------------------------------------------- |
| `login`      | Connect to a Kafka cluster         | `ok login`   |
| `logout`     | End the current session            | `ok logout`                                                          |
| `session`    | Display current session information| `ok session`                                                         |
| `metadata`   | Display cluster information         | `ok metadata`                                                        |
| `produce`    | Produce messages to a topic         | `ok produce [TOPIC NAME] --msg/-m <message> [flags]`               |
| `server`     | REST server commands                | `ok server <subcommand>`                                            |
| `topic`      | Topic management commands           | `ok topic <subcommand>`                                             |
| `broker`     | Broker management commands          | `ok broker <subcommand>`                                            |
| `help`       | Display available commands          | `ok help`                                                           |

### Topic Management

OpenKommander provides comprehensive topic management commands:

| Command                              | Description                         | Usage                                               |
| ------------------------------------ | ----------------------------------- | --------------------------------------------------- |
| `ok topic create [TOPIC NAME]`      | Create a new Kafka topic           | `ok topic create my-topic -p 3 -r 2`              |
| `ok topic list`                     | List all available topics          | `ok topic list`                                     |
| `ok topic delete [TOPIC NAME]`      | Delete an existing topic           | `ok topic delete my-topic`                         |
| `ok topic describe [TOPIC NAME]`    | Describe an existing topic         | `ok topic describe my-topic`                       |
| `ok topic update [TOPIC NAME]`      | Update topic partition count       | `ok topic update my-topic -p 5`                    |

**Topic Create Flags:**
- `-p, --partitions`: Number of partitions (interactive prompt if not provided)
- `-r, --replication-factor`: Replication factor (interactive prompt if not provided)

**Topic Update Flags:**
- `-p, --new-partitions`: New partition count (required)

### Message Production

The `produce` command allows you to send messages to Kafka topics:

**Usage:**
```bash
ok produce [TOPIC NAME] --msg/-m <message> [flags]
```

**Flags:**
- `-m, --msg string`: Message payload (required)
- `-k, --key string`: [optional] Message key
- `-p, --partition int`: [optional] Partition to write message to (default -1)
- `-a, --acks int`: [optional] Acks flag, default -1 (full ISR)

**Examples:**
```bash
# Send a simple message to a topic
ok produce my-topic -m "Hello, Kafka!"

# Send a message with a key
ok produce my-topic -k "user123" -m "User login event"

# Send a message to a specific partition
ok produce my-topic -m "Hello, Kafka!" -p 0

# Send a message with custom acks setting
ok produce my-topic -m "Important message" -a 1
```

### Server Management

Start the REST API server:

| Command                     | Description            | Usage                                    |
| --------------------------- | ---------------------- | ---------------------------------------- |
| `ok server start`          | Start the REST server  | `ok server start -p 8081`              |

**Server Start Flags:**
- `-p, --port`: Port number for the REST server (required)

### Broker Management

View broker information:

| Command              | Description           | Usage              |
| -------------------- | --------------------- | ------------------ |
| `ok broker info`    | List all broker info  | `ok broker info`   |

### REST API Endpoints

The REST server provides HTTP endpoints for topic management:

| Endpoint                | Method | Description        | Request Body                                       | Response                       |
| ----------------------- | ------ | ------------------ | -------------------------------------------------- | ------------------------------ |
| `/topics`             | GET    | List all topics    | None                                               | JSON object with topic details |
| `/topics`             | POST   | Create a new topic | JSON with name, partitions, and replication_factor | Success message                |
| `/topics/{topicName}` | DELETE | Delete a topic     | None                                               | Success message                |

#### REST API Examples

**List topics:**
```bash
curl -X GET http://localhost:8081/api/v1/topics
```

**Create a topic:**
```bash
curl -X POST http://localhost:8081/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{"name":"my-topic","partitions":2,"replication_factor":1}'
```

**Delete a topic:**
```bash
curl -X DELETE http://localhost:8081/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{"name":"my-topic"}'
```

**Broker status:**
```bash
curl -X GET http://localhost:8081/api/v1/status
```

**Broker management:**
```bash
curl -X GET http://localhost:8081/api/v1/brokers
```

### Example Workflow

1. Setup and build the CLI:

   ```bash
   make setup
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
   $ ok topic create my-new-topic -p 3 -r 2
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
   $ ok topic describe my-new-topic
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
10. Start REST server:

    ```bash
    $ ok server start -p 8081 --brokers kafka:9093
    ```
11. Update a topic:

    ```bash
    $ ok topics update -n my-new-topic -p 4
    Successfully updated topic 'my-new-topic' to 4 partitions.
    ```
12. End session and exit:

    ```bash
    $ ok logout
    Logged out successfully!
    ```
