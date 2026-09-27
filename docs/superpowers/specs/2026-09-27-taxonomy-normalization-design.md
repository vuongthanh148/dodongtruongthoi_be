# Taxonomy Normalization Design

Date: 2026-09-27
Status: Approved, pending implementation plan

## Problem

`categories.id`, `products.id`, `images.id` and `product_images.id` are
`TEXT PRIMARY KEY`. For categories and products specifically, the id is
**hand-typed by the admin** in an "ID (create only)" form field
(`dodongtruongthoi_fe/src/app/admin/categories/page.tsx`,
`dodongtruongthoi_fe/src/app/admin/products/[id]/page.tsx:422`) with no
slugify, no uniqueness preview, and no naming convention — the admin
mixes Vietnamese and English/latin text freely. A collision surfaces as
a raw 500 from the DB unique constraint.

Separately, several taxonomy fields have **no backing table at all**:

- `products.zodiac_ids`, `purpose_place`, `purpose_use`, `purpose_avoid`
  are `TEXT[]` columns. The valid option list for zodiac and purpose
  place is hardcoded in frontend TypeScript
  (`dodongtruongthoi_fe/src/lib/data.ts`: `ZODIAC`,
  `DEFAULT_PLACE_LABELS`) — adding or renaming an option requires a
  frontend code deploy. `purpose_use`/`purpose_avoid` have no fixed
  list at all; the admin form is a raw comma-separated text input, so
  the same concept ends up stored as different strings across products
  (e.g. "phong thủy" vs "phong-thuy").
- `product_image_attrs` (bg_tone/frame variant data) is a generic EAV
  table (`attr_key TEXT, attr_value TEXT`) with no validation. Vietnamese
  display labels for known values are resolved through an ad-hoc string
  convention in `site_settings` (keys like `bg_tone_label.gold`) — a
  flat key-value table doing double duty as a label dictionary.

Data volume is tiny (4 products, 7 categories at review time), so a
single-shot backfill migration is low risk.

## Goals

1. Categories and products get a real auto-increment integer primary
   key. Existing hand-typed ids are preserved as a `legacy_id` column
   and public URLs keep working via a `slug` column.
2. Zodiac, purpose place, purpose use, purpose avoid, background tone,
   and frame style all become admin-manageable lookup tables (int id,
   stable `code`, Vietnamese `name_vi`) with CRUD in the admin UI —
   no more hardcoded frontend arrays, no more comma-separated free text
   for purpose tags.
3. `product_image_attrs` keeps its EAV shape (still supports arbitrary
   future `attr_key`s) but `attr_value` for known keys (`bg_tone`,
   `frame`) is validated against the new lookup tables' `code`, and the
   admin UI presents a dropdown instead of free text. The
   `site_settings` label-key hack is retired in favor of reading
   `name_vi` from the lookup table.
4. No behavior change to storefront URLs, and no change to tables that
   are out of scope (see below).

## Out of scope

`campaigns`, `orders`, `reviews`, `wishlists`, `banners`,
`contact_links`, `admin_users` — already UUID PK, and their text fields
(banner titles, campaign names, review bodies) are genuine free-form
content, not enum-like taxonomy data.

`images.id` / `product_images.id` — already auto-generated (UUID
string via `uuid.NewString()` in `image_usecase.go`), not hand-typed,
not part of this problem.

## Schema design

### `categories`

```sql
ALTER TABLE categories RENAME COLUMN id TO legacy_id;
ALTER TABLE categories ADD COLUMN id SERIAL;
-- slug column already exists (TEXT UNIQUE) — reused as-is, no change
-- id becomes new PK after FK repoint (see migration plan)
```

### `products`

```sql
ALTER TABLE products RENAME COLUMN id TO legacy_id;
ALTER TABLE products ADD COLUMN id SERIAL;
ALTER TABLE products ADD COLUMN slug TEXT UNIQUE; -- backfilled from legacy_id, 1:1 copy
```

### New lookup tables

Common shape:

