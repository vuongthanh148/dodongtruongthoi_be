-- Migration 015: per-order lookup codes and order-bound lookup sessions.
-- Each order gets a random 6-character code the buyer keeps; it replaces the
-- order-id-derived code. Sessions are bound to one order, so a token opens only
-- the order it was verified for. Existing sessions carry no order and are cleared.

ALTER TABLE orders ADD COLUMN IF NOT EXISTS lookup_code TEXT;

-- Backfill each existing row with its own random code. A loop is used so every
-- row draws independently; an uncorrelated subquery would be evaluated once.
DO $$
DECLARE
    alphabet CONSTANT TEXT := 'ABCDEFGHJKMNPQRSTUVWXYZ23456789';
    r        RECORD;
    code     TEXT;
BEGIN
    FOR r IN SELECT id FROM orders WHERE lookup_code IS NULL LOOP
        code := '';
        FOR i IN 1..6 LOOP
            code := code || substr(alphabet, 1 + floor(random() * 31)::int, 1);
        END LOOP;
        UPDATE orders SET lookup_code = code WHERE id = r.id;
    END LOOP;
END $$;

ALTER TABLE orders ALTER COLUMN lookup_code SET NOT NULL;

DELETE FROM lookup_sessions;
ALTER TABLE lookup_sessions
    ADD COLUMN order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE;
