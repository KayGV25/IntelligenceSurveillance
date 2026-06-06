ALTER TABLE camera.discovered_devices
ADD COLUMN IF NOT EXISTS device_type VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN';

ALTER TABLE camera.discovered_devices
ADD COLUMN IF NOT EXISTS confidence DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE camera.discovered_devices
ADD COLUMN IF NOT EXISTS detection_reason TEXT;