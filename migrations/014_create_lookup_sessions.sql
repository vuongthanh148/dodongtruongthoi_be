-- Migration 014: lookup_sessions table
-- Short-lived sessions granted after a successful order lookup. Only the SHA-256
-- hash of the bearer token is stored; expired rows are rejected by time check.

CREATE TABLE IF NOT EXISTS lookup_sessions (
    token_hash  TEXT PRIMARY KEY,
    phone       TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
