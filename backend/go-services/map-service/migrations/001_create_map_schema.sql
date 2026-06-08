CREATE SCHEMA IF NOT EXISTS map;

CREATE TABLE IF NOT EXISTS map.buildings (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,

    address TEXT,

    color VARCHAR(20) NOT NULL DEFAULT '#2563EB',
    opacity DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    icon VARCHAR(100),
    label_visible BOOLEAN NOT NULL DEFAULT true,

    created_by UUID,
    updated_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_buildings_deleted_at
ON map.buildings(deleted_at);