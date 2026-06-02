# API Gateway

Reactive Spring Cloud Gateway service for the Advanced Home Surveillance backend. The gateway is the public HTTP entry point for backend APIs, applies per-client rate limiting, forwards trace headers, normalizes error responses, and optionally enforces JWT authentication.

## Responsibilities

- Routes `/api/v1/**` traffic to downstream services.
- Applies Redis-backed request rate limiting to proxied routes.
- Permits local development without JWT when `JWT_ENABLED=false`.
- Enforces OAuth2 resource server JWT authentication when `JWT_ENABLED=true`.
- Adds or preserves `X-Trace-Id`, `X-Request-Id`, and `X-Correlation-Id` headers on requests and responses.
- Writes gateway and audit logs for request start, completion, auth failures, authorization failures, rate limiting, and request failures.
- Returns shared `ApiResponse` JSON error bodies for gateway, security, and unhandled errors.

## Routes

| Path | Downstream property | Default target |
| --- | --- | --- |
| `/api/v1/auth/**` | `AUTH_SERVICE_URL` | `http://localhost:8081` |
| `/api/v1/cameras/**` | `CAMERA_SERVICE_URL` | `http://localhost:8101` |
| `/api/v1/maps/**` | `MAP_SERVICE_URL` | `http://localhost:8105` |
| `/api/v1/alerts/**` | `ALERT_SERVICE_URL` | `http://localhost:8106` |
| `/api/v1/incidents/**` | `INCIDENT_SERVICE_URL` | `http://localhost:8107` |
| `/api/v1/tracking/**` | `TRACKING_SERVICE_URL` | `http://localhost:8108` |
| `/api/v1/playback/**` | `PLAYBACK_SERVICE_URL` | `http://localhost:8104` |

The gateway health endpoint is served locally at `/api/v1/gateway/health`.

## Security

JWT validation is controlled by:

```properties
JWT_ENABLED=false
```

When disabled, all proxied routes are permitted. When enabled, these paths remain public:

- `OPTIONS /**`
- `/actuator/health`
- `/actuator/info`
- `/api/v1/gateway/health`
- `/api/v1/auth/**`
- `/v3/api-docs/**`
- `/swagger-ui.html`
- `/swagger-ui/**`
- `/webjars/**`

All other routes require a valid bearer token.

## Configuration

Common environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `SERVER_PORT` | `8080` | Gateway HTTP port. |
| `JWT_ENABLED` | `false` | Enables JWT authentication for protected routes. |
| `REDIS_HOST` | `localhost` | Redis host used by the rate limiter. |
| `REDIS_PORT` | `6379` | Redis port used by the rate limiter. |
| `AUTH_SERVICE_URL` | `http://localhost:8081` | Auth service base URL. |
| `CAMERA_SERVICE_URL` | `http://localhost:8101` | Camera service base URL. |
| `MAP_SERVICE_URL` | `http://localhost:8105` | Map service base URL. |
| `ALERT_SERVICE_URL` | `http://localhost:8106` | Alert service base URL. |
| `INCIDENT_SERVICE_URL` | `http://localhost:8107` | Incident service base URL. |
| `TRACKING_SERVICE_URL` | `http://localhost:8108` | Tracking service base URL. |
| `PLAYBACK_SERVICE_URL` | `http://localhost:8104` | Playback service base URL. |

Rate limiting uses a Redis token bucket configured in `RateLimitConfig` with replenish rate `10` and burst capacity `20`. The key resolver uses the remote IP address and falls back to `unknown`.

## Running Locally

From `backend/java-services`:

```powershell
.\gradlew.bat :api-gateway:bootRun
```

For the Docker Compose development stack, from `backend/deployment`:

```powershell
docker compose up api-gateway redis
```

The service listens on `http://localhost:8080` by default.

## Verification

Run the gateway unit tests from `backend/java-services`:

```powershell
.\gradlew.bat :api-gateway:test
```

Build the service JAR:

```powershell
.\gradlew.bat :api-gateway:bootJar
```

## Useful Endpoints

- Gateway health: `GET /api/v1/gateway/health`
- Actuator health: `GET /actuator/health`
- Actuator info: `GET /actuator/info`
- Gateway actuator: `GET /actuator/gateway`
- OpenAPI JSON: `GET /v3/api-docs`
- Swagger UI: `GET /swagger-ui.html`
