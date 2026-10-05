package usecase

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// fakeOrderRepo records UpdateStatus calls for usecase-level unit tests.
type fakeOrderRepo struct {
	domain.OrderRepository
	orders      map[string]domain.Order
	created     []domain.Order
	updateCalls int
	lastID      string
	lastStatus  string
	lastNote    *string
}

func (f *fakeOrderRepo) Get(ctx context.Context, id string) (domain.Order, bool, error) {
	ord, ok := f.orders[id]
	return ord, ok, nil
}

func (f *fakeOrderRepo) TransitionStatus(ctx context.Context, id string, from []string, to string) (bool, error) {
	ord, ok := f.orders[id]
	if !ok || !slices.Contains(from, ord.Status) {
		return false, nil
	}
	ord.Status = to
	f.orders[id] = ord
	return true, nil
}

func (f *fakeOrderRepo) Create(ctx context.Context, ord domain.Order) (domain.Order, error) {
	f.created = append(f.created, ord)
	return ord, nil
}

func (f *fakeOrderRepo) UpdateStatus(ctx context.Context, id string, status string, adminNote *string) error {
	f.updateCalls++
	f.lastID = id
	f.lastStatus = status
	f.lastNote = adminNote
	if ord, ok := f.orders[id]; ok {
		ord.Status = status
		f.orders[id] = ord
	}
	return nil
}

const testOrderID = "11111111-1111-4111-8111-111111111111"

func newFakeOrderRepoWith(status string) *fakeOrderRepo {
	return &fakeOrderRepo{orders: map[string]domain.Order{
		testOrderID: {ID: testOrderID, Status: status},
	}}
}

func TestCancelOrder_CancellableStatusesBecomeCancelled(t *testing.T) {
	for _, status := range []string{"pending_confirm", "confirmed"} {
		repo := newFakeOrderRepoWith(status)
		u := NewOrderUsecase(repo, nil)

		if err := u.CancelOrder(context.Background(), testOrderID); err != nil {
			t.Fatalf("status %q: CancelOrder returned error: %v", status, err)
		}
		if got := repo.orders[testOrderID].Status; got != "cancelled" {
			t.Errorf("status %q: stored status = %q, want %q", status, got, "cancelled")
		}
	}
}

func TestCancelOrder_KeepsExistingAdminNote(t *testing.T) {
	repo := newFakeOrderRepoWith("pending_confirm")
	note := "khach hen giao sau"
	ord := repo.orders[testOrderID]
	ord.AdminNote = &note
	repo.orders[testOrderID] = ord
	u := NewOrderUsecase(repo, nil)

	if err := u.CancelOrder(context.Background(), testOrderID); err != nil {
		t.Fatalf("CancelOrder returned error: %v", err)
	}
	got := repo.orders[testOrderID].AdminNote
	if got == nil || *got != note {
		t.Errorf("stored admin note = %v, want %q", got, note)
	}
}

// racingOrderRepo reports the order as cancellable on the first read, but the stored
// status has already moved on, as if another request changed it in between.
type racingOrderRepo struct {
	*fakeOrderRepo
	reads int
}

func (r *racingOrderRepo) Get(ctx context.Context, id string) (domain.Order, bool, error) {
	r.reads++
	if r.reads == 1 {
		return domain.Order{ID: id, Status: "pending_confirm"}, true, nil
	}
	return r.fakeOrderRepo.Get(ctx, id)
}

func TestCancelOrder_ConcurrentStatusChangeIsNotOverwritten(t *testing.T) {
	inner := newFakeOrderRepoWith("shipped")
	repo := &racingOrderRepo{fakeOrderRepo: inner}
	u := NewOrderUsecase(repo, nil)

	err := u.CancelOrder(context.Background(), testOrderID)
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("err = %v, want ErrConflict", err)
	}
	if got := inner.orders[testOrderID].Status; got != "shipped" {
		t.Errorf("stored status = %q, want %q (must not be overwritten)", got, "shipped")
	}
}

func TestCancelOrder_NonCancellableStatusConflictsWithoutRepoUpdate(t *testing.T) {
	for _, status := range []string{"processing", "shipped", "completed", "cancelled"} {
		repo := newFakeOrderRepoWith(status)
		u := NewOrderUsecase(repo, nil)

		err := u.CancelOrder(context.Background(), testOrderID)
		if !errors.Is(err, domain.ErrConflict) {
			t.Errorf("status %q: err = %v, want ErrConflict", status, err)
		}
		if repo.updateCalls != 0 {
			t.Errorf("status %q: repo.UpdateStatus called %d times, want 0", status, repo.updateCalls)
		}
	}
}

