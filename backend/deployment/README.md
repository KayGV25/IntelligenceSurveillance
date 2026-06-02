# Advanced Home Surveillance - Local Docker Infrastructure

This folder contains the development Docker environment for backend infrastructure services.

## Services

| Service | Purpose | URL / Port |
|---|---|---|
| PostgreSQL | Relational database | localhost:5432 |
| Redis | Cache, sessions, live state | localhost:6379 |
| Redpanda | Kafka-compatible event bus | localhost:19092 |
| Redpanda Console | Kafka UI | http://localhost:8088 |
| MinIO | S3-compatible object storage | http://localhost:9000 |
| MinIO Console | Object storage UI | http://localhost:9001 |
| Qdrant | Vector database | http://localhost:6333 |

## Start

```bash
cp .env.example .env
docker compose up -d
```

## Check status

```bash
docker compose ps
```

## Stop

```bash
docker compose down
```

## Stop and delete volumes

```bash
docker compose down -v
```

## Useful connection strings

PostgreSQL:

```text
postgresql://ahs_user:ahs_password@localhost:5432/ahs_dev
```

Redis:

```text
localhost:6379
```

Kafka / Redpanda broker from host machine:

```text
localhost:19092
```

Kafka / Redpanda broker from another Docker container in the same Compose network:

```text
redpanda:9092
```

MinIO credentials:

```text
user: minioadmin
password: minioadmin123
```
