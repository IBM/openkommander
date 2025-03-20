#!/bin/bash

echo "Stopping any existing containers..."
docker-compose down -v

echo "Removing any persisted data..."
docker volume prune -f

echo "Starting Kafka infrastructure..."
docker-compose up -d zookeeper kafka kafka-ui

echo "Waiting for Kafka to be ready..."
RETRIES=30
while [ $RETRIES -gt 0 ]
do
    if docker exec kafka /usr/bin/kafka-topics --bootstrap-server localhost:9092 --list 2>/dev/null; then
        break
    fi
    sleep 2
    RETRIES=$((RETRIES-1))
    echo "Waiting for Kafka... ($RETRIES attempts left)"
done

if [ $RETRIES -eq 0 ]; then
    echo "Failed to connect to Kafka. Exiting."
    docker-compose down -v
    exit 1
fi

echo "Creating Kafka topics..."
for TOPIC in order-created order-validated payment-processed inventory-updated shipping-prepared notification-sent order-completed order-failed
do
    docker exec kafka /usr/bin/kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic $TOPIC --partitions 3 --replication-factor 1
done

export KAFKAJS_NO_PARTITIONER_WARNING=1

echo "Starting all services..."
npm run start:analytics &
sleep 1
npm run start:inventory &
sleep 1
npm run start:payment &
sleep 1
npm run start:shipping &
sleep 1
npm run start:notification &
sleep 1
npm run start:order &

echo "All services started! Press Ctrl+C to stop"
wait