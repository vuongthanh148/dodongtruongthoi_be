-- Drop old product_images (has bg_tone/frame cols + UUID PK)
DROP TABLE IF EXISTS product_images CASCADE;

-- Central image store: upload once, reuse anywhere
CREATE TABLE images (
  id                   TEXT PRIMARY KEY,
  name                 TEXT NOT NULL DEFAULT '',
  url                  TEXT NOT NULL,
  cloudinary_public_id TEXT NOT NULL,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Product → image join (TEXT PK, no hardcoded variant cols)
CREATE TABLE product_images (
  id          TEXT PRIMARY KEY,
  product_id  TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  image_id    TEXT NOT NULL REFERENCES images(id) ON DELETE CASCADE,
  sort_order  INT NOT NULL DEFAULT 0,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Variant attrs per product_image (one value per key, enforced by PK)
CREATE TABLE product_image_attrs (
  product_image_id TEXT NOT NULL REFERENCES product_images(id) ON DELETE CASCADE,
  attr_key         TEXT NOT NULL,
  attr_value       TEXT NOT NULL,
  PRIMARY KEY (product_image_id, attr_key)
);
CREATE INDEX ON product_image_attrs(product_image_id);

-- Remove hardcoded variant cols from order_items; add generic selected_attrs JSONB
ALTER TABLE order_items
  DROP COLUMN IF EXISTS bg_tone,
  DROP COLUMN IF EXISTS bg_tone_label,
  DROP COLUMN IF EXISTS frame,
  DROP COLUMN IF EXISTS frame_label,
  ADD COLUMN IF NOT EXISTS selected_attrs JSONB NOT NULL DEFAULT '{}';

-- Remove hardcoded variant cols; add extensible JSONB ones
ALTER TABLE products
  DROP COLUMN IF EXISTS default_bg,
  DROP COLUMN IF EXISTS default_frame,
  DROP COLUMN IF EXISTS bg_tones,
  DROP COLUMN IF EXISTS frames,
  DROP COLUMN IF EXISTS requires_bg_tone,
  DROP COLUMN IF EXISTS requires_frame,
  ADD COLUMN IF NOT EXISTS variant_options JSONB NOT NULL DEFAULT '[]',
  ADD COLUMN IF NOT EXISTS default_variant JSONB NOT NULL DEFAULT '{}';
