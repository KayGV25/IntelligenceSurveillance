Go Services
===========

Summary
-------
This folder contains the Go microservices that implement core platform functionality for the Advanced Home Surveillance backend: camera management, map/building management, discovery, and shared libraries.

Structure
---------
- `camera-service` — Camera lifecycle, discovery, connection contracts, snapshots, stream validation.
- `map-service` — Building, floor, and zone management and map queries.
- `common` — Shared modules used by Go services (response helpers, security, observability, errors).

Responsibilities
----------------
- Expose service HTTP APIs (Gin) and health endpoints.
- Provide request validation and consistent JSON `ApiResponse` bodies.
- Integrate with Postgres, Redis, MinIO, and Redpanda as configured in the compose file.

Routes & Controllers
--------------------
Routes are registered in each service under `internal/router/router.go`.

Camera Service (default: `http://localhost:8101`)
- `GET /health`
- `/api/v1` prefix:
  - `POST /api/v1/cameras/discover` — discover cameras (CIDR scan)
  - `GET /api/v1/cameras/discovered` — list discovered devices
  - `GET /api/v1/cameras/discovered/:id`
  - `POST /api/v1/cameras/discovered/:discoveredDeviceId/connect`
  - `POST /api/v1/cameras` — create camera
  - `GET /api/v1/cameras` — list cameras
  - `GET /api/v1/cameras/:id`
  - `PUT /api/v1/cameras/:id`
  - `DELETE /api/v1/cameras/:id`
  - `GET /api/v1/cameras/:id/connection`
  - `POST /api/v1/cameras/:id/validate-stream`
  - `POST /api/v1/cameras/:id/snapshot`
  - `GET /api/v1/cameras/:id/health`
  - `GET /api/v1/cameras/:id/stream-info`

Map Service (default: `http://localhost:8105`)
- `GET /health`
- `/api/v1/maps` prefix:
  - Buildings: `POST /buildings`, `GET /buildings`, `GET /buildings/:id`, `PUT /buildings/:id`, `DELETE /buildings/:id`
  - Floors: `POST /floors`, `GET /floors`, `GET /floors/:id`, `PUT /floors/:id`, `DELETE /floors/:id`
  - Zones: `POST /zones`, `GET /zones`, `GET /zones/:id`, `PUT /zones/:id`, `DELETE /zones/:id`, `GET /floors/:floorId/zones`

API Contracts / DTOs
--------------------
- Camera DTOs: `backend/go-services/camera-service/internal/dto` — `CreateCameraRequest`, `UpdateCameraRequest`, discovery and snapshot response DTOs.
- Map DTOs: `backend/go-services/map-service/internal/dto` — building, floor, zone DTOs (see `CreateBuildingRequest`, `CreateFloorRequest`, `CreateZoneRequest`).

Configuration (env vars used by services)
---------------------------------------
Services read configuration from the environment (see `backend/deployment/docker-compose.yml` for examples):
- Database: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- MinIO: `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_SNAPSHOTS_BUCKET`
- Redpanda/Kafka brokers: `REDPANDA_BROKERS`
- Service ports: `APP_PORT` (e.g. `8101`, `8105`)

Building & Running Locally
--------------------------
From repository root you can run individual services for development.

Run camera service directly (Go must be installed):

```bash
cd backend/go-services/camera-service
# run the server
go run ./cmd/server
```

Run map service:

```bash
cd backend/go-services/map-service
go run ./cmd/server
```

Build a binary:

```bash
cd backend/go-services/camera-service
go build -o bin/camera-service ./cmd/server
./bin/camera-service
```

Running with Docker Compose (recommended for end-to-end)
-------------------------------------------------------
Use the development compose in `backend/deployment/docker-compose.yml`. To bring up only Go services (plus required infra):

```powershell
cd backend/deployment
docker compose up --build postgres redis minio redpanda camera-service map-service
```

Verification
------------
- Run unit tests for a service:

```bash
cd backend/go-services/camera-service
go test ./...
```

Useful files & locations
------------------------
- Camera routes: `backend/go-services/camera-service/internal/router/router.go`
- Map routes: `backend/go-services/map-service/internal/router/router.go`
- DTOs: `backend/go-services/*/internal/dto`
- Handlers: `backend/go-services/*/internal/handler`
- Shared response helpers: `backend/go-services/common/response`

Notes
-----
- The Go services use Gin for HTTP routing and a consistent JSON response helper from `common/response` to keep payloads uniform across services.
- If you want, I can generate OpenAPI specs for the Go services by parsing router/DTO files and producing a YAML/JSON spec.

Useful endpoints
----------------
- Camera health: `GET http://localhost:8101/health`
- Map health: `GET http://localhost:8105/health`

---
Generated README matching the Java services README style for `go-services`.
