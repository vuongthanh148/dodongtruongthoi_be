# Category ID Normalization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Give `categories` a real auto-increment integer primary key, and
collapse the two hand-typed, drift-prone text identifiers (`id`, `slug`)
into one auto-generated-but-editable slug — with almost no ripple into
frontend code, because the JSON wire contract keeps the same `"id"` key
and value shape it has today.

**Architecture:** `categories.id` becomes a `SERIAL` (internal DB PK,
used only for FK joins — never serialized to JSON). The existing
`categories.slug` column becomes the sole external identifier: it's what
`domain.Category.ID` maps to in Go, what `"id"` means in the JSON API,
and what the admin form edits. `products.category_id` is retyped from
`TEXT` to `INT`, resolved by joining on `categories.slug` at read/write
time so `domain.Product.CategoryID` also keeps meaning "the category's
slug string" with zero behavior change for any existing caller.

**Tech Stack:** Go 1.23, chi router, pgx/v5, PostgreSQL, Next.js 16 /
React (admin UI only touched here).

**Spec:** `docs/superpowers/specs/2026-09-27-taxonomy-normalization-design.md`
— this plan implements the categories portion of that spec. Products
(the larger piece — `product_id` fans out to 7 other tables) and the
zodiac/purpose/bg_tone/frame_style lookup tables are separate follow-up
plans; see "Out of scope" below.

## Global Constraints

- No new Go dependencies beyond promoting `golang.org/x/text` from
  `// indirect` to direct in `go.mod` (already vendored transitively;
  used here for Unicode diacritic stripping).
- Every SQL migration ships in one file, one transaction — this
  codebase's existing migrations (`internal/infrastructure/database`,
  `migrations/embed.go`) apply files in filename order automatically on
  server startup; there is no separate migration-runner command.
- No test framework is present in this repo (`go.mod` has no
  `testify`/mocking lib, zero existing `_test.go` files). New tests use
  stdlib `testing` only, with hand-written fakes for interfaces —
  matching the "no ORM, no framework, explicit code" style already used
  throughout `internal/`.
- The repo layer and HTTP layer have no automated test harness (no test
  DB fixture, no httptest setup anywhere in the codebase today). Steps
  that touch those layers are verified manually via `psql`/`curl`
  against the local dev Postgres (already running per this session —
  `postgres://postgres:postgres@localhost:5432/dodongtruongthoi`) and
  the local dev server (`go run ./cmd/server`, port 8080), the same way
  the rest of this codebase has been hand-verified so far. This is a
  deliberate scope choice, not an oversight — introducing DB-test
  infra is a separate concern from this feature.

## Out of scope (separate plans)

- `products.id` int PK migration (touches `product_images`,
  `product_sizes`, `product_skus`, `reviews`, `wishlists`,
  `campaign_products`, `order_items` — every FK fans out from here).
