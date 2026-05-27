ALTER TABLE orders
DROP COLUMN IF EXISTS payment_mode;

ALTER TABLE store_settings
DROP COLUMN IF EXISTS accept_delivery_cash,
DROP COLUMN IF EXISTS accept_delivery_card,
DROP COLUMN IF EXISTS accept_delivery_pix,
DROP COLUMN IF EXISTS accept_online_card,
DROP COLUMN IF EXISTS accept_online_pix;
