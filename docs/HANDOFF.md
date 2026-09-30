# Work handoff — backend view (2026-09-30)

Full handoff (frontend + design work, setup steps, open items): `dodongtruongthoi_fe/docs/HANDOFF.md`.

## Backend state
- Go API + Postgres via `docker compose up -d`. Compose user/password `postgres`/`postgres`, DB `dodongtruongthoi`, API on :8080 (`PORT: "8080"` is set in compose), Postgres on :5432. `.env` is not committed — copy it.
- Migrations 000–011 are embedded and applied at startup (`schema_migrations`). Compose initdb also applies 001 → if the API says `relation "categories" already exists`, insert the `001_create_categories_products_orders.sql` row into `schema_migrations`.
- Recent backend commits: category slug generation, products→categories int FK via slug join, partial product update fix.

## Local fake data (dev only, all idempotent)
```
docker compose exec -T postgres psql -U postgres -d dodongtruongthoi < scripts/seed-test-data.sql
docker compose exec -T postgres psql -U postgres -d dodongtruongthoi < scripts/seed-more-products.sql   # demo-sp-1..12
docker compose exec -T postgres psql -U postgres -d dodongtruongthoi < scripts/seed-order-items.sql     # items + totals
```
Note: `seed.sql` / `seed-test-data.sql` predate migrations 008–011 in places; the two newer scripts use the current schema.

## Open
- Categories have no `image_url` in the DB, so the frontend `/categories` tiles show placeholder art.
- The `docker-compose.yml` and `.claude/settings.json` changes in this commit were pending in the working tree before this session's handoff.
