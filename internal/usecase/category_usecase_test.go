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