func TestCancelOrder_MissingOrderIsNotFound(t *testing.T) {
	repo := newFakeOrderRepoWith("pending_confirm")
	u := NewOrderUsecase(repo, nil)

	err := u.CancelOrder(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	if repo.updateCalls != 0 {
		t.Errorf("repo.UpdateStatus called %d times, want 0", repo.updateCalls)
	}
}

func TestUpdateOrderStatus_ValidStatusPassedToRepo(t *testing.T) {
	repo := &fakeOrderRepo{}
	u := NewOrderUsecase(repo, nil)

	if err := u.UpdateOrderStatus(context.Background(), testOrderID, "shipped", "note"); err != nil {
		t.Fatalf("UpdateOrderStatus returned error: %v", err)
	}
	if repo.updateCalls != 1 {
		t.Fatalf("repo.UpdateStatus calls = %d, want 1", repo.updateCalls)
	}
	if repo.lastID != testOrderID || repo.lastStatus != "shipped" {
		t.Errorf("repo got id=%q status=%q, want id=%q status=%q", repo.lastID, repo.lastStatus, testOrderID, "shipped")
	}
	if repo.lastNote == nil || *repo.lastNote != "note" {
		t.Errorf("repo got admin note %v, want %q", repo.lastNote, "note")
	}
}

func TestUpdateOrderStatus_InvalidStatusRejectedWithoutRepoCall(t *testing.T) {
	for _, status := range []string{"garbage", "", "Completed", "pending"} {
		repo := &fakeOrderRepo{}
		u := NewOrderUsecase(repo, nil)

		err := u.UpdateOrderStatus(context.Background(), testOrderID, status, "")
		if !errors.Is(err, ErrInvalidOrderStatus) {
			t.Errorf("status %q: err = %v, want ErrInvalidOrderStatus", status, err)
		}
		if repo.updateCalls != 0 {
			t.Errorf("status %q: repo.UpdateStatus called %d times, want 0", status, repo.updateCalls)
		}
	}
}

// fakeProductPricer serves fixed storefront products to CreateOrder.
type fakeProductPricer struct {
	products map[string]ProductPublic
}

func (f *fakeProductPricer) GetProduct(ctx context.Context, id string, includeInactive bool) (ProductPublic, bool, error) {
	p, ok := f.products[id]
	return p, ok, nil
}

const testProductID = "vinh-quy-bai-to"

func newCreateOrderUsecase(product ProductPublic) (*OrderUsecase, *fakeOrderRepo) {
	repo := &fakeOrderRepo{}
	pricer := &fakeProductPricer{products: map[string]ProductPublic{testProductID: product}}
	return NewOrderUsecase(repo, pricer), repo
}

func listProduct(basePrice int64) ProductPublic {
	return ProductPublic{
		Product: domain.Product{ID: testProductID, Title: "Vinh Quy Bai To", BasePrice: basePrice, IsActive: true},
		Price:   basePrice,
	}
}

func orderRequest(item CreateOrderItem) CreateOrderRequest {
	return CreateOrderRequest{Phone: "0988000888", Items: []CreateOrderItem{item}}
}

func strPtr(s string) *string { return &s }

func TestCreateOrder_IgnoresClientUnitPrice(t *testing.T) {
	u, repo := newCreateOrderUsecase(listProduct(8500000))

	ord, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
		ProductID: testProductID,
		Quantity:  2,
		UnitPrice: 1,
	}))
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if len(repo.created) != 1 {
		t.Fatalf("created orders = %d, want 1", len(repo.created))
	}
	if got := ord.Items[0].UnitPrice; got != 8500000 {
		t.Errorf("unit price = %d, want 8500000 (server price)", got)
	}
	if got := ord.TotalAmount; got != 17000000 {
		t.Errorf("total = %d, want 17000000", got)
	}
}

