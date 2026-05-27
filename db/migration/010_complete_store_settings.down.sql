ALTER TABLE store_settings
DROP COLUMN IF EXISTS notifications,
DROP COLUMN IF EXISTS opening_hours,
DROP COLUMN IF EXISTS average_delivery_max_minutes,
DROP COLUMN IF EXISTS average_delivery_min_minutes,
DROP COLUMN IF EXISTS delivery_radius_km,
DROP COLUMN IF EXISTS allow_pickup,
DROP COLUMN IF EXISTS allow_delivery,
DROP COLUMN IF EXISTS instagram,
DROP COLUMN IF EXISTS brand_color,
DROP COLUMN IF EXISTS short_description,
DROP COLUMN IF EXISTS cnpj;
