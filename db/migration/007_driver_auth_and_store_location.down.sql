ALTER TABLE store_settings
DROP COLUMN IF EXISTS driver_pickup_radius_meters,
DROP COLUMN IF EXISTS longitude,
DROP COLUMN IF EXISTS latitude;

DROP INDEX IF EXISTS idx_drivers_email;

ALTER TABLE drivers
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS max_active_deliveries,
DROP COLUMN IF EXISTS must_change_password,
DROP COLUMN IF EXISTS password_hash,
DROP COLUMN IF EXISTS email;
