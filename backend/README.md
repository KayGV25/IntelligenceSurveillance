Advanced Home Surveillance — Backend
===================================

Summary
-------
This folder contains the backend services for the Advanced Home Surveillance platform. See the detailed product and engineering specification in [backend/BACKEND_SPECS.md](backend/BACKEND_SPECS.md).

Contents
--------
- **Services**: Go microservices under [backend/go-services](go-services/) (camera-service, map-service, common libs).
- **Java services**: API gateway and auth in [backend/java-services](../java-services).
- **Deployment**: docker-compose for local/dev in [backend/deployment/docker-compose.yml](deployment/docker-compose.yml).

Quick start (local, dev)
------------------------
Prerequisites: `docker`, `docker-compose` (or `docker compose`), Go (1.20+ for building services locally), Java (for java-services if building locally).

Bring up the full dev stack (Postgres, Redis, MinIO, Redpanda, services):

```bash
cd backend/deployment
docker compose up --build
```

The compose file is at [backend/deployment/docker-compose.yml](deployment/docker-compose.yml).

Core backend services and ports (dev compose)
- API Gateway: `8080` (java)
- Auth Service: `8081` (java)
- Camera Service: `8101` (go)
- Map Service: `8105` (go)

Java services (location & notes)
- The Java services are in the repository at `backend/java-services` and contain the API Gateway and Authentication service.
- Key Java modules:
	- `api-gateway` — routes public API traffic, validates JWTs, and proxies requests to internal services. See [backend/java-services/api-gateway](go-services/../java-services/api-gateway).
	- `auth-service` — identity, token issuance, and JWKS exposure. See [backend/java-services/auth-service](go-services/../java-services/auth-service).
- Ports used by compose: API Gateway `8080`, Auth `8081`.
- The gateway expects these environment variables (example snippets in compose): `JWT_ENABLED`, `JWT_ISSUER`, `JWT_JWK_SET_URI`, and service URLs like `CAMERA_SERVICE_URL`, `MAP_SERVICE_URL`.

Building & running Java services locally
- From repository root, build and run with Gradle wrappers inside each service folder. Examples:

Linux / macOS:

```bash
cd backend/java-services/api-gateway
./gradlew build
java -jar build/libs/*jar

cd backend/java-services/auth-service
./gradlew build
java -jar build/libs/*jar
```

Windows (PowerShell / CMD):

```powershell
cd backend/java-services
gradlew.bat :api-gateway:bootRun
```

Note: the compose file builds these services from `../java-services` when bringing up the stack.
Java services API contracts
--------------------------
Below are concise API contracts for the Java services in `backend/java-services`. These are derived from controller mappings and DTO records in the service code.

Auth Service (base: `http://localhost:8081`)
- `GET /api/v1/auth/health`
	- Response: `ApiResponse` wrapper with health payload `{ service: "auth-service", status: "UP", timestamp: "..." }`

- `POST /api/v1/auth/login`
	- Request: `LoginRequest`
		- `{ "email": "user@example.com", "password": "secret" }`
	- Response: `ApiResponse<AuthResponse>`
		- `AuthResponse` fields: `accessToken`, `refreshToken`, `expiresInSeconds`, `tokenType`, `user` (id, email, roles[], permissions[])

- `POST /api/v1/auth/refresh`
	- Request: `RefreshTokenRequest` `{ "refreshToken": "..." }`
	- Response: `ApiResponse<AuthResponse>` (same as login)

- `POST /api/v1/auth/logout`
	- Request: `LogoutRequest` `{ "refreshToken": "..." }`
	- Response: `ApiResponse<String>` — success message `"Logged out"`

- `GET /api/v1/auth/me`
	- Header: `Authorization: Bearer <accessToken>` (optional header accepted by controller)
	- Response: `ApiResponse<MeResponse>`
		- `MeResponse` fields: `id`, `email`, `roles[]`, `permissions[]`

JWKS (public key for token verification)
- `GET /.well-known/jwks.json`
	- Response: standard JWKS JSON used by the gateway and clients to verify JWT signatures. Endpoint exposed by `auth-service`.

API Gateway (base: `http://localhost:8080`)
- `GET /api/v1/gateway/health`
	- Response: `ApiResponse` wrapper with `{ service: "api-gateway", status: "UP", timestamp: "..." }`
- The gateway proxies and routes to internal services (camera, map, alert, etc.) according to configuration. It validates JWTs and exposes the unified public API surface under `/api/v1`.

Notes on `ApiResponse` wrapper
- Many Java controllers return `ApiResponse<T>` (see `backend/java-services/*/src/main/java/com/ahs/common/response`). Typical JSON shape:

```json
{
	"success": true,
	"data": { /* payload */ },
	"error": null
}
```

Next steps I can do for you
- generate an OpenAPI/Swagger spec for `auth-service` (and expose YAML/JSON)
- create example curl commands for each Java endpoint
- produce a Postman collection for the Java services

Environment variables used by compose
- `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_PORT`
- `REDIS_PORT`
- `MINIO_ROOT_USER`, `MINIO_ROOT_PASSWORD`, `MINIO_SNAPSHOTS_BUCKET`, etc.
See the compose file for full list.

