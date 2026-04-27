-- Migration 002: customer_photos table
-- Stores real-home lifestyle photos for the homepage gallery section

CREATE TABLE IF NOT EXISTS customer_photos (
    id          TEXT PRIMARY KEY,
    image_url   TEXT NOT NULL,
    caption     TEXT,               -- e.g. "Hà Nội · Phòng khách"
    sort_order  INTEGER NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT true,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_customer_photos_is_active ON customer_photos (is_active);
CREATE INDEX IF NOT EXISTS idx_customer_photos_sort_order ON customer_photos (sort_order, created_at DESC);
