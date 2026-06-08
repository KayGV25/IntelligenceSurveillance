CREATE SCHEMA IF NOT EXISTS camera;

CREATE TABLE IF NOT EXISTS camera.cameras (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    rtsp_url TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN',

    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,

    building_id UUID,
    floor_id UUID,
    zone_id UUID,

    position_x DOUBLE PRECISION,
    position_y DOUBLE PRECISION,
    direction_angle DOUBLE PRECISION,
    fov_angle DOUBLE PRECISION,

    created_by UUID,
    updated_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_cameras_floor_id
ON camera.cameras(floor_id);

CREATE INDEX IF NOT EXISTS idx_cameras_status
ON camera.cameras(status);

CREATE INDEX IF NOT EXISTS idx_cameras_deleted_at
ON camera.cameras(deleted_at);