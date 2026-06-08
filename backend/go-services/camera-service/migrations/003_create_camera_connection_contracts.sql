CREATE TABLE IF NOT EXISTS camera.camera_connection_contracts (
    id UUID PRIMARY KEY,
    camera_id UUID NOT NULL REFERENCES camera.cameras(id),

    discovered_device_id UUID REFERENCES camera.discovered_devices(id),

    connection_type VARCHAR(50) NOT NULL,

    ip_address VARCHAR(100) NOT NULL,
    rtsp_url TEXT,
    onvif_url TEXT,

    username VARCHAR(255),
    encrypted_password TEXT,

    stream_profile_token VARCHAR(255),
    main_stream_url TEXT,
    sub_stream_url TEXT,
    snapshot_url TEXT,

    capabilities JSONB,
    metadata JSONB,

    created_by UUID,
    updated_by UUID,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_camera_connection_contracts_camera_id
ON camera.camera_connection_contracts(camera_id);

CREATE INDEX IF NOT EXISTS idx_camera_connection_contracts_discovered_device_id
ON camera.camera_connection_contracts(discovered_device_id);

CREATE INDEX IF NOT EXISTS idx_camera_connection_contracts_connection_type
ON camera.camera_connection_contracts(connection_type);

CREATE UNIQUE INDEX IF NOT EXISTS ux_camera_connection_contracts_discovered_device_id
ON camera.camera_connection_contracts(discovered_device_id)
WHERE discovered_device_id IS NOT NULL
AND deleted_at IS NULL;