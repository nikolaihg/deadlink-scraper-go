# Part 3 – Scalable Dead Link Microservice for Media Workflows

This part evolves the project into a **production-ready, scalable microservice**.

It focuses on:

- Event-driven architecture
- Real-time reporting
- Observability
- Microservice design patterns

## Key Enhancements

### 1. gRPC Service

Expose scanning via gRPC API:

```proto
service LinkScanner {
  rpc ScanURL(ScanRequest) returns (ScanResponse);
}
```

Supports structured, efficient service-to-service communication.

### 2. Kafka Integration

* Emit scan results to a Kafka topic: `link-scan-results`
* Schema:

```json
{
  "scan_id": "uuid",
  "url": "https://tv2.no",
  "timestamp": "2025-06-19T10:00:00Z",
  "dead_links": [
    {"href": "https://...", "status_code": 404}
  ]
}
```

Tools:

* [`segmentio/kafka-go`](https://github.com/segmentio/kafka-go)
* Simulate topic with `Redpanda` locally

### 3. Observability (Prometheus-ready)

* Add `promhttp` endpoint on `/metrics`
* Track:

  * # URLs scanned
  * Avg scan duration
  * # dead links per run

### 4. Deployment Ready

* Build Docker images
* Helm chart or Compose for dev
* Kubernetes-ready manifest

## How It Relates to TV 2

* Mimics their Go + microservice + Kafka stack
* Shows experience with system observability, monitoring, and structured APIs
* Demonstrates how link QA can fit into editorial or media content pipelines
* Optionally integrate OpenSearch or dashboard-style UI for visual QA

## 🛠️ Technologies

| Stack     | Tool                 |
| --------- | -------------------- |
| Language  | Go                   |
| Messaging | Kafka / Redpanda     |
| API       | gRPC                 |
| Metrics   | Prometheus + Grafana |
| Database  | PostgreSQL           |
| Infra     | Docker, Kubernetes   |

## Stretch Ideas

* gRPC-to-HTTP gateway (via grpc-gateway)
* Internal job scheduler
* Rate limiting via Redis
* Screenshot capture of broken pages (headless Chrome)

Back to [Part 2 – Persistent Backend](./part2.md)
Or the [README](./README.md)