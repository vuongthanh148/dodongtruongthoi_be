-- Local-dev only: extra fake products to exercise listing pagination ("Xem thêm sản phẩm").
-- Idempotent. Run with:
--   docker compose exec -T postgres psql -U postgres -d dodongtruongthoi < scripts/seed-more-products.sql

INSERT INTO products (
  id, title, subtitle, category_id, base_price, description, meaning,
  requires_size, is_active, sort_order
)
SELECT
  'demo-sp-' || n,
  'Sản phẩm mẫu ' || n,
  'Sản phẩm dữ liệu thử số ' || n,
  (SELECT id FROM categories ORDER BY id OFFSET ((n - 1) % GREATEST((SELECT COUNT(*) FROM categories), 1)) LIMIT 1),
  1000000 + n * 250000,
  'Mô tả thử nghiệm cho sản phẩm mẫu ' || n || '.',
  'Ý nghĩa thử nghiệm.',
  true, true, 100 + n
FROM generate_series(1, 12) AS n
ON CONFLICT (id) DO NOTHING;

INSERT INTO product_sizes (product_id, size_label, size_code, price, sort_order)
SELECT p.id, s.label, s.code, p.base_price + s.bump, s.ord
FROM products p
CROSS JOIN (VALUES
  ('0.8m × 0.6m', 's', 0, 1),
  ('1.2m × 0.8m', 'm', 800000, 2),
  ('1.5m × 1.0m', 'l', 1800000, 3)
) AS s(label, code, bump, ord)
WHERE p.id LIKE 'demo-sp-%'
ON CONFLICT (product_id, size_code) DO NOTHING;
