## E-commerce Microservice-based demo application powered by Apache Kafka

Built using:

- **Node.js**: Runtime for all microservices
- **KafkaJS**: Client library for Apache Kafka
- **Docker & Docker Compose**: Containerization and orchestration
- **Winston**: For logging

## Running the Demo

The demo can be started using the provided script, which:

1. Starts Zookeeper and Kafka containers
2. Creates the required Kafka topics
3. Starts all microservices in the correct order
4. Provides a Kafka UI for monitoring topics and messages at http://localhost:8090

The script ensures no persistence between runs, giving a clean state for each demonstration.