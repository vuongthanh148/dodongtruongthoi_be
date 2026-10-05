-- Migration 012: contact_messages table
-- Stores messages submitted through the public contact form; handled flags
-- whether the shop team has followed up.

CREATE TABLE IF NOT EXISTS contact_messages (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    phone       TEXT NOT NULL,
    message     TEXT NOT NULL,
    handled     BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_contact_messages_created_at ON contact_messages (created_at DESC);