Building and running a single Go service locally
------------------------------------------------
From repository root:

```bash
cd backend/go-services/camera-service
go build ./...
./camera-service  # service config expects env vars, see docker-compose for names

cd ../map-service
go build ./...
./map-service
```

API surface (summary)
---------------------
All service routers expose health checks and APIs under `/api/v1` (map-service uses `/api/v1/maps`). Below are the primary routes and request/response contracts inferred from code (DTO structs).

Camera Service (http://localhost:8101)
- `GET /health` — health check
- `POST /api/v1/cameras/discover` — discovery request
	- Request: `{ "method": "CIDR_SCAN", "network_cidr": "192.168.1.0/24", "ports": [80,554], "timeout_ms": 5000 }`
	- Response: `{ "devices": [ { /* discovered device */ } ], "count": 1 }`
- `GET /api/v1/cameras/discovered` — list discovered devices
- `GET /api/v1/cameras/discovered/:id` — discovered device by id
- `POST /api/v1/cameras/discovered/:discoveredDeviceId/connect` — connect discovered device
	- Request: `{ "name": "Cam 1", "description": "", "username": "admin", "password": "pass", "preferred_stream": "main" }`
	- Response: `{ "camera": { /* camera object */ }, "contract": { /* connection contract */ } }`
- `POST /api/v1/cameras` — create camera
	- Request: `CreateCameraRequest` fields: `name` (required), `description`, `rtsp_url`, `latitude`, `longitude`, `building_id`, `floor_id`, `zone_id`, `position_x`, `position_y`, `direction_angle`, `fov_angle`.
	- Response: created camera object (HTTP 201)
- `GET /api/v1/cameras` — list cameras
- `GET /api/v1/cameras/:id` — get camera by id
- `PUT /api/v1/cameras/:id` — update camera (same fields as create, optional)
- `DELETE /api/v1/cameras/:id` — delete camera
- `GET /api/v1/cameras/:id/connection` — get camera connection contract
- `POST /api/v1/cameras/:id/validate-stream` — validate camera stream
- `POST /api/v1/cameras/:id/snapshot` — capture snapshot (returns `SnapshotResponse` with `object_path`)
- `GET /api/v1/cameras/:id/health` — get camera health (`CameraHealthResponse`)
- `GET /api/v1/cameras/:id/stream-info` — get `StreamInfoResponse`

Map Service (http://localhost:8105)
- `GET /health` — health check
- Base prefix: `/api/v1/maps`
- Building endpoints
	- `POST /api/v1/maps/buildings` — create building
		- Request: `{ "name": "HQ", "description": "", "address": "" }` (see DTO in code)
	- `GET /api/v1/maps/buildings` — list
	- `GET /api/v1/maps/buildings/:id` — get
	- `PUT /api/v1/maps/buildings/:id` — update
	- `DELETE /api/v1/maps/buildings/:id` — delete
- Floor endpoints
	- `POST /api/v1/maps/floors` — create floor (requires `building_id`, `name`)
	- `GET /api/v1/maps/floors` — list
	- `GET /api/v1/maps/floors/:id` — get
	- `PUT /api/v1/maps/floors/:id` — update
	- `DELETE /api/v1/maps/floors/:id` — delete
- Zone endpoints
	- `POST /api/v1/maps/zones` — create zone (requires `floor_id`, `name`, `polygon` array of points)
	- `GET /api/v1/maps/zones` — list (optional `floor_id` query)
	- `GET /api/v1/maps/zones/:id` — get
	- `GET /api/v1/maps/floors/:floorId/zones` — list zones for a floor
	- `PUT /api/v1/maps/zones/:id` — update
	- `DELETE /api/v1/maps/zones/:id` — delete

API contracts / DTO references
- Camera DTOs: see [backend/go-services/camera-service/internal/dto](go-services/camera-service/internal/dto) for `CreateCameraRequest`, `UpdateCameraRequest`, discovery and snapshot responses.
- Map DTOs: see [backend/go-services/map-service/internal/dto](go-services/map-service/internal/dto) for building/floor/zone request models.

Error responses
- Services use a consistent error response helper in `backend/go-services/common/response`. Typical shape:

```json
{
	"code": "BAD_REQUEST",
	"reason": "INVALID_CAMERA_REQUEST",
	"message": "...",
	"data": null
}
```

Deployment notes
----------------
- The compose file boots supporting infra (Postgres, Redis, MinIO, Redpanda, Qdrant) and the services used by the platform. Use it for local end-to-end testing.
- For production, replace images with pinned release images, use managed cloud services where appropriate, secure secrets (do not use env files in plain text), and set resource limits.

Where to look next
------------------
- Product & API specs: [backend/BACKEND_SPECS.md](backend/BACKEND_SPECS.md)
- Dev compose: [backend/deployment/docker-compose.yml](deployment/docker-compose.yml)
- Camera service implementation: [backend/go-services/camera-service](go-services/camera-service)
- Map service implementation: [backend/go-services/map-service](go-services/map-service)

