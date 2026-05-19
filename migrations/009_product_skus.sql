CREATE TABLE IF NOT EXISTS product_skus (
  id          TEXT PRIMARY KEY,
  product_id  TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  size_code   TEXT,
  attrs       JSONB NOT NULL DEFAULT '{}',
  price       BIGINT NOT NULL,
  sort_order  INT NOT NULL DEFAULT 0,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS product_skus_product_id_idx ON product_skus(product_id);
