-- Migration 013: lookup_attempts table
-- Tracks failed order-lookup code attempts per phone. After too many failures
-- the phone is locked until locked_until.

CREATE TABLE IF NOT EXISTS lookup_attempts (
    phone         TEXT PRIMARY KEY,
    failed_count  INT NOT NULL DEFAULT 0,
    locked_until  TIMESTAMPTZ NULL,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
