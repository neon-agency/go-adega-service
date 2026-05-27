ALTER TABLE products
ADD COLUMN cost_cents INTEGER NOT NULL DEFAULT 0 CHECK (cost_cents >= 0);

UPDATE stock_movements sm
SET unit_cost_cents = p.cost_cents
FROM products p
WHERE sm.product_id = p.id
  AND sm.movement_type = 'sale'
  AND sm.unit_cost_cents IS NULL;

