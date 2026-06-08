CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS map.zones (
    id UUID PRIMARY KEY,
    floor_id UUID NOT NULL REFERENCES map.floors(id),

    name VARCHAR(255) NOT NULL,
    description TEXT,

    zone_type VARCHAR(50) NOT NULL DEFAULT 'GENERAL',
    monitored BOOLEAN NOT NULL DEFAULT true,

    geometry geometry(POLYGON) NOT NULL,

    color VARCHAR(20) NOT NULL DEFAULT '#22C55E',
    opacity DOUBLE PRECISION NOT NULL DEFAULT 0.35,
    icon VARCHAR(100),
    label_visible BOOLEAN NOT NULL DEFAULT true,

    created_by UUID,
    updated_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_zones_floor_id
ON map.zones(floor_id);

CREATE INDEX IF NOT EXISTS idx_zones_deleted_at
ON map.zones(deleted_at);

CREATE INDEX IF NOT EXISTS idx_zones_geometry
ON map.zones
USING GIST (geometry);