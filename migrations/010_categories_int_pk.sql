-- migrations/010_categories_int_pk.sql
BEGIN;

-- 1. Rename the old hand-typed TEXT id out of the way, add a real
--    auto-increment internal PK. Existing rows get sequential ints
--    backfilled automatically by SERIAL's DEFAULT nextval(...).
ALTER TABLE categories RENAME COLUMN id TO legacy_id;
ALTER TABLE categories ADD COLUMN id SERIAL;

-- 2. Retype products.category_id from TEXT (matching the old
--    categories.legacy_id) to INT (matching the new categories.id).
ALTER TABLE products ADD COLUMN category_id_new INT;

UPDATE products p
SET category_id_new = c.id
FROM categories c
WHERE c.legacy_id = p.category_id;

-- Every product must have resolved to a category; abort if not.
DO $$
DECLARE unresolved INT;
BEGIN
	SELECT COUNT(*) INTO unresolved FROM products WHERE category_id_new IS NULL;
	IF unresolved > 0 THEN
		RAISE EXCEPTION '% products have no matching category after remap', unresolved;
	END IF;
END $$;

ALTER TABLE products ALTER COLUMN category_id_new SET NOT NULL;
ALTER TABLE products DROP COLUMN category_id;
ALTER TABLE products RENAME COLUMN category_id_new TO category_id;

-- 3. Promote the new column to primary key.
ALTER TABLE categories DROP CONSTRAINT categories_pkey;
ALTER TABLE categories ADD PRIMARY KEY (id);

-- 4. Re-add the FK + index that were dropped along with the old
--    TEXT category_id column.
ALTER TABLE products
	ADD CONSTRAINT products_category_id_fkey
	FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;
CREATE INDEX idx_products_category ON products(category_id);

COMMIT;
