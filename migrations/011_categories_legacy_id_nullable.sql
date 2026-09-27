-- Allow legacy_id to be NULL for categories created after migration 010
ALTER TABLE categories ALTER COLUMN legacy_id DROP NOT NULL;
