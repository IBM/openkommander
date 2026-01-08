# OpenKommander API Reference

## Overview

The OpenKommander REST API provides programmatic access to Kafka cluster management operations. All endpoints return JSON responses and follow RESTful conventions.

**Base URL**: `http://localhost:8081/api/v1/{broker}`

**Note**: Most endpoints require a `{broker}` path parameter specifying the broker address (e.g., `localhost:9092`).

**Content-Type**: `application/json`

## Authentication

Currently, the API does not require authentication. Access control is managed through Kafka broker security settings.

## Response Format

### Success Response
```json
{
  "status": "ok",
  "data": { ... }
}
```

### Error Response
```json
{
  "status_code": 404,
  "message": "Resource not found",
  "error": "detailed error message"
}
```

## Common HTTP Status Codes

| Code | Meaning | Description |
|------|---------|-------------|
| 200 | OK | Request successful |
| 201 | Created | Resource created successfully |
| 202 | Accepted | Request accepted for processing |
| 400 | Bad Request | Invalid request parameters |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists |
| 500 | Internal Server Error | Server error occurred |

---

## Health & Status Endpoints

### Health Check

Check if the API is running and healthy.

**Endpoint**: `GET /{broker}/health`

**Path Parameters**:
- `broker` (string): Broker address (e.g., `localhost:9092`)

**Response**:
```json
{
  "status": "ok",
  "message": "Health check successful"
}
```

**Example**:
```bash
curl http://localhost:8081/api/v1/localhost:9092/health
```

---

## Broker Endpoints

### List Brokers

Get information about all brokers in the cluster.

**Endpoint**: `GET /{broker}/brokers`

**Path Parameters**:
- `broker` (string): Broker address (e.g., `localhost:9092`)

**Response**:
```json
{
  "status": "ok",
  "data": [
    {
      "id": 1,
      "addr": "localhost:9092",
      "connected": true,
      "rack": "",
      "state": null
    }
  ]
}
```

**Response Fields**:
- `id` (integer): Broker ID
- `addr` (string): Broker address
- `connected` (boolean): Connection status
- `rack` (string): Rack assignment
- `state` (object): TLS connection state (if applicable)

**Example**:
```bash
curl http://localhost:8081/api/v1/localhost:9092/brokers
```

---

## Cluster Endpoints

### List Clusters

Get information about all available clusters/brokers.

**Endpoint**: `GET /clusters`

**Response**:
```json
{
  "status": "ok",
  "data": [
    {
      "id": 1,
      "address": "localhost:9092",
      "status": "Connected",
      "rack": "rack-1",
      "connected": true
    }
  ]
}
```

**Response Fields**:
- `id` (integer): Cluster/broker ID
- `address` (string): Network address (host:port)
- `status` (string): Connection status ("Connected" or "Disconnected")
- `rack` (string): Rack assignment (or "N/A")
- `connected` (boolean): Connection state

**Example**:
```bash
curl http://localhost:8081/api/v1/clusters
```

### Get Cluster Metadata

Get detailed metadata about a specific cluster.

**Endpoint**: `GET /clusters/{clusterId}/metadata`

**Path Parameters**:
- `clusterId` (string): The cluster identifier

**Response**:
```json
{
  "status": "ok",
  "data": {
    "broker_count": 3,
    "cluster_id": "openkommander-client",
    "brokers": [
      {
        "id": 1,
        "address": "localhost:9092",
        "connected": true,
        "rack": "rack-1"
      }
    ]
  }
}
```

**Example**:
```bash
curl http://localhost:8081/api/v1/clusters/1/metadata
```

---

## Topic Endpoints

### List Topics

Get all topics in the cluster.

**Endpoint**: `GET /topics`

**Response**:
```json
[
  {
    "name": "test-topic",
    "partitions": 3,
    "replication_factor": 1,
    "internal": false,
    "replicas": 3,
    "in_sync_replicas": 3,
    "cleanup_policy": "delete"
  }
]
```

**Response Fields**:
- `name` (string): Topic name
- `partitions` (integer): Number of partitions
- `replication_factor` (integer): Replication factor
- `internal` (boolean): Whether topic is internal
- `replicas` (integer): Total replicas
- `in_sync_replicas` (integer): In-sync replicas
- `cleanup_policy` (string): Cleanup policy ("delete" or "compact")

**Example**:
```bash
curl http://localhost:8081/api/v1/topics
```

### Create Topic

Create a new Kafka topic.

**Endpoint**: `POST /topics`

