-- Categories and products sorted by sort_order on every list call
CREATE INDEX IF NOT EXISTS idx_categories_sort ON categories(sort_order, created_at);
CREATE INDEX IF NOT EXISTS idx_categories_active_sort ON categories(is_active, sort_order);

-- Products filtered by is_active + sort_order (existing idx_products_active only covers is_active)
CREATE INDEX IF NOT EXISTS idx_products_active_sort ON products(is_active, sort_order, created_at);

-- Banners: always fetched by is_active + sort_order
CREATE INDEX IF NOT EXISTS idx_banners_active_sort ON banners(is_active, sort_order);

-- Campaigns: date range + active status filter
CREATE INDEX IF NOT EXISTS idx_campaigns_active_dates ON campaigns(is_active, starts_at, ends_at);

-- Order items: product-based lookup (e.g. "how many times was this product ordered?")
CREATE INDEX IF NOT EXISTS idx_order_items_product ON order_items(product_id);

-- Reviews: approved reviews per product (public listing filters is_approved = true)
CREATE INDEX IF NOT EXISTS idx_reviews_product_approved ON reviews(product_id, is_approved);

-- Customer photos: active + sort
CREATE INDEX IF NOT EXISTS idx_customer_photos_active_sort ON customer_photos(is_active, sort_order);
