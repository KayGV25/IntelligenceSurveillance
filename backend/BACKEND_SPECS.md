# Advanced Home Surveillance — Backend Product & Engineering Specification

## 1. Backend Scope

The backend provides the core server-side platform for an AI-powered NVR/VMS system.

It is responsible for:

* Authentication and authorization
* API Gateway routing
* Camera management
* RTSP stream registration
* Recording and playback orchestration
* Building/floor/map metadata
* Alert management
* Incident management
* Person tracking metadata
* AI event ingestion
* Storage coordination
* Real-time WebSocket event delivery

Technology direction:

* Java Spring Boot: Authentication Service and API Gateway
* Go: Core platform services
* PostgreSQL: persistent relational data
* Redis: live state and caching
* Kafka/Redpanda/NATS: event streaming
* MinIO/S3: video clips, recordings, snapshots
* Qdrant: person embedding/vector search

---

# 2. Backend Product Specification

## 2.1 Product Goals

The backend shall support a complete AI surveillance platform with NVR/VMS functionality.

The backend must allow users to:

* Register and manage IP cameras
* View camera status
* Configure buildings, floors, zones, and camera locations
* Store footage for a configurable retention time
* Search and replay historical footage
* Receive real-time AI alerts
* Track people across cameras
* Manage incidents and evidence
* Support web and mobile clients through APIs

---

# 3. User Roles

## 3.1 Super Admin

Can:

* Manage all users
* Manage all buildings
* Manage all cameras
* Configure retention policies
* Configure AI models
* Access all recordings
* Access all incidents
* View audit logs

## 3.2 Admin

Can:

* Manage cameras
* Manage maps
* Manage users under assigned organization
* Configure alert rules
* View playback

## 3.3 Security Operator

Can:

* View live cameras
* View alerts
* Acknowledge alerts
* Track persons
* Create incidents

## 3.4 Investigator

Can:

* Search historical footage
* Export clips
* View incidents
* Add notes
* Generate evidence packages

## 3.5 Viewer

Can:

* View assigned live cameras
* View assigned maps
* Cannot modify data

---

# 4. Backend Functional Requirements

## 4.1 Authentication & Authorization

### Requirements

The backend shall support:

* User login
* User logout
* JWT access tokens
* Refresh tokens
* Role-based access control
* Permission-based access control
* Password reset
* Account lockout
* Session invalidation
* Multi-device sessions
* Audit logging

### Authentication Service

Implemented in:

```text
Java Spring Boot
```

Responsibilities:

* User identity
* Password hashing
* Token issuing
* Token refreshing
* Token revocation
* Role validation
* Permission validation

### Security Requirements

* Passwords must be hashed using BCrypt or Argon2
* JWT access tokens should be short-lived
* Refresh tokens should be stored securely
* All protected APIs must require authorization
* Admin APIs must require elevated permissions

---

## 4.2 API Gateway

### Requirements

The API Gateway shall:

* Expose a unified public API
* Validate JWT tokens
* Forward requests to internal services
* Enforce rate limits
* Handle CORS
* Provide centralized request logging
* Provide API versioning

Implemented in:

```text
Java Spring Boot
```

Example public API prefix:

```text
/api/v1
```

Example internal service routing:

```text
/api/v1/auth      -> Auth Service
/api/v1/cameras   -> Camera Service
/api/v1/maps      -> Map Service
/api/v1/alerts    -> Alert Service
/api/v1/playback  -> Playback Service
/api/v1/incidents -> Incident Service
/api/v1/tracking  -> Tracking Service
```

---

## 4.3 Camera Management

### Requirements

The backend shall support:

* Add camera
* Edit camera
* Delete camera
* Enable/disable camera
* Assign camera to building/floor/zone
* Store RTSP URL
* Store ONVIF configuration
* Store camera credentials securely
* Store position on map
* Store camera direction and field of view
* Monitor camera health

### Camera Fields

```text
camera_id
name
description
ip_address
rtsp_url
onvif_url
username
encrypted_password
building_id
floor_id
zone_id
status
resolution_width
resolution_height
fps
bitrate
fov_angle
direction_angle
position_x
position_y
created_at
updated_at
```

### Camera Health Status

Supported statuses:

```text
ONLINE
OFFLINE
DEGRADED
RECORDING
NOT_RECORDING
UNKNOWN
```

---

## 4.4 Recording Management

### Requirements

The backend shall support:

