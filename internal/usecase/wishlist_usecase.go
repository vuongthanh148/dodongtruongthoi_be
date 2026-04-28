package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type WishlistUsecase struct {
	wishlistRepo domain.WishlistRepository
}

func NewWishlistUsecase(wishlistRepo domain.WishlistRepository) *WishlistUsecase {
	return &WishlistUsecase{
		wishlistRepo: wishlistRepo,
	}
}

func (u *WishlistUsecase) SyncWishlist(ctx context.Context, phone string, productIDs []string) ([]domain.WishlistItem, error) {
	// Get existing wishlist items
	existingItems, _ := u.wishlistRepo.GetByPhone(ctx, phone)
	existingMap := make(map[string]bool)
	for _, item := range existingItems {
		existingMap[item.ProductID] = true
	}

	// Add new items that don't exist yet
	for _, id := range productIDs {
		if existingMap[id] {
			continue
		}
		item := domain.WishlistItem{
			ID:        uuid.NewString(),
			Phone:     phone,
			ProductID: id,
			CreatedAt: time.Now(),
		}
		_, err := u.wishlistRepo.Add(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	// Get updated wishlist
	items, err := u.wishlistRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
	return items, nil
}

func (u *WishlistUsecase) GetWishlist(ctx context.Context, phone string) ([]domain.WishlistItem, error) {
	return u.wishlistRepo.GetByPhone(ctx, phone)
}

func (u *WishlistUsecase) DeleteWishlistItem(ctx context.Context, phone, productID string) error {
	return u.wishlistRepo.Remove(ctx, phone, productID)
}
