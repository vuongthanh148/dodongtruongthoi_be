-- Local-dev only: give every seeded order without items one fake line item, and sync total_amount.
-- Run: docker compose exec -T postgres psql -U postgres -d dodongtruongthoi < scripts/seed-order-items.sql
INSERT INTO order_items (order_id, product_id, product_title, product_subtitle, size_code, size_label, quantity, unit_price)
SELECT o.id, p.id, p.title, p.subtitle, 'm', '1.2m × 0.8m', 1 + (row_number() OVER (ORDER BY o.created_at) % 2), p.base_price
FROM orders o
JOIN LATERAL (SELECT * FROM products WHERE id NOT LIKE 'demo-sp-%' ORDER BY md5(o.id::text || id) LIMIT 1) p ON true
WHERE NOT EXISTS (SELECT 1 FROM order_items i WHERE i.order_id = o.id);

UPDATE orders o SET total_amount = s.t
FROM (SELECT order_id, SUM(unit_price * quantity) t FROM order_items GROUP BY order_id) s
WHERE s.order_id = o.id;