* Continuous recording
* Event-based recording
* Hybrid recording
* Per-camera retention policy
* Per-zone retention policy
* Storage quota tracking
* Automatic cleanup
* Clip preservation
* Recording gap detection

### Recording Modes

```text
CONTINUOUS
EVENT_ONLY
HYBRID
DISABLED
```

### Retention Examples

```text
7 days
14 days
30 days
90 days
custom
```

### Recording Metadata

The backend stores metadata, not raw video bytes, in PostgreSQL.

Video files are stored in MinIO/S3.

```text
recording_segment_id
camera_id
start_time
end_time
duration_seconds
storage_path
file_size_bytes
codec
resolution
fps
has_audio
is_preserved
created_at
```

---

## 4.5 Playback Management

### Requirements

The backend shall support:

* Query recordings by camera and time range
* Generate playback manifest
* Seek to timestamp
* Return available footage ranges
* Support multi-camera synchronized playback
* Support clip export
* Support snapshot export
* Support event-based playback

### Playback Query Example

```text
camera_id = CAM-001
start_time = 2026-05-29T10:00:00Z
end_time = 2026-05-29T10:30:00Z
```

Expected result:

```text
List of video segments and playback manifest URL
```

---

## 4.6 Building & Map Management

### Requirements

The backend shall support:

* Create building
* Create floor
* Upload floor plan
* Store map coordinates
* Store walls, doors, rooms, zones
* Store camera placements
* Store restricted zones
* Store camera coverage cone metadata
* Store walkthrough mapping results

### Map Entity Types

```text
CAMERA
WALL
ROOM
DOOR
ZONE
RESTRICTED_ZONE
STAIR
ELEVATOR
ENTRANCE
EXIT
```

### Map Coordinate System

Each floor has its own coordinate space.

```text
origin_x = 0
origin_y = 0
unit = meter or pixel-normalized coordinate
```

Recommended internal format:

```text
normalized coordinates between 0.0 and 1.0
```

This allows floor plans with different image sizes to render consistently.

---

## 4.7 Alert Management

### Requirements

The backend shall support:

* Receive AI alerts
* Store alerts
* Send real-time notifications
* Display alert status
* Assign alert severity
* Allow operators to acknowledge alerts
* Mark alerts as resolved
* Mark alerts as false positive
* Link alerts to incidents

### Alert Types

```text
FIRE
SMOKE
FALL
FIGHTING
INTRUSION
LOITERING
ABANDONED_OBJECT
CROWD_DENSITY
CAMERA_TAMPERING
CAMERA_OFFLINE
PERSON_TRACKED
```

### Alert Severity

```text
LOW
MEDIUM
HIGH
CRITICAL
```

### Alert Status

```text
NEW
ACKNOWLEDGED
INVESTIGATING
RESOLVED
FALSE_POSITIVE
```

---

## 4.8 Incident Management

### Requirements

The backend shall support:

* Create incident from alert
* Create incident manually
* Assign incident to user
* Add notes
* Attach video clips
* Attach snapshots
* Attach person tracking timeline
* Export incident report
* Store evidence package

### Incident Status

```text
OPEN
IN_PROGRESS
RESOLVED
CLOSED
ARCHIVED
```

---

## 4.9 Person Tracking Metadata

### Requirements

The backend shall support:

* Store detected person records
* Store per-camera tracking IDs
* Store global person IDs
* Store person location on map
* Store movement trajectory
* Store camera sightings
* Store confidence score
* Link person sightings to video segments

### Tracking Concepts

```text
local_track_id
global_track_id
camera_id
floor_id
zone_id
position_x
position_y
timestamp
confidence
```

### Product Behavior

When a person is selected:

The backend shall return:

* current position
* current camera
* related cameras
* movement trail
* timeline of sightings
* matching playback clips

---

## 4.10 AI Event Ingestion

### Requirements

The backend shall expose internal APIs or message topics for AI services to publish:

* detections
* tracking updates
* alert candidates
* re-identification results
* model confidence scores

The backend should not run heavy AI inference directly.

AI inference is handled by the AI system.

The backend stores, validates, indexes, and distributes AI results.

---

## 4.11 Notification System

### Requirements

The backend shall support:

* WebSocket notifications
* Mobile push notification integration
* Email notification
* Webhook notification

Notification events:

```text
alert.created
alert.updated
incident.created
camera.offline
tracking.person.updated
recording.failed
```

---

## 4.12 Audit Logging

### Requirements

The backend shall log important actions:

