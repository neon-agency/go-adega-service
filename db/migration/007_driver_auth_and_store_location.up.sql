ALTER TABLE drivers
ADD COLUMN email VARCHAR(160),
ADD COLUMN password_hash TEXT,
ADD COLUMN must_change_password BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN max_active_deliveries INTEGER NOT NULL DEFAULT 1 CHECK (max_active_deliveries > 0),
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE UNIQUE INDEX idx_drivers_email ON drivers(LOWER(email)) WHERE email IS NOT NULL;

ALTER TABLE store_settings
ADD COLUMN latitude NUMERIC(10, 7),
ADD COLUMN longitude NUMERIC(10, 7),
ADD COLUMN driver_pickup_radius_meters INTEGER NOT NULL DEFAULT 150 CHECK (driver_pickup_radius_meters > 0);
