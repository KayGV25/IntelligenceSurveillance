CREATE TABLE IF NOT EXISTS camera.discovered_devices (
    id UUID PRIMARY KEY,

    ip_address VARCHAR(100) NOT NULL,
    mac_address VARCHAR(100),
    hostname VARCHAR(255),

    manufacturer VARCHAR(255),
    model VARCHAR(255),
    firmware_version VARCHAR(255),

    onvif_supported BOOLEAN NOT NULL DEFAULT false,
    rtsp_supported BOOLEAN NOT NULL DEFAULT false,
    http_supported BOOLEAN NOT NULL DEFAULT false,

    http_port INT,
    rtsp_port INT,
    onvif_port INT,

    discovery_method VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'DISCOVERED',

    raw_metadata JSONB,

    discovered_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_discovered_devices_ip_address
ON camera.discovered_devices(ip_address);

CREATE INDEX IF NOT EXISTS idx_discovered_devices_status
ON camera.discovered_devices(status);

CREATE INDEX IF NOT EXISTS idx_discovered_devices_discovered_at
ON camera.discovered_devices(discovered_at);

CREATE UNIQUE INDEX IF NOT EXISTS ux_discovered_devices_ip_address
ON camera.discovered_devices(ip_address);