* login
* logout
* failed login
* camera created
* camera deleted
* playback accessed
* clip exported
* incident exported
* user permission changed

Audit records must include:

```text
user_id
action
resource_type
resource_id
ip_address
user_agent
timestamp
metadata
```

---

# 5. Backend Engineering Specification

---

# 5.1 Backend Service Architecture

## Service Overview

```text
Java Services
├── api-gateway
└── auth-service

Go Services
├── camera-service
├── video-ingest-service
├── recording-service
├── playback-service
├── map-service
├── alert-service
├── incident-service
├── tracking-service
├── notification-service
└── storage-service
```

---

# 5.2 Java Services

## 5.2.1 API Gateway

Technology:

```text
Java 21
Spring Boot 3
Spring Cloud Gateway
Spring Security
```

Responsibilities:

* Request routing
* JWT validation
* Rate limiting
* CORS
* Request logging
* API versioning
* Reverse proxy to internal services

Suggested port:

```text
8080
```

---

## 5.2.2 Auth Service

Technology:

```text
Java 21
Spring Boot 3
Spring Security
PostgreSQL
Redis
```

Responsibilities:

* Login
* Refresh token
* Logout
* User management
* Role management
* Permission management
* Session management

Suggested port:

```text
8081
```

---

# 5.3 Go Services

## 5.3.1 Camera Service

Responsibilities:

* Camera CRUD
* Camera grouping
* Camera health metadata
* ONVIF discovery integration
* Camera credential references

Suggested port:

```text
8101
```

Database:

```text
PostgreSQL
```

---

## 5.3.2 Video Ingest Service

Responsibilities:

* Connect to RTSP streams
* Validate camera stream availability
* Publish stream health
* Provide stream input to recorder and AI pipeline
* Handle reconnect logic

Suggested technologies:

```text
Go
FFmpeg
GStreamer
Pion WebRTC
```

Suggested port:

```text
8102
```

---

## 5.3.3 Recording Service

Responsibilities:

* Record stream segments
* Generate HLS/fMP4/MP4 chunks
* Store segments in MinIO/S3
* Write recording metadata
* Enforce retention policies
* Detect recording gaps

Suggested port:

```text
8103
```

---

## 5.3.4 Playback Service

Responsibilities:

* Query recording segments
* Generate playback manifests
* Serve signed playback URLs
* Support synchronized multi-camera playback
* Generate clip exports

Suggested port:

```text
8104
```

---

## 5.3.5 Map Service

Responsibilities:

* Building CRUD
* Floor CRUD
* Floor plan upload metadata
* Map entity CRUD
* Camera placement
* Zone definitions
* Walkthrough mapping data persistence

Suggested port:

```text
8105
```

---

## 5.3.6 Alert Service

Responsibilities:

* Receive alerts from AI pipeline
* Store alerts
* Manage alert lifecycle
* Publish real-time alert events
* Trigger notifications

Suggested port:

```text
8106
```

---

## 5.3.7 Incident Service

Responsibilities:

* Incident CRUD
* Link alerts to incidents
* Evidence package management
* Report generation coordination

Suggested port:

```text
8107
```

---

## 5.3.8 Tracking Service

Responsibilities:

* Store tracking events
* Store global person tracks
* Provide movement timeline
* Query person sightings
* Interface with vector database

Suggested port:

```text
8108
```

---

## 5.3.9 Notification Service

Responsibilities:

* WebSocket event delivery
* Push notification dispatch
* Email notification dispatch
* Webhook delivery

Suggested port:

```text
8109
```

---

## 5.3.10 Storage Service

Responsibilities:

* Generate upload URLs
* Generate download URLs
* Manage MinIO/S3 paths
* Store snapshot and clip references
* Enforce storage policies

Suggested port:

```text
8110
```

---

# 6. Backend Data Storage

## 6.1 PostgreSQL

Primary relational database.

Stores:

* users
* roles
* permissions
* cameras
* buildings
* floors
* maps
* alerts
* incidents
* recording metadata
* tracking metadata
* audit logs

---

## 6.2 Redis

Used for:

* token/session cache
* live camera state
* temporary tracking state
* rate limiting
* pub/sub if needed

---

## 6.3 MinIO/S3

Stores:

* recording segments
* exported clips
* snapshots
* floor plan images
* evidence packages

---

## 6.4 Qdrant

Stores:

* person embeddings
* appearance vectors
* re-identification search index

---

# 7. Event Bus Specification