- `zodiac_signs`, `purpose_places`, `purpose_use_tags`,
  `purpose_avoid_tags`, `bg_tones`, `frame_styles` lookup tables and
  their admin CRUD pages (depend on `products.id` being int first, per
  the design spec's join-table shape).

---

### Task 1: Vietnamese-aware slugify utility

**Files:**
- Create: `pkg/slug/slug.go`
- Test: `pkg/slug/slug_test.go`

**Interfaces:**
- Produces: `slug.Generate(name string) string` — later tasks (3, 5)
  call this to turn a category name into a URL-safe identifier.

- [ ] **Step 1: Write the failing test**

```go
// pkg/slug/slug_test.go
package slug

import "testing"

func TestGenerate(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"vietnamese diacritics", "Tranh Đồng", "tranh-dong"},
		{"multi-word vietnamese", "Tượng Đồng Trang Trí", "tuong-dong-trang-tri"},
		{"punctuation and extra spaces", "  Hello,   World!!  ", "hello-world"},
		{"ampersand", "Café & Bar", "cafe-bar"},
		{"already a slug", "tranh-dong", "tranh-dong"},
		{"empty", "", ""},
		{"mixed case english", "Best Seller 2026", "best-seller-2026"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Generate(tc.in)
			if got != tc.want {
				t.Errorf("Generate(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd dodongtruongthoi_be && go test ./pkg/slug/...`
Expected: FAIL — `package slug: build constraints exclude ...` or
`undefined: Generate` (the package/function doesn't exist yet).

- [ ] **Step 3: Write minimal implementation**

```go
// pkg/slug/slug.go
package slug

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	nonAlnum  = regexp.MustCompile(`[^a-z0-9]+`)
	trimDash  = regexp.MustCompile(`^-+|-+$`)
	stripMarks = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

// Generate converts a display name (Vietnamese or English) into a
// lowercase, URL-safe, latin slug: diacritics stripped, everything
// that isn't a letter or digit collapsed to a single hyphen, leading
// and trailing hyphens trimmed. "Tranh Đồng" -> "tranh-dong".
func Generate(name string) string {
	// đ/Đ don't decompose under NFD (they're distinct letters, not a
	// base letter + combining mark), so fold them explicitly first.
	folded := strings.NewReplacer("đ", "d", "Đ", "D").Replace(name)

	ascii, _, err := transform.String(stripMarks, folded)
	if err != nil {
		ascii = folded
	}

	lower := strings.ToLower(ascii)
	dashed := nonAlnum.ReplaceAllString(lower, "-")
	return trimDash.ReplaceAllString(dashed, "")
}
```

- [ ] **Step 4: Add the new import and tidy go.mod**

Run:
```bash
cd dodongtruongthoi_be
go mod tidy
```
Expected: `golang.org/x/text` moves from the `// indirect` require block
to the direct require block in `go.mod`.

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./pkg/slug/... -v`
Expected: PASS, all 7 subtests green.

- [ ] **Step 6: Commit**

```bash
git add pkg/slug/slug.go pkg/slug/slug_test.go go.mod go.sum
git commit -m "feat: add Vietnamese-aware slug generator"
```

---

### Task 2: Migration — categories int PK, products.category_id retype

**Files:**
- Create: `migrations/010_categories_int_pk.sql`

**Interfaces:**
- Consumes: nothing (pure SQL, runs standalone).
- Produces: `categories.id INT PRIMARY KEY` (new internal PK),
  `categories.legacy_id TEXT` (old hand-typed id, kept for audit only,
  unused by application code after this plan), `categories.slug`
  unchanged (still `TEXT NOT NULL UNIQUE`, now the sole external key),
  `products.category_id INT REFERENCES categories(id)` (was `TEXT`).

- [ ] **Step 1: Write the migration file**

```sql
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
```

- [ ] **Step 2: Apply the migration and verify**

Run:
```bash
cd dodongtruongthoi_be
go run ./cmd/server &
sleep 2
kill %1
```
Expected in server output: `✓ Migrations up to date` (the embedded
`010_categories_int_pk.sql` was picked up and applied — `internal/infrastructure/database/migrate.go`
applies any `.sql` file in `migrations/embed.go`'s embed FS that
hasn't run yet, in filename order).

- [ ] **Step 3: Verify the data manually**

Run:
```bash
PGPASSWORD=postgres psql -h localhost -U postgres -d dodongtruongthoi -c \
  "SELECT id, legacy_id, slug FROM categories ORDER BY id;"
PGPASSWORD=postgres psql -h localhost -U postgres -d dodongtruongthoi -c \
  "SELECT p.title, p.category_id, c.slug FROM products p JOIN categories c ON c.id = p.category_id;"
```
Expected: `categories.id` is `1..8` (or however many rows exist), each
row's `legacy_id`/`slug` match what was there before the migration
(`tranh-dong`, `tranh-phong-thuy`, etc. — verified via `psql` earlier in
this session), and every product's `category_id` resolves to the
correct category slug (`tranh-dong` for all 4 current products).

- [ ] **Step 4: Commit**

```bash
git add migrations/010_categories_int_pk.sql
git commit -m "migrate: give categories a real int PK, retype products.category_id"
```

---

### Task 3: Domain struct cleanup — drop redundant Category.Slug field

**Files:**
- Modify: `internal/domain/entities.go:7-19`

**Interfaces:**
- Produces: `domain.Category` with fields `ID, Name, Description, Tone,
  ImageURL, SortOrder, IsActive, ProductCount, CreatedAt, UpdatedAt` —
  no `Slug` field (it was a second hand-typed identifier duplicating
  `ID`'s job; Task 2's migration + Task 4's repo rewrite make `ID` be
  the slug value, so the separate field is now dead weight). `ID`'s doc
  comment records that it's the URL-safe slug, not the raw DB int PK.

- [ ] **Step 1: Edit the struct**

In `internal/domain/entities.go`, replace lines 7-19:

```go
type Category struct {
	ID           string    `json:"id"` // URL-safe slug; auto-generated from Name if not provided, editable until first save
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"`
	Tone         string    `json:"tone"`
	ImageURL     *string   `json:"image_url,omitempty"`
	SortOrder    int       `json:"sort_order"`
	IsActive     bool      `json:"is_active"`
	ProductCount int       `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
```

- [ ] **Step 2: Confirm it doesn't compile yet (expected — dependents still reference .Slug)**

Run: `cd dodongtruongthoi_be && go build ./... 2>&1 | head -20`
Expected: FAIL, errors in `internal/repository/postgres/category_repo.go`,
`internal/delivery/http/handler/admin_category_handler.go`,
`internal/usecase/category_usecase.go` — all `unknown field Slug`.
These get fixed in Tasks 4-6.

- [ ] **Step 3: Commit** (after Tasks 4-6 make the build pass again — see Task 6's commit step, which commits this change together with theirs since Go won't compile with only Task 3 applied)

Skip committing here; carried into Task 6's commit.

---

### Task 4: Category repository — key everything on `slug`

**Files:**
- Modify: `internal/repository/postgres/category_repo.go` (whole file)

**Interfaces:**
- Consumes: `domain.Category` (from Task 3, no `.Slug` field).
- Produces: same method signatures as before —
  `List(ctx, includeInactive) ([]domain.Category, error)`,
  `Get(ctx, id string, includeInactive) (domain.Category, bool, error)`,
  `Create(ctx, domain.Category) (domain.Category, error)`,
  `Update(ctx, id string, domain.Category) (domain.Category, error)`,
  `Delete(ctx, id string) error` — but every `id string` parameter now
  means "slug", and every SQL query filters/writes the `slug` column
  instead of the retired `legacy_id`/old `id` column.

- [ ] **Step 1: Rewrite the file**

```go
// internal/repository/postgres/category_repo.go
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) List(ctx context.Context, includeInactive bool) ([]domain.Category, error) {
	query := `SELECT c.slug, c.name, c.description, c.tone, c.image_url, c.sort_order, c.is_active,
		COUNT(p.id) FILTER (WHERE p.is_active = true) AS product_count,
		c.created_at, c.updated_at
		FROM categories c
		LEFT JOIN products p ON p.category_id = c.id`

	if !includeInactive {
		query += " WHERE c.is_active = true"
	}
	query += " GROUP BY c.id ORDER BY c.sort_order, c.created_at DESC"

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.Tone, &c.ImageURL,
			&c.SortOrder, &c.IsActive, &c.ProductCount, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *CategoryRepository) Get(ctx context.Context, slug string, includeInactive bool) (domain.Category, bool, error) {
	query := "SELECT slug, name, description, tone, image_url, sort_order, is_active, created_at, updated_at FROM categories WHERE slug = $1"

	if !includeInactive {
		query += " AND is_active = true"
	}

	var c domain.Category
	err := r.pool.QueryRow(ctx, query, slug).Scan(
		&c.ID, &c.Name, &c.Description, &c.Tone, &c.ImageURL,
		&c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Category{}, false, nil
		}
		return domain.Category{}, false, err
	}
	return c, true, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c domain.Category) (domain.Category, error) {
	query := `INSERT INTO categories (slug, name, description, tone, image_url, sort_order, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		c.ID, c.Name, c.Description, c.Tone, c.ImageURL,
		c.SortOrder, c.IsActive, c.CreatedAt, c.UpdatedAt,
	).Scan(&c.CreatedAt, &c.UpdatedAt)

	return c, err
}

func (r *CategoryRepository) Update(ctx context.Context, slug string, c domain.Category) (domain.Category, error) {
	if slug == "" {
		return domain.Category{}, errors.New("category slug is required")
	}

	query := `UPDATE categories SET
		name = COALESCE(NULLIF($2, ''), name),
		description = $3,
		tone = COALESCE(NULLIF($4, ''), tone),
		image_url = COALESCE(NULLIF($5, ''), image_url),
		sort_order = $6,
		is_active = $7,
		updated_at = now()
	WHERE slug = $1
	RETURNING slug, name, description, tone, image_url, sort_order, is_active, created_at, updated_at`

	var result domain.Category
	err := r.pool.QueryRow(ctx, query,
		slug, c.Name, c.Description, c.Tone, c.ImageURL,
		c.SortOrder, c.IsActive,
	).Scan(
		&result.ID, &result.Name, &result.Description, &result.Tone, &result.ImageURL,
		&result.SortOrder, &result.IsActive, &result.CreatedAt, &result.UpdatedAt,
	)

	return result, err
}

func (r *CategoryRepository) Delete(ctx context.Context, slug string) error {
	result, err := r.pool.Exec(ctx, "UPDATE categories SET is_active = false, updated_at = now() WHERE slug = $1", slug)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("category not found")
	}
	return nil
}
```

Note the `Update` query no longer accepts a slug change (dropped the
old `slug = COALESCE(NULLIF($3, ''), slug)` line) — renaming a
category's slug after creation would break existing product/category
URLs and is out of scope; the admin UI (Task 9) keeps the identifier
field disabled after creation, same as it disables `id` today.

- [ ] **Step 2: Confirm compile progress**

Run: `cd dodongtruongthoi_be && go build ./... 2>&1 | head -20`
Expected: `category_repo.go` errors gone; remaining errors only in
`admin_category_handler.go` and `category_usecase.go` (fixed next).

---

### Task 5: Category usecase — auto-generate + de-duplicate slug on create

**Files:**
- Modify: `internal/usecase/category_usecase.go` (whole file)
- Test: `internal/usecase/category_usecase_test.go`

**Interfaces:**
- Consumes: `slug.Generate(string) string` (Task 1),
  `domain.CategoryRepository` interface (unchanged shape, Task 4's impl).
- Produces: `CreateCategory(ctx, domain.Category) (domain.Category, error)`
  now accepts a `Category` with a blank `ID` (auto-generates from
  `Name`) or a caller-supplied `ID` (still normalized through
  `slug.Generate` and de-duplicated) — this replaces the old
  "id and name are required" hard requirement.

- [ ] **Step 1: Write the failing test**

```go
// internal/usecase/category_usecase_test.go
package usecase

import (
	"context"
	"testing"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// fakeCategoryRepo is an in-memory domain.CategoryRepository for
// usecase-level unit tests — no database needed.
type fakeCategoryRepo struct {
	byID map[string]domain.Category
}

func newFakeCategoryRepo() *fakeCategoryRepo {
	return &fakeCategoryRepo{byID: map[string]domain.Category{}}
}

func (f *fakeCategoryRepo) List(ctx context.Context, includeInactive bool) ([]domain.Category, error) {
	var out []domain.Category
	for _, c := range f.byID {
		out = append(out, c)
	}
	return out, nil
}

func (f *fakeCategoryRepo) Get(ctx context.Context, id string, includeInactive bool) (domain.Category, bool, error) {
	c, ok := f.byID[id]
	return c, ok, nil
}

func (f *fakeCategoryRepo) Create(ctx context.Context, c domain.Category) (domain.Category, error) {
	f.byID[c.ID] = c
	return c, nil
}

func (f *fakeCategoryRepo) Update(ctx context.Context, id string, c domain.Category) (domain.Category, error) {
	c.ID = id
	f.byID[id] = c
	return c, nil
}

func (f *fakeCategoryRepo) Delete(ctx context.Context, id string) error {
	delete(f.byID, id)
	return nil
}

func TestCreateCategory_AutoGeneratesSlugFromName(t *testing.T) {
	repo := newFakeCategoryRepo()
	u := NewCategoryUsecase(repo)

	got, err := u.CreateCategory(context.Background(), domain.Category{Name: "Tranh Đồng"})
	if err != nil {
		t.Fatalf("CreateCategory returned error: %v", err)
	}
	if got.ID != "tranh-dong" {
		t.Errorf("ID = %q, want %q", got.ID, "tranh-dong")
	}
}

func TestCreateCategory_DeduplicatesCollidingSlug(t *testing.T) {
	repo := newFakeCategoryRepo()
	u := NewCategoryUsecase(repo)

	first, err := u.CreateCategory(context.Background(), domain.Category{Name: "Tranh Phong Cảnh"})
	if err != nil {
		t.Fatalf("first CreateCategory returned error: %v", err)
	}
	if first.ID != "tranh-phong-canh" {
		t.Fatalf("first.ID = %q, want %q", first.ID, "tranh-phong-canh")
	}

	second, err := u.CreateCategory(context.Background(), domain.Category{Name: "Tranh Phong Cảnh"})
	if err != nil {
		t.Fatalf("second CreateCategory returned error: %v", err)
	}
	if second.ID != "tranh-phong-canh-2" {
		t.Errorf("second.ID = %q, want %q", second.ID, "tranh-phong-canh-2")
	}
}

func TestCreateCategory_RequiresName(t *testing.T) {
	repo := newFakeCategoryRepo()
	u := NewCategoryUsecase(repo)

	_, err := u.CreateCategory(context.Background(), domain.Category{})
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd dodongtruongthoi_be && go test ./internal/usecase/... -run TestCreateCategory -v`
Expected: FAIL — `CreateCategory` still requires `c.ID != ""` (old
behavior), so `TestCreateCategory_AutoGeneratesSlugFromName` fails with
an error instead of a generated slug.

- [ ] **Step 3: Rewrite the usecase**

```go
// internal/usecase/category_usecase.go
package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
	"github.com/vuongthanh148/dodongtruongthoi_be/pkg/slug"
)

type CategoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx, false)
}

