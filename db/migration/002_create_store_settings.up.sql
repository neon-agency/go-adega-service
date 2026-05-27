CREATE TABLE store_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    store_name VARCHAR(120) NOT NULL DEFAULT 'Adega Flow',
    phone VARCHAR(32) NOT NULL,
    whatsapp VARCHAR(32),
    address_street VARCHAR(160) NOT NULL,
    address_number VARCHAR(30) NOT NULL,
    address_complement VARCHAR(120),
    address_neighborhood VARCHAR(120) NOT NULL,
    address_city VARCHAR(120) NOT NULL,
    address_state VARCHAR(2) NOT NULL,
    address_zip_code VARCHAR(12) NOT NULL,
    delivery_fee_cents INTEGER NOT NULL DEFAULT 0 CHECK (delivery_fee_cents >= 0),
    free_delivery_from_cents INTEGER NOT NULL DEFAULT 0 CHECK (free_delivery_from_cents >= 0),
    min_order_cents INTEGER NOT NULL DEFAULT 0 CHECK (min_order_cents >= 0),
    is_open BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO store_settings (
    store_name,
    phone,
    whatsapp,
    address_street,
    address_number,
    address_neighborhood,
    address_city,
    address_state,
    address_zip_code,
    delivery_fee_cents,
    free_delivery_from_cents,
    min_order_cents
) VALUES (
    'Adega Flow',
    '5511931153811',
    '5511931153811',
    'Rua da Adega',
    '100',
    'Centro',
    'Sao Paulo',
    'SP',
    '01000-000',
    790,
    12000,
    0
);