Recommended:

```text
Kafka or Redpanda
```

For simpler deployment:

```text
NATS
```

## 7.1 Core Topics

```text
camera.status.updated
camera.stream.health
recording.segment.created
recording.failed
ai.detection.created
ai.alert.candidate
alert.created
alert.updated
tracking.person.updated
tracking.person.sighted
incident.created
notification.dispatch
audit.event.created
```

---

# 8. Core Event Schemas

## 8.1 AI Detection Event

```json
{
  "event_id": "evt_001",
  "camera_id": "cam_001",
  "timestamp": "2026-05-29T10:00:00Z",
  "type": "PERSON",
  "bbox": {
    "x": 0.12,
    "y": 0.20,
    "w": 0.18,
    "h": 0.42
  },
  "confidence": 0.94,
  "local_track_id": "track_123"
}
```

---

## 8.2 Alert Created Event

```json
{
  "alert_id": "alert_001",
  "type": "FIRE",
  "severity": "CRITICAL",
  "camera_id": "cam_001",
  "floor_id": "floor_001",
  "zone_id": "zone_001",
  "timestamp": "2026-05-29T10:00:00Z",
  "confidence": 0.91
}
```

---

## 8.3 Tracking Update Event

```json
{
  "global_track_id": "person_001",
  "camera_id": "cam_003",
  "floor_id": "floor_001",
  "zone_id": "zone_002",
  "timestamp": "2026-05-29T10:00:05Z",
  "position": {
    "x": 0.42,
    "y": 0.65
  },
  "confidence": 0.88
}
```

---

# 9. API Specification Overview

## 9.1 Auth APIs

```http
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/auth/me
```

---

## 9.2 User APIs

```http
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/{id}
PUT    /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

---

## 9.3 Camera APIs

```http
GET    /api/v1/cameras
POST   /api/v1/cameras
GET    /api/v1/cameras/{id}
PUT    /api/v1/cameras/{id}
DELETE /api/v1/cameras/{id}

GET    /api/v1/cameras/{id}/health
POST   /api/v1/cameras/discover
```

---

## 9.4 Recording APIs

```http
GET  /api/v1/recordings
GET  /api/v1/recordings/ranges
POST /api/v1/recordings/policies
PUT  /api/v1/recordings/policies/{id}
```

---

## 9.5 Playback APIs

```http
GET  /api/v1/playback/manifest
GET  /api/v1/playback/segments
POST /api/v1/playback/export
GET  /api/v1/playback/export/{id}
```

---

## 9.6 Map APIs

```http
GET    /api/v1/buildings
POST   /api/v1/buildings
GET    /api/v1/buildings/{id}

GET    /api/v1/floors
POST   /api/v1/floors

GET    /api/v1/maps/{floor_id}/entities
POST   /api/v1/maps/{floor_id}/entities
PUT    /api/v1/maps/entities/{id}
DELETE /api/v1/maps/entities/{id}
```

---

## 9.7 Alert APIs

```http
GET  /api/v1/alerts
GET  /api/v1/alerts/{id}
POST /api/v1/alerts/{id}/acknowledge
POST /api/v1/alerts/{id}/resolve
POST /api/v1/alerts/{id}/false-positive
```

---

## 9.8 Incident APIs

```http
GET    /api/v1/incidents
POST   /api/v1/incidents
GET    /api/v1/incidents/{id}
PUT    /api/v1/incidents/{id}
POST   /api/v1/incidents/{id}/notes
POST   /api/v1/incidents/{id}/export
```

---

## 9.9 Tracking APIs

```http
GET  /api/v1/tracking/persons/{global_track_id}
GET  /api/v1/tracking/persons/{global_track_id}/timeline
GET  /api/v1/tracking/persons/{global_track_id}/sightings
POST /api/v1/tracking/search
```

---

# 10. Database Design Overview

## 10.1 Core Tables

```text
users
roles
permissions
user_roles
role_permissions
refresh_tokens
sessions

cameras
camera_credentials
camera_health
buildings
floors
map_entities
zones

recording_segments
recording_policies
playback_exports

alerts
incidents
incident_notes
incident_evidence

person_tracks
person_sightings
tracking_positions

