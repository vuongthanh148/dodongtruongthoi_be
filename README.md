# Đồ Đồng Trường Thơi — Backend

Go REST API for Vietnamese bronze craft e-commerce. Handles catalog, campaigns, orders, reviews, wishlist sync, admin CMS, JWT auth, and site settings.

**Status:** Phase 5 Complete ✅ All endpoints implemented and tested, PostgreSQL required, JWT auth enabled, optional Cloudinary image upload.

---

## Quick Start

### Local Run (No Docker)
```bash
cp .env.example .env
# edit .env
go run ./cmd/server
```

Server: http://localhost:8080

---

## Docker Setup

### Prerequisites
- Docker
- Docker Compose

### Local Development
```bash
# From BE repo root
cd /Users/stephen/Documents/Projects/dodongtruongthoi/dodongtruongthoi_be

# Start postgres + backend
docker compose up -d

# View logs
docker compose logs -f backend

# Stop services
docker compose down

# Reset DB volume (destructive)
docker compose down -v
docker compose up -d
```

Backend: http://localhost:8080

Database:
- host: localhost
- port: 5432
- user: postgres
- password: postgres
- db: dodongtruongthoi

Notes:
- Migrations are mounted from `./migrations` to `/docker-entrypoint-initdb.d`.
- `DATABASE_URL` is overridden in compose to use the internal hostname `postgres`.

### Production Build
```bash
# From BE repo root
docker build -t dodongtruongthoi-api:1.0 .
```

---

## API Overview

### Public
- `GET /api/v1/categories`
- `GET /api/v1/categories/{id}`
- `GET /api/v1/products?category=&sort=&limit=&offset=`
- `GET /api/v1/products/{id}`
- `GET /api/v1/products/{id}/reviews`
- `POST /api/v1/products/{id}/reviews`
- `POST /api/v1/orders`
- `GET /api/v1/orders?phone=`
- `GET /api/v1/orders/{id}`
- `GET /api/v1/wishlists?phone=`
- `POST /api/v1/wishlists`
- `DELETE /api/v1/wishlists/{phone}/{productId}`
- `GET /api/v1/banners`
- `GET /api/v1/contacts`
- `GET /api/v1/settings`

### Admin (Bearer JWT)
- `POST /api/v1/admin/login`
- Product CRUD + image upload + sizes
- Category CRUD
- Campaign CRUD + campaign-product mapping
- Banner CRUD
- Contact CRUD
- `GET /api/v1/admin/orders?status=`
- `GET /api/v1/admin/orders/{id}`
- `PUT /api/v1/admin/orders/{id}/status`
- Review moderation endpoints
- Settings read/update

---

## Environment Variables

Copy and edit:
```bash
cp .env.example .env
```

Variables:
- `APP_NAME` optional app label
- `APP_ENV` `development` or `production`
- `PORT` API port, default `8080`
- `JWT_SECRET` required in production
- `DATABASE_URL` required, PostgreSQL connection string
- `CLOUDINARY_CLOUD_NAME` optional
- `CLOUDINARY_API_KEY` optional
- `CLOUDINARY_API_SECRET` optional

Example:
```bash
APP_NAME=dodongtruongthoi-be
APP_ENV=development
PORT=8080
JWT_SECRET=dev-secret-change-in-production
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dodongtruongthoi
CLOUDINARY_CLOUD_NAME=
CLOUDINARY_API_KEY=
CLOUDINARY_API_SECRET=
```

---

## Development Commands

```bash
# Run API
go run ./cmd/server

# Build
go build ./...

# Tests
go test ./...
```

---

## Project Structure

```text
cmd/server/main.go

internal/
  config/                        # Env config
  delivery/http/                 # Router + handlers + middleware
  domain/                        # Entities + repository interfaces
  infrastructure/database/       # DB connection
  infrastructure/storage/        # Cloudinary uploader
  repository/postgres/           # PostgreSQL repositories
  usecase/                       # Business rules

migrations/                      # Schema + bootstrap SQL
scripts/                         # Seed helpers
pkg/response/                    # JSON response helpers
```

---

## Roadmap

| Phase | Status | Scope |
|-------|--------|-------|
| 1 | ✅ | Health, structure, conventions |
| 2 | ✅ | Products, categories, images |
| 3 | ✅ | Orders, campaigns, reviews, JWT, themes |
| 4 | ✅ | PostgreSQL repos, migrations, image upload |
| 5 | ✅ | Admin CRUD, public review submission, order filtering, UUID fixes |
| 6 | 🔄 | Docker, docs, deployment prep |
