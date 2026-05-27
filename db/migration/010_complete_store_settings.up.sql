ALTER TABLE store_settings
ADD COLUMN IF NOT EXISTS cnpj VARCHAR(18),
ADD COLUMN IF NOT EXISTS short_description TEXT,
ADD COLUMN IF NOT EXISTS brand_color VARCHAR(20) NOT NULL DEFAULT '#6E1F2C',
ADD COLUMN IF NOT EXISTS instagram VARCHAR(120),
ADD COLUMN IF NOT EXISTS allow_delivery BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS allow_pickup BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN IF NOT EXISTS delivery_radius_km NUMERIC(6, 2) NOT NULL DEFAULT 5,
ADD COLUMN IF NOT EXISTS average_delivery_min_minutes INTEGER NOT NULL DEFAULT 30 CHECK (average_delivery_min_minutes > 0),
ADD COLUMN IF NOT EXISTS average_delivery_max_minutes INTEGER NOT NULL DEFAULT 45 CHECK (average_delivery_max_minutes > 0),
ADD COLUMN IF NOT EXISTS opening_hours JSONB NOT NULL DEFAULT '[
  {"id":"seg","label":"Segunda","open":true,"from":"17:00","to":"23:00"},
  {"id":"ter","label":"Terça","open":true,"from":"17:00","to":"23:00"},
  {"id":"qua","label":"Quarta","open":true,"from":"17:00","to":"23:00"},
  {"id":"qui","label":"Quinta","open":true,"from":"17:00","to":"00:00"},
  {"id":"sex","label":"Sexta","open":true,"from":"17:00","to":"01:00"},
  {"id":"sab","label":"Sábado","open":true,"from":"14:00","to":"01:00"},
  {"id":"dom","label":"Domingo","open":false,"from":"14:00","to":"22:00"}
]'::jsonb,
ADD COLUMN IF NOT EXISTS notifications JSONB NOT NULL DEFAULT '{
  "sound_new": true,
  "sound_confirm": true,
  "browser": true,
  "daily": true,
  "weekly": true,
  "low_stock": true,
  "customer_msg": true
}'::jsonb;

UPDATE store_settings
SET short_description = COALESCE(short_description, 'Delivery de bebidas'),
    brand_color = COALESCE(NULLIF(brand_color, ''), '#6E1F2C'),
    delivery_radius_km = CASE WHEN delivery_radius_km <= 0 THEN 5 ELSE delivery_radius_km END,
    average_delivery_min_minutes = CASE WHEN average_delivery_min_minutes <= 0 THEN 30 ELSE average_delivery_min_minutes END,
    average_delivery_max_minutes = CASE WHEN average_delivery_max_minutes <= 0 THEN 45 ELSE average_delivery_max_minutes END;
