ALTER TABLE camera.cameras
ALTER COLUMN rtsp_url DROP NOT NULL;

ALTER TABLE camera.camera_connection_contracts
ALTER COLUMN rtsp_url DROP NOT NULL;