audit_logs
notifications
```

---

# 11. Security Engineering Requirements

## 11.1 API Security

* All external APIs must go through API Gateway
* All protected APIs require JWT
* Internal service APIs should not be publicly exposed
* Service-to-service authentication should be used

## 11.2 Credential Security

Camera passwords must not be stored in plaintext.

Recommended:

```text
AES-GCM encryption
```

Secrets should be stored through:

```text
Vault
Kubernetes Secrets
Cloud Secret Manager
```

## 11.3 Audit Security

The backend must log:

* authentication events
* permission changes
* camera changes
* footage access
* evidence export

---

# 12. Performance Requirements

## 12.1 API Performance

Target response time:

```text
p95 < 300ms for metadata APIs
p95 < 1s for playback manifest queries
```

## 12.2 Real-Time Events

Tracking update delivery:

```text
< 500ms
```

Alert delivery:

```text
< 1 second after backend receives alert
```

## 12.3 Recording

Recording service must tolerate:

* stream reconnects
* packet loss
* temporary storage errors
* camera offline periods

---

# 13. Scalability Requirements

Initial target:

```text
100 cameras
```

Future target:

```text
1000+ cameras
```

The backend shall support:

* horizontal service scaling
* distributed recording workers
* distributed AI workers
* partitioned event topics
* object storage scaling

---

# 14. Deployment Requirements

## 14.1 Development Deployment

Recommended:

```text
Docker Compose
```

Services:

```text
PostgreSQL
Redis
MinIO
Kafka or Redpanda
Qdrant
API Gateway
Auth Service
Go Services
```

## 14.2 Production Deployment

Recommended:

```text
Kubernetes
```

Required:

* service discovery
* config management
* secrets management
* horizontal pod autoscaling
* persistent volumes
* monitoring

---

# 15. Observability Requirements

## 15.1 Logging

Use structured JSON logs.

Fields:

```text
timestamp
service
level
trace_id
user_id
message
metadata
```

## 15.2 Metrics

Track:

* API latency
* error rate
* stream health
* recording failures
* alert count
* notification delivery
* storage usage
* Kafka lag

## 15.3 Tracing

Recommended:

```text
OpenTelemetry
```

## 15.4 Monitoring Stack

Recommended:

```text
Prometheus
Grafana
Loki
Tempo
```

---

# 16. Backend Testing Requirements

## 16.1 Unit Tests

Required for:

* auth logic
* permissions
* retention policy logic
* alert status transitions
* map coordinate validation

## 16.2 Integration Tests

Required for:

* PostgreSQL repositories
* Redis cache
* Kafka producers/consumers
* MinIO upload/download
* API Gateway routing

## 16.3 End-to-End Tests

Required flows:

* login → create camera → stream health update
* AI alert → alert display → acknowledge → incident creation
* recording segment → playback manifest → clip export
* person tracking update → map position update

---

# 17. Six-Month Backend Delivery Plan

## Month 1 — Backend Foundation

Deliver:

* API Gateway
* Auth Service
* RBAC
* PostgreSQL schema
* Redis setup
* Kafka/Redpanda setup
* MinIO setup
* Camera Service basic CRUD

## Month 2 — Camera, Streaming, Recording

Deliver:

* RTSP validation
* camera health monitoring
* recording service
* recording metadata
* retention policy engine
* storage service

## Month 3 — Playback and Mapping

Deliver:

* playback service
* playback manifest API
* clip export API
* building service
* floor service
* map entity service
* camera placement API

## Month 4 — Alert and AI Event Ingestion

Deliver:

* AI event ingestion API/topics
* alert service
* alert workflow
* notification service
* WebSocket delivery
* audit logging

## Month 5 — Tracking and Investigation

Deliver:

* tracking service
* person sightings
* trajectory APIs
* Qdrant integration
* incident service
* evidence package support

## Month 6 — Production Hardening

Deliver:

* service-to-service auth
* performance testing
* retention cleanup optimization
* observability stack
* deployment scripts
* Kubernetes manifests
* CI/CD
* backend documentation

---

# 18. Backend Acceptance Criteria

The backend is considered complete when:

* Users can authenticate and access APIs by role
* Cameras can be registered, edited, and monitored
* RTSP camera health can be checked
* Footage metadata can be stored for configurable retention periods
* Playback APIs return historical footage ranges
* Maps, floors, zones, and camera positions can be managed
* AI detections can be ingested and converted into alerts
* Alerts can be acknowledged, resolved, and linked to incidents
* Person tracking metadata can be stored and queried
* Incidents can be created and exported
* Web and Android clients can receive real-time updates
* All critical actions are audit logged
* Services can run through Docker Compose locally
* Services can be deployed to Kubernetes for production