```sql
CREATE TABLE <name> (
  id          SERIAL PRIMARY KEY,
  code        TEXT NOT NULL UNIQUE,   -- stable machine key, e.g. 'ty', 'gold'
  name_vi     TEXT NOT NULL,          -- display label, e.g. 'Tý', 'Nền Vàng'
  sort_order  INT NOT NULL DEFAULT 0,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Tables: `zodiac_signs` (+ `years TEXT` column), `purpose_places`,
`purpose_use_tags`, `purpose_avoid_tags`, `bg_tones`, `frame_styles`.

Seed data:
- `zodiac_signs`: 12 rows from `ZODIAC` in `data.ts`
- `purpose_places`: 5 rows from `DEFAULT_PLACE_LABELS`
- `bg_tones`: 4 rows from `BG_TONES`
- `frame_styles`: 4 rows from `FRAME_STYLES`
- `purpose_use_tags` / `purpose_avoid_tags`: derived by splitting and
  de-duplicating the existing comma-separated values across all
  current products (small dataset, done by the migration script)

### Join tables

```sql
CREATE TABLE product_zodiac (
  product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  zodiac_id  INT NOT NULL REFERENCES zodiac_signs(id) ON DELETE CASCADE,
  PRIMARY KEY (product_id, zodiac_id)
);
-- same pattern for product_purpose_place, product_purpose_use,
-- product_purpose_avoid
```

Old `zodiac_ids`, `purpose_place`, `purpose_use`, `purpose_avoid`
`TEXT[]` columns are dropped after data is migrated into the join
tables.

### FK retyping

Every column referencing `products.id` / `categories.id` as TEXT gets
retyped to INT against the new serial id:
`products.category_id`, `product_images.product_id`,
`product_sizes.product_id`, `campaign_products.product_id`,
`reviews.product_id`, `wishlists.product_id`, `order_items.product_id`.

### `product_image_attrs`

No schema change. Validation moves into the Go usecase layer: when
`attr_key = 'bg_tone'`, `attr_value` must equal a `bg_tones.code`; when
`attr_key = 'frame'`, must equal a `frame_styles.code`. Other
`attr_key`s remain unvalidated (EAV extensibility preserved).

## Migration plan

Single SQL migration file (`010_normalize_taxonomies.sql`), one
transaction:

1. Create all 6 new lookup tables + 4 join tables.
2. Seed lookup tables from the hardcoded frontend constants (values
   copied into the migration as literal INSERTs) plus derived tag
   values scanned from existing `purpose_use`/`purpose_avoid` arrays.
3. Add `categories.id SERIAL` / `products.id SERIAL` / `products.slug`,
   backfill `products.slug` from the current TEXT id.
4. For every FK column listed above: add new INT column, populate by
   joining old TEXT id → new SERIAL id, drop old TEXT column, rename
   new column into place, add FK constraint.
5. Migrate `products.zodiac_ids`/`purpose_place`/`purpose_use`/
   `purpose_avoid` array contents into the 4 join tables by matching
   against lookup `code`, then drop the 4 array columns.
6. Migrate existing `product_image_attrs` bg_tone/frame free-text
   values to the matching lookup `code` (should already match 1:1
   since the FE hardcoded ids were already used as attr_value).
7. Rename `categories.id`/`products.id` old TEXT columns to
   `legacy_id`, promote new SERIAL columns to PRIMARY KEY.

## Backend changes

- `internal/domain`: `Category.ID`, `Product.ID` become `int`; add
  `Product.Slug string`; new domain types + repository interfaces for
  the 6 lookup tables.
- `internal/repository/postgres`: update category/product repos for
  int PK + slug lookups (storefront routes by slug); new repos for
  lookup tables (simple CRUD, mirrors `category_repo.go`).
- `internal/delivery/http/handler`: new `admin_zodiac_handler.go`,
  `admin_purpose_handler.go` (or one handler per table — match existing
  granularity), `admin_bg_tone_handler.go`, `admin_frame_style_handler.go`;
  public read endpoints for the storefront to fetch option lists.
- `internal/usecase`: bg_tone/frame `attr_value` validation on
  product-image-attr writes.

## Frontend changes

- `src/lib/types.ts`: `Category.id`/`Product.id` → `number`; add
  `Product.slug`.
- Routing: product/category pages already use the string route param
  as a lookup key — switch that lookup to query by `slug` instead of
  `id` (URLs unchanged: `/products/vinh-quy-bai-to` still resolves,
  just via a `slug` column lookup instead of the old TEXT PK).
- `src/lib/data.ts`: remove `ZODIAC`, `BG_TONES`, `FRAME_STYLES`,
  `DEFAULT_PLACE_LABELS` hardcoded arrays (or keep as typed fallback
  only, matching the existing `CATEGORIES` fallback-data comment
  pattern) — options now fetched from the new public API endpoints.
- New admin pages: `/admin/zodiac`, `/admin/purpose-places`,
  `/admin/purpose-tags` (use/avoid can share one page with a type
  toggle), `/admin/bg-tones`, `/admin/frame-styles` — same list/create/
  edit shape as the existing `/admin/categories` page.
- `admin/products/[id]/page.tsx`: Zodiac/Purpose Place selects switch
  from `ZODIAC`/`DEFAULT_PLACE_LABELS` constants to API-fetched
  options; Purpose Use/Avoid switch from comma-separated text input to
  multi-select (existing options + "create new" inline, matching the
  Zodiac field's `react-select` pattern already in the file).
- Product edit form's bg_tone/frame variant pickers switch from free
  text to dropdowns sourced from the new lookup endpoints.

## Testing

- Migration: run against a copy of the current dev DB, verify row
  counts match pre/post for products/categories/product_images, verify
  every old TEXT id has a resolvable new INT id, verify no orphaned FK
  after the swap.
- Backend: repo-level tests for new lookup CRUD; usecase test for
  bg_tone/frame `attr_value` validation (accepts known code, rejects
  unknown).
- Frontend: admin CRUD pages get the same manual smoke-test pass done
  for the existing admin review (login → create/edit/delete a row in
  each new lookup page); storefront product/category pages verified to
  still resolve by their existing slugs post-migration.

## Risks

- FK retyping touches 7 columns across 6 tables in one migration —
  mistakes here are the highest-risk part. Mitigated by tiny data
  volume (safe to hand-verify every row) and running in a single
  transaction (all-or-nothing).
- `product_image_attrs` validation is app-layer, not a DB constraint —
  a direct DB write could still insert an invalid code. Acceptable
  given the EAV table intentionally supports arbitrary future keys; a
  DB-level `CHECK` isn't feasible without hardcoding the key list back
  in, which defeats the EAV extensibility goal from migration 008.
