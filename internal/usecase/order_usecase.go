package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type OrderUsecase struct {
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
}

func NewOrderUsecase(
	orderRepo domain.OrderRepository,
	productRepo domain.ProductRepository,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, req CreateOrderRequest) (domain.Order, error) {
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		return domain.Order{}, errors.New("phone is required")
	}
	if len(req.Items) == 0 {
		return domain.Order{}, errors.New("items are required")
	}

	now := time.Now()
	orderID := uuid.NewString()
	items := make([]domain.OrderItem, 0, len(req.Items))
	var total int64

	for _, in := range req.Items {
		product, ok, err := u.productRepo.Get(ctx, in.ProductID, false)
		if err != nil {
			return domain.Order{}, err
		}
		if !ok || !product.IsActive {
			return domain.Order{}, fmt.Errorf("invalid product: %s", in.ProductID)
		}
		bgTone := ""
		if in.BGTone != nil {
			bgTone = strings.TrimSpace(*in.BGTone)
		}
		if product.RequiresBGTone && bgTone == "" {
			return domain.Order{}, fmt.Errorf("bg_tone is required for product: %s", in.ProductID)
		}
		frame := ""
		if in.Frame != nil {
			frame = strings.TrimSpace(*in.Frame)
		}
		if product.RequiresFrame && frame == "" {
			return domain.Order{}, fmt.Errorf("frame is required for product: %s", in.ProductID)
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

		price := in.UnitPrice
		if price <= 0 {
			// Cannot build ProductPublic without campaign repo, so just use base price
			price = product.BasePrice
		}

		line := domain.OrderItem{
			ID:              uuid.NewString(),
			OrderID:         orderID,
			ProductID:       in.ProductID,
			ProductTitle:    product.Title,
			ProductSubtitle: product.Subtitle,
			SizeCode:        in.SizeCode,
			SizeLabel:       in.SizeLabel,
			BGTone:          in.BGTone,
			BGToneLabel:     in.BGToneLabel,
			Frame:           in.Frame,
			FrameLabel:      in.FrameLabel,
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
	}
	return u.orderRepo.Create(ctx, order)
}

func (u *OrderUsecase) ListOrdersByPhone(ctx context.Context, phone string) ([]domain.Order, error) {
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

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, id, status, adminNote string) error {
	return u.orderRepo.UpdateStatus(ctx, id, status, &adminNote)
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
