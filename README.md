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
   make dev
   ```

4. **Execute into the container**
    ```bash
    podman exec -it openkommander-app-1 bash
    ```

    Note: Replace `openkommander-app-1` in case your is different

5. **Run the application**
    ```bash
    go run main.go
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

The configuration file is loaded when the application starts. If you need to connect to a different broker after starting the application, you can use the `login` command with a custom broker address:

```bash
> login
Enter broker address [kafka:9093]: localhost:9092
```

## CLI Usage

After running the application, you can use the following commands:


### Commands

| Command | Description | Arguments |
|---------|-------------|-----------|
| `login` | Connect to a Kafka cluster | Optional broker address (default from config.yaml) |
| `logout` | End the current session | None |
| `session` | Display current session information | None |
| `metadata` | Display cluster information | None |
| `help` | Display available commands | None |
| `exit` | Exit the application | None |


### Example Workflow

1. Start the application:
   ```bash
   go run main.go
   ```

2. Connect to the cluster:
   ```bash
   > login
   Connected to Kafka cluster at kafka:9093
   ```

3. View session information:
   ```bash
   > session
   Current session: Brokers: [kafka:9093], Authenticated: true
   ```

4. View cluster information:
   ```bash
   > metadata
   Cluster Brokers:
    - kafka:9093 (ID: 1)
   ```

5. End session and exit:
   ```bash
   > logout
   Logged out successfully!
   > exit
   Exiting application.
   ```
