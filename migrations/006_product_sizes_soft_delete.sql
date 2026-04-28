ALTER TABLE product_sizes ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT true;
CREATE INDEX IF NOT EXISTS idx_product_sizes_active ON product_sizes(product_id, is_active);
