CREATE TABLE IF NOT EXISTS map.floors (
    id UUID PRIMARY KEY,
    building_id UUID NOT NULL REFERENCES map.buildings(id),

    name VARCHAR(255) NOT NULL,
    description TEXT,
    floor_number INT,

    width DOUBLE PRECISION,
    height DOUBLE PRECISION,
    unit VARCHAR(50) NOT NULL DEFAULT 'meters',

    color VARCHAR(20) NOT NULL DEFAULT '#64748B',
    opacity DOUBLE PRECISION NOT NULL DEFAULT 1.0,
    icon VARCHAR(100),
    label_visible BOOLEAN NOT NULL DEFAULT true,

    created_by UUID,
    updated_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_floors_building_id
ON map.floors(building_id);

CREATE INDEX IF NOT EXISTS idx_floors_deleted_at
ON map.floors(deleted_at);