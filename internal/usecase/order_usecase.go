package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

// ProductPricer returns a product with the prices the storefront shows, including
// active campaign discounts. ProductUsecase implements it.
type ProductPricer interface {
	GetProduct(ctx context.Context, id string, includeInactive bool) (ProductPublic, bool, error)
}

type OrderUsecase struct {
	orderRepo domain.OrderRepository
	pricer    ProductPricer
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	pricer ProductPricer,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
		pricer:    pricer,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, req CreateOrderRequest) (domain.Order, error) {
	phone := normalizeLookupPhone(strings.TrimSpace(req.Phone))
	if phone == "" {
		return domain.Order{}, errors.New("phone is required")
	}
	if len(req.Items) == 0 {
		return domain.Order{}, errors.New("items are required")
	}

	now := time.Now()
	orderID := uuid.NewString()
	lookupCode, err := newLookupCode()
	if err != nil {
		return domain.Order{}, err
	}
	items := make([]domain.OrderItem, 0, len(req.Items))
	var total int64

	for _, in := range req.Items {
		product, ok, err := u.pricer.GetProduct(ctx, in.ProductID, false)
		if err != nil {
			return domain.Order{}, err
		}
		if !ok || !product.IsActive {
			return domain.Order{}, fmt.Errorf("invalid product: %s", in.ProductID)
		}
		sizeCode := ""
		if in.SizeCode != nil {
			sizeCode = strings.TrimSpace(*in.SizeCode)
		}
		if product.RequiresSize && sizeCode == "" {
			return domain.Order{}, fmt.Errorf("size_code is required for product: %s", in.ProductID)
		}

		qty := in.Quantity
		if qty <= 0 {
			qty = 1
		}

		// The client's unitPrice is never read: the server computes every line price.
		price := resolveStorefrontPrice(product, sizeCode, in.SelectedAttrs)
		if price <= 0 {
			return domain.Order{}, fmt.Errorf("invalid price for product: %s", in.ProductID)
		}

		line := domain.OrderItem{
			ID:              uuid.NewString(),
			OrderID:         orderID,
			ProductID:       in.ProductID,
			ProductTitle:    product.Title,
			ProductSubtitle: product.Subtitle,
			SizeCode:        in.SizeCode,
			SizeLabel:       in.SizeLabel,
			SelectedAttrs:   in.SelectedAttrs,
			Quantity:        qty,
			UnitPrice:       price,
			VariantImageURL: in.VariantImageURL,
		}
		items = append(items, line)
		total += int64(qty) * price
	}

	order := domain.Order{
		ID:           orderID,
		Phone:        phone,
		CustomerName: req.CustomerName,
		Address:      req.Address,
		Note:         req.Note,
		Status:       "pending_confirm",
		TotalAmount:  total,
		Items:        items,
		CreatedAt:    now,
		UpdatedAt:    now,
		LookupCode:   lookupCode,
	}
	return u.orderRepo.Create(ctx, order)
}

func (u *OrderUsecase) ListOrdersByPhone(ctx context.Context, phone string) ([]domain.Order, error) {
	phone = normalizeLookupPhone(phone)
	orders, err := u.orderRepo.List(ctx, &phone, nil, 0, 0)
	if err != nil {
		return nil, err
	}

	for i, ord := range orders {
		items, getItemsErr := u.orderRepo.GetItems(ctx, ord.ID)
		if getItemsErr != nil {
			return nil, getItemsErr
		}
		orders[i].Items = items
		orders[i].TotalAmount = calculateOrderTotal(items)
	}

	return orders, nil
}

func (u *OrderUsecase) ListOrders(ctx context.Context, status *string) ([]domain.Order, error) {
	orders, err := u.orderRepo.List(ctx, nil, status, 0, 0)
	if err != nil {
		return nil, err
	}

	for i, ord := range orders {
		items, getItemsErr := u.orderRepo.GetItems(ctx, ord.ID)
		if getItemsErr != nil {
			return nil, getItemsErr
		}
		orders[i].Items = items
		orders[i].TotalAmount = calculateOrderTotal(items)
	}

	return orders, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (domain.Order, bool, error) {
	ord, ok, err := u.orderRepo.Get(ctx, id)
	if err != nil || !ok {
		return ord, ok, err
	}

	items, getItemsErr := u.orderRepo.GetItems(ctx, id)
	if getItemsErr != nil {
		return domain.Order{}, false, getItemsErr
	}
	ord.Items = items
	ord.TotalAmount = calculateOrderTotal(items)

	return ord, true, nil
}

// validOrderStatuses is the order lifecycle, matching the admin CMS and
// customer-facing order pages.
var validOrderStatuses = map[string]struct{}{
	"pending_confirm": {},
	"confirmed":       {},
	"processing":      {},
	"shipped":         {},
	"completed":       {},
	"cancelled":       {},
}

// ErrInvalidOrderStatus is returned when an order status is not in the lifecycle.
var ErrInvalidOrderStatus = errors.New("invalid order status")

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, id, status, adminNote string) error {
	if _, ok := validOrderStatuses[status]; !ok {
		return fmt.Errorf("%w: %q", ErrInvalidOrderStatus, status)
	}
	return u.orderRepo.UpdateStatus(ctx, id, status, &adminNote)
}

