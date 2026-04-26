-- Phase 1: Categories
CREATE TABLE categories (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    description TEXT,
    tone        TEXT DEFAULT 'bronze',
    image_url   TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 2: Products
CREATE TABLE products (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    subtitle        TEXT,
    category_id     TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    badge           TEXT,
    base_price      BIGINT NOT NULL,
    description     TEXT,
    meaning         TEXT,
    default_bg      TEXT NOT NULL DEFAULT 'gold',
    default_frame   TEXT NOT NULL DEFAULT 'bronze',
    bg_tones        TEXT[] NOT NULL DEFAULT '{}',
    frames          TEXT[] NOT NULL DEFAULT '{}',
    zodiac_ids      TEXT[] NOT NULL DEFAULT '{}',
    purpose_place   TEXT[] DEFAULT '{}',
    purpose_use     TEXT[] DEFAULT '{}',
    purpose_avoid   TEXT[] DEFAULT '{}',
    specs           JSONB DEFAULT '{}',
    requires_bg_tone BOOLEAN NOT NULL DEFAULT true,
    requires_frame   BOOLEAN NOT NULL DEFAULT true,
    requires_size    BOOLEAN NOT NULL DEFAULT true,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_active ON products(is_active) WHERE is_active = true;

-- Phase 3: Product Images
CREATE TABLE product_images (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    bg_tone     TEXT,
    frame       TEXT,
    url         TEXT NOT NULL,
    alt_text    TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_product_images_product ON product_images(product_id);
CREATE INDEX idx_product_images_variant ON product_images(product_id, bg_tone, frame);

-- Phase 4: Product Sizes
CREATE TABLE product_sizes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    size_label  TEXT NOT NULL,
    size_code   TEXT NOT NULL,
    price       BIGINT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    UNIQUE(product_id, size_code)
);
CREATE INDEX idx_product_sizes_product ON product_sizes(product_id);

-- Phase 5: Campaigns
CREATE TABLE campaigns (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    description     TEXT,
    discount_type   TEXT NOT NULL,
    discount_value  BIGINT NOT NULL,
    starts_at       TIMESTAMPTZ NOT NULL,
    ends_at         TIMESTAMPTZ NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE campaign_products (
    campaign_id UUID NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    product_id  TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    PRIMARY KEY (campaign_id, product_id)
);

-- Phase 6: Reviews
CREATE TABLE reviews (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id    TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    reviewer_name TEXT NOT NULL,
    rating        INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    body          TEXT,
    is_approved   BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_reviews_product ON reviews(product_id);

-- Phase 7: Wishlists
CREATE TABLE wishlists (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone       TEXT NOT NULL,
    product_id  TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(phone, product_id)
);
CREATE INDEX idx_wishlists_phone ON wishlists(phone);

-- Phase 8: Orders
CREATE TABLE orders (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone           TEXT NOT NULL,
    customer_name   TEXT,
    note            TEXT,
    status          TEXT NOT NULL DEFAULT 'pending_confirm',
    admin_note      TEXT,
    total_amount    BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_orders_phone ON orders(phone);
CREATE INDEX idx_orders_status ON orders(status);

CREATE TABLE order_items (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id        UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id      TEXT NOT NULL REFERENCES products(id),
    product_title   TEXT NOT NULL,
    product_subtitle TEXT,
    size_code       TEXT,
    size_label      TEXT,
    bg_tone         TEXT,
    bg_tone_label   TEXT,
    frame           TEXT,
    frame_label     TEXT,
    quantity        INT NOT NULL DEFAULT 1,
    unit_price      BIGINT NOT NULL,
    variant_image_url TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_order_items_order ON order_items(order_id);

-- Phase 9: Banners
CREATE TABLE banners (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT,
    subtitle    TEXT,
    image_url   TEXT,
    link_url    TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 10: Contact Links
CREATE TABLE contact_links (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    platform    TEXT NOT NULL,
    label       TEXT NOT NULL,
    url         TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Phase 11: Site Settings
CREATE TABLE site_settings (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO site_settings (key, value) VALUES
    ('active_theme', 'default'),
    ('hotline', '0899012288'),
    ('shop_name', 'Do Dong Truong Thoi')
ON CONFLICT (key) DO NOTHING;

-- Phase 12: Admin Users
CREATE TABLE admin_users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    display_name    TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login_at   TIMESTAMPTZ
);