**Request Body**:
```json
{
  "name": "new-topic",
  "partitions": 3,
  "replication_factor": 1
}
```

**Request Fields**:
- `name` (string, required): Topic name
- `partitions` (integer, required): Number of partitions (min: 1)
- `replication_factor` (integer, required): Replication factor (min: 1)

**Response**:
```json
{
  "status": "created"
}
```

**Example**:
```bash
curl -X POST http://localhost:8081/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-topic",
    "partitions": 3,
    "replication_factor": 1
  }'
```

**Error Responses**:
- `400`: Invalid parameters (partitions < 1, replication_factor > broker count)
- `500`: Topic already exists or creation failed

### Get Topic Details

Get detailed information about a specific topic.

**Endpoint**: `GET /topics/{name}`

**Path Parameters**:
- `name` (string): Topic name

**Response**:
```json
{
  "name": "test-topic",
  "partitions": 3,
  "replication_factor": 1,
  "partition_ids": [0, 1, 2]
}
```

**Example**:
```bash
curl http://localhost:8081/api/v1/topics/test-topic
```

### Delete Topic

Delete an existing topic.

**Endpoint**: `DELETE /topics/{name}`

**Path Parameters**:
- `name` (string): Topic name

**Response**:
```json
{
  "status": "deleted"
}
```

**Example**:
```bash
curl -X DELETE http://localhost:8081/api/v1/topics/test-topic
```

**Error Responses**:
- `500`: Deletion failed

---

## Error Handling

### Error Response Format

All errors follow this format:

```json
{
  "status_code": 404,
  "message": "User-friendly error message",
  "error": "Detailed technical error"
}
```

### Common Errors

#### 400 Bad Request
```json
{
  "status_code": 400,
  "message": "Invalid request parameters",
  "error": "Partitions must be greater than 0"
}
```

#### 404 Not Found
```json
{
  "status_code": 404,
  "message": "Topic 'test-topic' not found",
  "error": "resource not found"
}
```

#### 409 Conflict
```json
{
  "status_code": 409,
  "message": "Consumer ID already exists",
  "error": "duplicate resource"
}
```

#### 500 Internal Server Error
```json
{
  "status_code": 500,
  "message": "Error connecting to cluster",
  "error": "connection refused"
}
```

---

## Rate Limiting

Currently, there are no rate limits enforced. For production use, consider implementing rate limiting at the reverse proxy level.

---

## CORS

The API supports CORS for web browser access. Configure CORS settings based on your deployment requirements.

---

## Pagination

Currently, the API does not support pagination. All list endpoints return complete results. For large datasets, consider implementing pagination in future versions.

---

## Versioning

The API is versioned through the URL path (`/api/v1`). Breaking changes will result in a new version (`/api/v2`).

---

## SDK Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

const api = axios.create({
  baseURL: 'http://localhost:8081/api/v1',
  headers: { 'Content-Type': 'application/json' }
});

// List topics
const topics = await api.get('/topics');
console.log(topics.data);

// Create topic
await api.post('/topics', {
  name: 'my-topic',
  partitions: 3,
  replication_factor: 1
});

// Produce message
await api.post('/messages/my-topic', {
  message: 'Hello, Kafka!'
});
```

### Python

```python
import requests

BASE_URL = 'http://localhost:8081/api/v1'

# List topics
response = requests.get(f'{BASE_URL}/topics')
topics = response.json()

# Create topic
requests.post(f'{BASE_URL}/topics', json={
    'name': 'my-topic',
    'partitions': 3,
    'replication_factor': 1
})

# Produce message
requests.post(f'{BASE_URL}/messages/my-topic', json={
    'message': 'Hello, Kafka!'
})
```

### cURL

```bash
# List topics
curl http://localhost:8081/api/v1/topics

# Create topic
curl -X POST http://localhost:8081/api/v1/topics \
  -H "Content-Type: application/json" \
  -d '{"name":"my-topic","partitions":3,"replication_factor":1}'

# Produce message
curl -X POST http://localhost:8081/api/v1/messages/my-topic \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello, Kafka!"}'
```

---

## WebSocket Support

WebSocket support is planned for future releases to enable real-time updates without polling.

---

## OpenAPI Specification

The complete OpenAPI 3.0 specification is available at `docs/openapi.yaml`. You can use tools like Swagger UI or Postman to explore the API interactively.

---

## Support

For API issues or questions:
- GitHub Issues: https://github.com/IBM/openkommander/issues
- Documentation: `/docs` directory