func (u *CategoryUsecase) GetCategory(ctx context.Context, id string) (domain.Category, bool, error) {
	return u.categoryRepo.Get(ctx, id, false)
}

func (u *CategoryUsecase) ListAllCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx, true)
}

func (u *CategoryUsecase) CreateCategory(ctx context.Context, c domain.Category) (domain.Category, error) {
	if c.Name == "" {
		return domain.Category{}, errors.New("name is required")
	}

	base := c.ID
	if base == "" {
		base = c.Name
	}
	c.ID = u.uniqueSlug(ctx, slug.Generate(base))

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.IsActive = true

	return u.categoryRepo.Create(ctx, c)
}

// uniqueSlug returns base, or base-2, base-3, ... — whichever is the
// first not already taken by an existing category (active or not).
func (u *CategoryUsecase) uniqueSlug(ctx context.Context, base string) string {
	candidate := base
	for i := 2; ; i++ {
		_, exists, err := u.categoryRepo.Get(ctx, candidate, true)
		if err != nil || !exists {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

func (u *CategoryUsecase) UpdateCategory(ctx context.Context, id string, updates domain.Category) (domain.Category, error) {
	if id == "" {
		return domain.Category{}, errors.New("category id is required")
	}

	existing, ok, err := u.categoryRepo.Get(ctx, id, true)
	if err != nil {
		return domain.Category{}, err
	}
	if !ok {
		return domain.Category{}, errors.New("category not found")
	}

	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.Description != nil {
		existing.Description = updates.Description
	}
	if updates.Tone != "" {
		existing.Tone = updates.Tone
	}
	if updates.ImageURL != nil && *updates.ImageURL != "" {
		existing.ImageURL = updates.ImageURL
	}
	existing.SortOrder = updates.SortOrder
	existing.IsActive = updates.IsActive
	existing.UpdatedAt = time.Now()

	return u.categoryRepo.Update(ctx, id, existing)
}

func (u *CategoryUsecase) DeleteCategory(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category id is required")
	}
	return u.categoryRepo.Delete(ctx, id)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/usecase/... -run TestCreateCategory -v`
Expected: PASS, all 3 tests green.

- [ ] **Step 5: Commit**

```bash
git add internal/usecase/category_usecase.go internal/usecase/category_usecase_test.go
git commit -m "feat: auto-generate and de-duplicate category slugs on create"
```

---

### Task 6: Admin handler — drop manual ID requirement, drop Slug field

**Files:**
- Modify: `internal/delivery/http/handler/admin_category_handler.go:22-51`

**Interfaces:**
- Consumes: `PlatformUsecase.CreateCategory` (Task 5 — now tolerates
  blank `ID`).
- Produces: `POST /api/v1/admin/categories` request body no longer has
  a `slug` field; `id` is now optional (an override, not a requirement).

- [ ] **Step 1: Edit `CreateCategory`**

Replace lines 22-51 of `admin_category_handler.go`:

```go
func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Tone        string `json:"tone"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	category := domain.Category{
		ID:          body.ID,
		Name:        body.Name,
		Description: ptrIfNotEmpty(body.Description),
		Tone:        body.Tone,
		ImageURL:    ptrIfNotEmpty(body.ImageURL),
		SortOrder:   body.SortOrder,
	}

	result, err := h.platform.CreateCategory(r.Context(), category)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(w, http.StatusCreated, result)
}
```

- [ ] **Step 2: Edit `UpdateCategory`**

In the same file, `UpdateCategory`'s body struct currently has a `Slug`
field (line ~68) — remove it:

```go
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Tone        string `json:"tone"`
		ImageURL    string `json:"image_url"`
		SortOrder   int    `json:"sort_order"`
		IsActive    bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	updates := domain.Category{
		Name:        body.Name,
		Description: ptrIfNotEmpty(body.Description),
		Tone:        body.Tone,
		ImageURL:    ptrIfNotEmpty(body.ImageURL),
		SortOrder:   body.SortOrder,
		IsActive:    body.IsActive,
	}
```

- [ ] **Step 3: Verify the whole project compiles**

Run: `cd dodongtruongthoi_be && go build ./...`
Expected: exits 0, no errors — this is the point where Tasks 3-6's
changes all land together.

- [ ] **Step 4: Run full test suite**

Run: `go test ./...`
Expected: PASS (`pkg/slug`, `internal/usecase` tests; all other
packages report `no test files`, which is fine).

- [ ] **Step 5: Commit**

```bash
git add internal/domain/entities.go internal/repository/postgres/category_repo.go internal/delivery/http/handler/admin_category_handler.go
git commit -m "refactor: categories use slug as sole external id, drop redundant field"
```

---

### Task 7: Product repository — resolve category_id through the slug join

**Files:**
- Modify: `internal/repository/postgres/product_repo.go:23-26` (the
  `productColumns` constant and every query built from it)

**Interfaces:**
- Consumes: `categories(id, slug)` (Task 2's schema).
- Produces: `domain.Product.CategoryID` still holds the category's
  *slug* string (unchanged meaning for every caller in
  `product_usecase.go`, `admin_product_handler.go`, and the FE) — only
  the SQL underneath changes to join through the new int FK.

- [ ] **Step 1: Change the SELECT column list to join for the slug**

Replace lines 23-26:

```go
const productColumns = `
	p.id, p.title, p.subtitle, c.slug AS category_id, p.badge, p.base_price, p.description, p.meaning,
	p.variant_options, p.default_variant, p.zodiac_ids, p.purpose_place, p.purpose_use,
	p.purpose_avoid, p.specs, p.requires_size, p.is_active, p.sort_order, p.created_at, p.updated_at`

const productFrom = `FROM products p JOIN categories c ON c.id = p.category_id`
```

- [ ] **Step 2: Update `List` to use the join**

In `List`, change:
```go
query := "SELECT " + productColumns + " FROM products WHERE 1=1"
```
to:
```go
query := "SELECT " + productColumns + " " + productFrom + " WHERE 1=1"
```
and change the category filter (currently `category_id = $N` matching
the old TEXT column) to filter by slug through the join:
```go
if q.Category != "" {
	query += fmt.Sprintf(" AND c.slug = $%d", argCount)
	args = append(args, q.Category)
	argCount++
}
```
Every other unqualified column reference in `List` (`is_active`,
`base_price`, `sort_order`, `created_at`) needs a `p.` prefix now that
the query joins two tables — update the `!includeInactive`, `orderBy`,
and `LIMIT`/`OFFSET` blocks accordingly: `p.is_active = true`,
`p.base_price ASC`, etc.

- [ ] **Step 3: Update `Get` to use the join**

```go
func (r *ProductRepository) Get(ctx context.Context, id string, includeInactive bool) (domain.Product, bool, error) {
	query := "SELECT " + productColumns + " " + productFrom + " WHERE p.id = $1"
	if !includeInactive {
		query += " AND p.is_active = true"
	}

	p, err := scanProduct(r.pool.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Product{}, false, nil
	}
	if err != nil {
		return domain.Product{}, false, err
	}
	return p, true, nil
}
```

- [ ] **Step 4: Update `Create` to resolve the slug to an int FK on insert**

In `Create`, change the INSERT's `category_id` placeholder to a
subquery against `categories.slug`:

```go
query := `INSERT INTO products (
	id, title, subtitle, category_id, badge, base_price, description, meaning,
	variant_options, default_variant, zodiac_ids, purpose_place, purpose_use,
	purpose_avoid, specs, requires_size, is_active, sort_order, created_at, updated_at
) VALUES (
	$1,$2,$3,(SELECT id FROM categories WHERE slug = $4),$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20
) RETURNING created_at, updated_at`

err := r.pool.QueryRow(ctx, query,
	p.ID, p.Title, p.Subtitle, p.CategoryID, p.Badge, p.BasePrice,
	p.Description, p.Meaning,
	variantOptionsJSON, defaultVariantJSON,
	p.ZodiacIDs, p.PurposePlace, p.PurposeUse, p.PurposeAvoid,
	p.Specs, p.RequiresSize, p.IsActive, p.SortOrder, p.CreatedAt, p.UpdatedAt,
).Scan(&p.CreatedAt, &p.UpdatedAt)
```

(Note: the Go `args` slice passed to `QueryRow` is unchanged — still
`p.ID, p.Title, p.Subtitle, p.CategoryID, p.Badge, ...` in that order.
Only the SQL text changes: the `$4` placeholder, which already meant
`p.CategoryID`, is now wrapped in `(SELECT id FROM categories WHERE
slug = $4)` instead of being inserted raw. No argument positions move.)

- [ ] **Step 5: Update `Update` similarly**

In `Update`, change:
```go
category_id    = COALESCE(NULLIF($4, ''), category_id),
```
to:
```go
category_id    = COALESCE((SELECT id FROM categories WHERE slug = NULLIF($4, '')), category_id),
```
and change the final `RETURNING` clause from `RETURNING ` +
`productColumns` (which now requires the join, unavailable in a plain
`UPDATE ... RETURNING`) to a two-step approach: run the `UPDATE`
without `RETURNING`, then call `r.Get(ctx, id, true)` for the result:

```go
func (r *ProductRepository) Update(ctx context.Context, id string, p domain.Product) (domain.Product, error) {
	if id == "" {
		return domain.Product{}, errors.New("product id is required")
	}

	variantOptionsJSON, _ := json.Marshal(p.VariantOptions)
	defaultVariantJSON, _ := json.Marshal(p.DefaultVariant)

	query := `UPDATE products SET
		title          = COALESCE(NULLIF($2, ''), title),
		subtitle       = $3,
		category_id    = COALESCE((SELECT id FROM categories WHERE slug = NULLIF($4, '')), category_id),
		badge          = $5,
		base_price     = CASE WHEN $6 > 0 THEN $6 ELSE base_price END,
		description    = $7,
		meaning        = $8,
		variant_options = $9,
		default_variant = $10,
		zodiac_ids     = CASE WHEN array_length($11::text[], 1) > 0 THEN $11 ELSE zodiac_ids END,
		purpose_place  = $12,
		purpose_use    = $13,
		purpose_avoid  = $14,
		specs          = COALESCE(NULLIF($15::jsonb, '{}'), specs),
		requires_size  = $16,
		is_active      = $17,
		sort_order     = $18,
		updated_at     = now()
	WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		id, p.Title, p.Subtitle, p.CategoryID, p.Badge, p.BasePrice,
		p.Description, p.Meaning,
		variantOptionsJSON, defaultVariantJSON,
		p.ZodiacIDs, p.PurposePlace, p.PurposeUse, p.PurposeAvoid,
		p.Specs, p.RequiresSize, p.IsActive, p.SortOrder,
	)
	if err != nil {
		return domain.Product{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.Product{}, errors.New("product not found")
	}

	result, _, err := r.Get(ctx, id, true)
	return result, err
}
```

- [ ] **Step 6: Build and manually verify**

Run:
```bash
go build ./...
go run ./cmd/server &
sleep 2
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/admin/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['token'])")
curl -s http://localhost:8080/api/v1/products | python3 -m json.tool | head -20
curl -s "http://localhost:8080/api/v1/products?category=tranh-dong" | python3 -c "import sys,json; print(len(json.load(sys.stdin)['data']))"
kill %1
```
Expected: product list returns products with `"category_id": "tranh-dong"`
(string slug, unchanged shape from before this plan), and the
`?category=tranh-dong` filter still returns all 4 existing products.

- [ ] **Step 7: Commit**

```bash
git add internal/repository/postgres/product_repo.go
git commit -m "fix: resolve products.category_id through the new int FK via slug join"
```

---

### Task 8: Frontend types — drop redundant AdminCategory.slug

**Files:**
- Modify: `src/app/../dodongtruongthoi_fe/src/lib/types.ts:151-160`

**Interfaces:**
- Produces: `AdminCategory` without a `slug` field (the `id` field is
  now the single editable identifier, matching the BE's Task 3 change).

- [ ] **Step 1: Edit the interface**

In `dodongtruongthoi_fe/src/lib/types.ts`, replace lines 151-160:

```ts
export interface AdminCategory {
  id: string
  name: string
  description: string | null
  tone: string
  image_url: string | null
  sort_order: number
  is_active: boolean
}
```

- [ ] **Step 2: Verify no other file reads `.slug` on a category object**

Run: `cd dodongtruongthoi_fe && grep -rn "\.slug\b" src`
Expected: zero remaining matches (Task 9 removes the only current
usage, in `admin/categories/page.tsx`).

- [ ] **Step 3: Commit**

```bash
git add src/lib/types.ts
git commit -m "refactor: drop redundant AdminCategory.slug field"
```

(commit path is relative to `dodongtruongthoi_fe/`)

---

### Task 9: Admin categories page — one auto-suggested Slug field

**Files:**
- Modify: `src/app/admin/categories/page.tsx` (form state, `saveCategory`,
  `startEdit`, and the JSX fields at lines 10-18, 44-52, 60-67, 100-108)

**Interfaces:**
- Consumes: `slug.Generate`-equivalent client-side preview (a small
  local helper, not imported from the backend) — cosmetic only; the
  server (Task 5) is the source of truth and re-normalizes/de-duplicates
  regardless of what the client sends.
- Produces: form has one identifier field, labeled "Slug", auto-filled
  from the Name field while creating (until the admin edits it by
  hand), disabled after creation — same UX shape as the old "ID (create
  only)" field, minus the second redundant "Slug" input.

- [ ] **Step 1: Add a local slug-preview helper**

At the top of `src/app/admin/categories/page.tsx`, after the existing
imports, add:

```ts
// Client-side preview only — the server (CategoryUsecase.CreateCategory)
// re-normalizes and de-duplicates the slug regardless of what's sent.
function previewSlug(name: string): string {
  return name
    .normalize('NFD')
    .replace(/[̀-ͯ]/g, '')
    .replace(/đ/g, 'd')
    .replace(/Đ/g, 'D')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}
```

- [ ] **Step 2: Simplify form state**

Replace `emptyForm` (lines 10-18):

```ts
const emptyForm = {
  id: '',
  name: '',
  description: '',
  tone: 'gold',
  image_url: '',
  sort_order: '0',
  is_active: true,
}
```

Add a tracker for whether the admin has hand-edited the slug (so
auto-suggest stops overwriting their edit):

```ts
const [slugTouched, setSlugTouched] = useState(false)
```//_ place this alongside the other `useState` calls near the top of `AdminCategoriesPage`.

- [ ] **Step 3: Update `startEdit` and `startCreate`**

`startEdit` (around line 44): drop the `slug: category.slug` line:

```ts
function startEdit(category: AdminCategory) {
	setEditingId(category.id)
	setSlugTouched(true) // editing an existing category never auto-regenerates its slug
	setForm({
		id: category.id,
		name: category.name,
		description: category.description ?? '',
		tone: category.tone,
		image_url: category.image_url ?? '',
		sort_order: String(category.sort_order),
		is_active: category.is_active,
	})
}
```

`startCreate` should also reset `slugTouched`:

```ts
function startCreate() {
	setEditingId(null)
	setSlugTouched(false)
	setForm(emptyForm)
}
```

- [ ] **Step 4: Update `saveCategory`**

Remove `slug: form.slug.trim()` from the payload (around line 60):

```ts
async function saveCategory() {
	const payload = {
		name: form.name.trim(),
		description: form.description.trim(),
		tone: form.tone,
		image_url: form.image_url.trim(),
		sort_order: Number(form.sort_order) || 0,
		is_active: form.is_active,
	}

	const saved = editingId
		? await adminPut<AdminCategory>(`/categories/${editingId}`, payload)
		: await adminPost<AdminCategory>('/categories', { id: form.id.trim(), ...payload })

	if (saved) {
		toast.success(editingId ? 'Category updated' : 'Category created')
		await loadCategories()
		startCreate()
	} else {
		toast.error('Failed to save category')
	}
}
```

- [ ] **Step 5: Update the JSX fields**

Replace the "ID (create only)" + "Slug" field pair (lines 100-108) with
one field, and wire the Name field to drive the auto-suggestion:

```tsx
<Field label="Name">
  <input
    value={form.name}
    onChange={(event) => {
      const name = event.target.value
      setForm((prev) => ({
        ...prev,
        name,
        id: slugTouched || editingId ? prev.id : previewSlug(name),
      }))
    }}
    style={inputStyle}
  />
</Field>
<Field label="Slug">
  <input
    value={form.id}
    onChange={(event) => {
      setSlugTouched(true)
      setForm((prev) => ({ ...prev, id: event.target.value }))
    }}
    style={inputStyle}
    disabled={Boolean(editingId)}
  />
</Field>
```

Remove the old separate "ID (create only)" field block entirely (it's
replaced by the pair above — Name now comes first in the grid, Slug
second, matching the order these two fields already occupied).

- [ ] **Step 6: Manual verification**

Run:
```bash
cd dodongtruongthoi_fe
npm run dev &
```
Then in a browser (or via `curl` against the already-running BE from
Task 7's verification): open `/admin/login`, sign in (`admin` /
`admin123`), go to `/admin/categories`, type a Vietnamese name like
"Tranh Sơn Mài" into the Name field, and confirm the Slug field
auto-fills with `tranh-son-mai` before you save. Create it, confirm it
appears in the list with that slug as its `id`. Edit an existing
category and confirm the Slug field is disabled and unchanged.

- [ ] **Step 7: Commit**

```bash
git add src/app/admin/categories/page.tsx
git commit -m "feat: auto-suggest category slug from name, drop redundant id/slug pair"
```

---

### Task 10: End-to-end verification

**Files:** none (verification only)

- [ ] **Step 1: Full stack smoke test**

With both `go run ./cmd/server` (port 8080) and `npm run dev` (port
3000) running:

1. Storefront: open `http://localhost:3000/categories/tranh-dong` —
   confirm it still loads products (proves `GetCategory` and the
   `?category=` product filter both still resolve by the same slug
   string as before this plan).
2. Storefront: open `http://localhost:3000/products/vinh-quy-bai-to` —
   confirm the category link at the bottom of the page still points to
   `/categories/tranh-dong` and works.
3. Admin: `/admin/categories` — create a category, edit it, delete it
   (soft-delete via the existing Delete button), confirm each action
   succeeds with a toast and the list refreshes correctly.
4. Admin: `/admin/products/new` — create a product, pick a category
   from the dropdown, save, confirm it appears correctly in
   `/admin/products` with the right category name shown (proves Task
   7's `product_repo.go` join works end-to-end through the full admin
   create path, not just the `curl` check from Task 7).

- [ ] **Step 2: Confirm no regressions in existing admin sections**

Quickly click through `/admin/orders`, `/admin/reviews`,
`/admin/campaigns` — none of these were touched by this plan, but they
all read `product.category_id` transitively through product listings,
so a quick look confirms nothing broke.

- [ ] **Step 3: Report status**

No commit for this task — it's verification of Tasks 1-9's combined
commits, already each committed individually.

---

## Self-Review Notes

- **Spec coverage:** This plan implements the "categories" portion of
  the design spec's Schema design section — the `ALTER TABLE categories`
  shape, plus the decision (refined during planning, not a spec change)
  to keep the JSON wire contract's `"id"` key mapped to the slug value
  rather than exposing the raw int, which achieves the spec's stated
  goal ("public URLs keep working via a slug column... Old TEXT ids ...
  URLs currently use products.id as the slug in routes, so slug becomes
  the public-facing key") with far less code churn than serializing both
  fields would have required. Products (int PK, 7 fan-out FKs) and the
  6 lookup tables are explicitly deferred — see "Out of scope" above.
- **Placeholder scan:** none found — every step has real code or a real
  command with expected output.
- **Type consistency:** `domain.Category.ID string` (Task 3) is written
  and read consistently as "the slug" across Tasks 4 (repo), 5
  (usecase), 6 (handler), matching `AdminCategory.id: string` on the FE
  (Task 8) and the `form.id` field in Task 9. `slug.Generate(string)
  string` (Task 1) is called with the same signature in Task 5.

---

**Plan complete and saved to `docs/superpowers/plans/2026-09-27-category-id-normalization.md`.**