// orderNotCancellableError reports a cancel attempt on an order past confirmation.
// It unwraps to domain.ErrConflict so handlers map it to 409 with a clean message.
type orderNotCancellableError struct {
	status string
}

func (e orderNotCancellableError) Error() string {
	return fmt.Sprintf("order cannot be cancelled at status %s", e.status)
}

func (e orderNotCancellableError) Unwrap() error {
	return domain.ErrConflict
}

var cancellableOrderStatuses = []string{"pending_confirm", "confirmed"}

// CancelOrder lets a buyer cancel an order that has not moved past confirmation.
func (u *OrderUsecase) CancelOrder(ctx context.Context, id string) error {
	// Order ids are UUIDs; anything else cannot match an order, so report not found
	// instead of letting the database reject the value.
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: order not found", domain.ErrNotFound)
	}
	ord, ok, err := u.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: order not found", domain.ErrNotFound)
	}
	if !slices.Contains(cancellableOrderStatuses, ord.Status) {
		return orderNotCancellableError{status: ord.Status}
	}
	// Conditional update: a status change that lands between the read above and this write
	// makes the cancel a no-op instead of overwriting it. admin_note is left untouched.
	changed, err := u.orderRepo.TransitionStatus(ctx, id, cancellableOrderStatuses, "cancelled")
	if err != nil {
		return err
	}
	if changed {
		return nil
	}
	latest, ok, err := u.orderRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: order not found", domain.ErrNotFound)
	}
	return orderNotCancellableError{status: latest.Status}
}

// resolveStorefrontPrice mirrors the storefront price (src/lib/sku.ts resolveSKUPrice, with the
// fallback the product and cart pages pass: size price, then discount price, then list price).
// Product prices already have campaign discounts applied by ProductUsecase.
func resolveStorefrontPrice(product ProductPublic, sizeCode string, selectedAttrs map[string]string) int64 {
	// Exact match: size and every SKU attr present in the selection with the same value.
	// Extra selected attrs are display-only and ignored.
	for _, s := range product.SKUs {
		if skuSizeCodeMatches(s.SizeCode, sizeCode) && skuAttrsMatch(s.Attrs, selectedAttrs) && s.Price > 0 {
			return s.Price
		}
	}
	// Size-only match: ignore the variant attrs.
	for _, s := range product.SKUs {
		if skuSizeCodeMatches(s.SizeCode, sizeCode) && s.Price > 0 {
			return s.Price
		}
	}
	if sizeCode != "" {
		for _, size := range product.Sizes {
			if size.SizeCode == sizeCode && size.Price > 0 {
				return size.Price
			}
		}
	}
	if product.DiscountPrice != nil && *product.DiscountPrice > 0 {
		return *product.DiscountPrice
	}
	return product.Price
}

// skuSizeCodeMatches compares a SKU size code with the selected one. An empty selected code
// means "no size", which matches only a SKU with a null size code (as in the storefront's null check).
func skuSizeCodeMatches(skuSize *string, sizeCode string) bool {
	if skuSize == nil || sizeCode == "" {
		return skuSize == nil && sizeCode == ""
	}
	return *skuSize == sizeCode
}

// skuAttrsMatch reports whether every SKU attr is selected with the same value.
// A missing selected attr never matches, even when the SKU value is empty.
func skuAttrsMatch(attrs map[string]string, selectedAttrs map[string]string) bool {
	for k, v := range attrs {
		selected, ok := selectedAttrs[k]
		if !ok || selected != v {
			return false
		}
	}
	return true
}

func calculateOrderTotal(items []domain.OrderItem) int64 {
	var total int64
	for _, item := range items {
		if item.Quantity <= 0 || item.UnitPrice <= 0 {
			continue
		}
		total += int64(item.Quantity) * item.UnitPrice
	}
	return total
}
