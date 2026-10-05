package usecase

import (
	"context"
	"testing"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// fakeProductRepo is an in-memory domain.ProductRepository for usecase-level tests.
type fakeProductRepo struct {
	domain.ProductRepository
	byID map[string]domain.Product
}

func (f *fakeProductRepo) Get(ctx context.Context, id string, includeInactive bool) (domain.Product, bool, error) {
	p, ok := f.byID[id]
	return p, ok, nil
}

func (f *fakeProductRepo) Update(ctx context.Context, id string, p domain.Product) (domain.Product, error) {
	f.byID[id] = p
	return p, nil
}

// fakeProductSKURepo is an in-memory domain.ProductSKURepository for usecase-level tests.
type fakeProductSKURepo struct {
	domain.ProductSKURepository
	byProduct map[string][]domain.ProductSKU
}

func (f *fakeProductSKURepo) ListByProduct(ctx context.Context, productID string) ([]domain.ProductSKU, error) {
	return f.byProduct[productID], nil
}

func (f *fakeProductSKURepo) SetByProduct(ctx context.Context, productID string, skus []domain.ProductSKU) error {
	f.byProduct[productID] = skus
	return nil
}

func newProductUsecaseForTest(basePrice int64) (*ProductUsecase, *fakeProductRepo, *fakeProductSKURepo) {
	productRepo := &fakeProductRepo{byID: map[string]domain.Product{
		"vinh-quy": {ID: "vinh-quy", Title: "Vinh Quy", BasePrice: basePrice, IsActive: true},
	}}
	skuRepo := &fakeProductSKURepo{byProduct: map[string][]domain.ProductSKU{}}
	u := NewProductUsecase(productRepo, nil, nil, skuRepo, nil, nil, nil, nil)
	return u, productRepo, skuRepo
}

func TestSetSKUs_BasePriceFollowsMinSKUPrice(t *testing.T) {
	u, productRepo, _ := newProductUsecaseForTest(1000)

	skus := []domain.ProductSKU{{Price: 9350000}, {Price: 8500000}}
	if err := u.SetSKUs(context.Background(), "vinh-quy", skus); err != nil {
		t.Fatalf("SetSKUs returned error: %v", err)
	}
	if got := productRepo.byID["vinh-quy"].BasePrice; got != 8500000 {
		t.Errorf("base_price = %d, want 8500000", got)
	}
}

func TestSetSKUs_EmptyListLeavesBasePriceUnchanged(t *testing.T) {
	u, productRepo, _ := newProductUsecaseForTest(1000)

	if err := u.SetSKUs(context.Background(), "vinh-quy", nil); err != nil {
		t.Fatalf("SetSKUs returned error: %v", err)
	}
	if got := productRepo.byID["vinh-quy"].BasePrice; got != 1000 {
		t.Errorf("base_price = %d, want 1000 (unchanged)", got)
	}
}

func TestSetSKUs_UnknownProductIsNotFound(t *testing.T) {
	u, _, skuRepo := newProductUsecaseForTest(1000)

	err := u.SetSKUs(context.Background(), "missing", []domain.ProductSKU{{Price: 5}})
	if err == nil {
		t.Fatal("SetSKUs on unknown product returned nil error")
	}
	if _, ok := skuRepo.byProduct["missing"]; ok {
		t.Error("SKUs were written for an unknown product")
	}
}

func TestUpdateProduct_IgnoresBasePriceWhenSKUsExist(t *testing.T) {
	u, productRepo, skuRepo := newProductUsecaseForTest(1000)
	skuRepo.byProduct["vinh-quy"] = []domain.ProductSKU{{Price: 9350000}, {Price: 8500000}}

	updates := domain.Product{Title: "Vinh Quy", BasePrice: 7000000}
	if _, err := u.UpdateProduct(context.Background(), "vinh-quy", updates, nil, nil); err != nil {
		t.Fatalf("UpdateProduct returned error: %v", err)
	}
	if got := productRepo.byID["vinh-quy"].BasePrice; got != 8500000 {
		t.Errorf("base_price = %d, want 8500000 (min SKU price)", got)
	}
}

func TestUpdateProduct_UsesBasePriceWhenNoSKUs(t *testing.T) {
	u, productRepo, _ := newProductUsecaseForTest(1000)

	updates := domain.Product{Title: "Vinh Quy", BasePrice: 7000000}
	if _, err := u.UpdateProduct(context.Background(), "vinh-quy", updates, nil, nil); err != nil {
		t.Fatalf("UpdateProduct returned error: %v", err)
	}
	if got := productRepo.byID["vinh-quy"].BasePrice; got != 7000000 {
		t.Errorf("base_price = %d, want 7000000", got)
	}
}

func TestCampaignDiscountApply(t *testing.T) {
	cases := []struct {
		name  string
		rule  campaignDiscount
		price int64
		want  int64
	}{
		{"percentage 30%", campaignDiscount{kind: "percentage", value: 30}, 8500000, 5950000},
		{"fixed amount", campaignDiscount{kind: "fixed_amount", value: 500000}, 8500000, 8000000},
		{"fixed amount larger than price clamps to zero", campaignDiscount{kind: "fixed_amount", value: 900}, 500, 0},
		{"unknown kind leaves price unchanged", campaignDiscount{kind: "bogus", value: 10}, 8500000, 8500000},
	}
	for _, tc := range cases {
		if got := tc.rule.apply(tc.price); got != tc.want {
			t.Errorf("%s: apply(%d) = %d, want %d", tc.name, tc.price, got, tc.want)
		}
	}
}