func TestCreateOrder_SKUMatchWinsOverBasePrice(t *testing.T) {
	product := listProduct(8500000)
	product.SKUs = []domain.ProductSKU{
		{SizeCode: strPtr("L"), Attrs: map[string]string{"frame": "gold"}, Price: 9350000},
		{SizeCode: strPtr("L"), Attrs: map[string]string{"frame": "bronze"}, Price: 8900000},
	}
	u, _ := newCreateOrderUsecase(product)

	ord, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
		ProductID:     testProductID,
		SizeCode:      strPtr("L"),
		SelectedAttrs: map[string]string{"frame": "gold", "bg_tone": "red"},
		Quantity:      1,
	}))
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if got := ord.Items[0].UnitPrice; got != 9350000 {
		t.Errorf("unit price = %d, want 9350000 (matching SKU)", got)
	}
}

func TestCreateOrder_SKUAttrMissingFromSelectionDoesNotMatch(t *testing.T) {
	product := listProduct(8500000)
	product.SKUs = []domain.ProductSKU{
		{SizeCode: strPtr("L"), Attrs: map[string]string{"frame": ""}, Price: 9350000},
	}
	u, _ := newCreateOrderUsecase(product)

	// The SKU's frame is "" but the selection has no frame: that is not an exact match,
	// and the size-only match still returns the SKU price.
	ord, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
		ProductID: testProductID,
		SizeCode:  strPtr("L"),
		Quantity:  1,
	}))
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if got := ord.Items[0].UnitPrice; got != 9350000 {
		t.Errorf("unit price = %d, want 9350000 (size-only SKU match)", got)
	}
}

func TestCreateOrder_SizePriceUsedWhenNoSKUMatches(t *testing.T) {
	product := listProduct(8500000)
	product.Sizes = []domain.ProductSize{
		{SizeCode: "L", SizeLabel: "Lớn", Price: 7200000},
	}
	product.SKUs = []domain.ProductSKU{
		{SizeCode: strPtr("S"), Attrs: map[string]string{}, Price: 5000000},
	}
	u, _ := newCreateOrderUsecase(product)

	ord, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
		ProductID: testProductID,
		SizeCode:  strPtr("L"),
		Quantity:  1,
	}))
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if got := ord.Items[0].UnitPrice; got != 7200000 {
		t.Errorf("unit price = %d, want 7200000 (size price)", got)
	}
}

func TestCreateOrder_CampaignDiscountedPriceUsed(t *testing.T) {
	product := listProduct(8500000)
	sale := int64(6800000)
	product.Price = sale
	product.DiscountPrice = &sale
	u, _ := newCreateOrderUsecase(product)

	ord, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
		ProductID: testProductID,
		Quantity:  1,
		UnitPrice: 8500000,
	}))
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if got := ord.Items[0].UnitPrice; got != 6800000 {
		t.Errorf("unit price = %d, want 6800000 (campaign price)", got)
	}
}

func TestCreateOrder_ZeroOrInvalidPriceRejected(t *testing.T) {
	cases := map[string]ProductPublic{
		"zero base price": listProduct(0),
		"zero size price": func() ProductPublic {
			p := listProduct(8500000)
			p.Sizes = []domain.ProductSize{{SizeCode: "L", Price: 0}}
			p.Price = 0
			return p
		}(),
	}
	for name, product := range cases {
		u, repo := newCreateOrderUsecase(product)

		_, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
			ProductID: testProductID,
			SizeCode:  strPtr("L"),
			Quantity:  1,
			UnitPrice: 1,
		}))
		if err == nil {
			t.Errorf("%s: CreateOrder returned nil error, want rejection", name)
		}
		if len(repo.created) != 0 {
			t.Errorf("%s: %d orders created, want 0", name, len(repo.created))
		}
	}
}

func TestCreateOrder_ReturnsAndStoresLookupCode(t *testing.T) {
	u, repo := newCreateOrderUsecase(listProduct(8500000))

	ord, err := u.CreateOrder(context.Background(), orderRequest(CreateOrderItem{
		ProductID: testProductID,
		Quantity:  1,
	}))
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}
	if len(ord.LookupCode) != lookupCodeLen {
		t.Fatalf("lookup code = %q, want %d characters", ord.LookupCode, lookupCodeLen)
	}
	if len(repo.created) != 1 {
		t.Fatalf("created orders = %d, want 1", len(repo.created))
	}
	if repo.created[0].LookupCode != ord.LookupCode {
		t.Errorf("stored lookup code = %q, want the returned code %q", repo.created[0].LookupCode, ord.LookupCode)
	}
	if strings.ContainsAny(ord.LookupCode, "01IOL") {
		t.Errorf("lookup code %q contains an ambiguous character", ord.LookupCode)
	}
}
