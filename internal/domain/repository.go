package domain

import "context"

// CategoryRepository defines all category data operations
type CategoryRepository interface {
	List(ctx context.Context, includeInactive bool) ([]Category, error)
	Get(ctx context.Context, id string, includeInactive bool) (Category, bool, error)
	Create(ctx context.Context, c Category) (Category, error)
	Update(ctx context.Context, id string, c Category) (Category, error)
	Delete(ctx context.Context, id string) error
}

// ProductRepository defines all product data operations including images and sizes
type ProductRepository interface {
	List(ctx context.Context, category string, includeInactive bool) ([]Product, error)
	Get(ctx context.Context, id string, includeInactive bool) (Product, bool, error)
	Create(ctx context.Context, p Product) (Product, error)
	Update(ctx context.Context, id string, p Product) (Product, error)
	Delete(ctx context.Context, id string) error
}

// ProductImageRepository defines product image data operations
type ProductImageRepository interface {
	ListByProduct(ctx context.Context, productID string) ([]ProductImage, error)
	Get(ctx context.Context, id string) (ProductImage, bool, error)
	Create(ctx context.Context, img ProductImage) (ProductImage, error)
	Update(ctx context.Context, id string, img ProductImage) (ProductImage, error)
	Delete(ctx context.Context, id string) error
	DeleteByProduct(ctx context.Context, productID string) error
}

// ProductSizeRepository defines product size data operations
type ProductSizeRepository interface {
	ListByProduct(ctx context.Context, productID string) ([]ProductSize, error)
	Get(ctx context.Context, id string) (ProductSize, bool, error)
	Create(ctx context.Context, s ProductSize) (ProductSize, error)
	Update(ctx context.Context, id string, s ProductSize) (ProductSize, error)
	Delete(ctx context.Context, id string) error
	DeleteByProduct(ctx context.Context, productID string) error
}

// CampaignRepository defines campaign data operations
type CampaignRepository interface {
	List(ctx context.Context, includeInactive bool) ([]Campaign, error)
	Get(ctx context.Context, id string) (Campaign, bool, error)
	Create(ctx context.Context, c Campaign) (Campaign, error)
	Update(ctx context.Context, id string, c Campaign) (Campaign, error)
	Delete(ctx context.Context, id string) error
	SetProducts(ctx context.Context, campaignID string, productIDs []string) error
	GetProductIDs(ctx context.Context, campaignID string) ([]string, error)
}

// ReviewRepository defines review data operations
type ReviewRepository interface {
	List(ctx context.Context, productID string) ([]Review, error)
	ListAll(ctx context.Context, approved *bool) ([]Review, error)
	Create(ctx context.Context, rev Review) (Review, error)
	UpdateApproval(ctx context.Context, id string, isApproved bool) error
	Delete(ctx context.Context, id string) error
}

// OrderRepository defines order data operations
type OrderRepository interface {
	List(ctx context.Context, phone *string, status *string, limit int, offset int) ([]Order, error)
	Get(ctx context.Context, id string) (Order, bool, error)
	Create(ctx context.Context, ord Order) (Order, error)
	UpdateStatus(ctx context.Context, id string, status string, adminNote *string) error
	AddItem(ctx context.Context, item OrderItem) (OrderItem, error)
	GetItems(ctx context.Context, orderID string) ([]OrderItem, error)
}

// WishlistRepository defines wishlist data operations
type WishlistRepository interface {
	GetByPhone(ctx context.Context, phone string) ([]WishlistItem, error)
	Add(ctx context.Context, w WishlistItem) (WishlistItem, error)
	Remove(ctx context.Context, phone string, productID string) error
	RemoveByPhone(ctx context.Context, phone string) error
}

// BannerRepository defines banner data operations
type BannerRepository interface {
	List(ctx context.Context, includeInactive bool) ([]Banner, error)
	Get(ctx context.Context, id string) (Banner, bool, error)
	Create(ctx context.Context, b Banner) (Banner, error)
	Update(ctx context.Context, id string, b Banner) (Banner, error)
	Delete(ctx context.Context, id string) error
}

// ContactLinkRepository defines contact link data operations
type ContactLinkRepository interface {
	List(ctx context.Context, includeInactive bool) ([]ContactLink, error)
	Get(ctx context.Context, id string) (ContactLink, bool, error)
	Create(ctx context.Context, c ContactLink) (ContactLink, error)
	Update(ctx context.Context, id string, c ContactLink) (ContactLink, error)
	Delete(ctx context.Context, id string) error
}

// AdminUserRepository defines admin user data operations
type AdminUserRepository interface {
	GetByUsername(ctx context.Context, username string) (AdminUser, bool, error)
	Get(ctx context.Context, id string) (AdminUser, bool, error)
	Create(ctx context.Context, user AdminUser) (AdminUser, error)
	UpdateLastLogin(ctx context.Context, id string) error
	UpdatePassword(ctx context.Context, id string, passwordHash string) error
}

// SiteSettingsRepository defines site settings data operations
type SiteSettingsRepository interface {
	Get(ctx context.Context, key string) (string, error)
	GetAll(ctx context.Context) (map[string]string, error)
	Set(ctx context.Context, key string, value string) error
	SetBulk(ctx context.Context, settings map[string]string) error
}

// CustomerPhotoRepository defines customer lifestyle photo operations
type CustomerPhotoRepository interface {
	List(ctx context.Context, includeInactive bool) ([]CustomerPhoto, error)
	Get(ctx context.Context, id string) (CustomerPhoto, bool, error)
	Create(ctx context.Context, p CustomerPhoto) (CustomerPhoto, error)
	Update(ctx context.Context, id string, p CustomerPhoto) (CustomerPhoto, error)
	Delete(ctx context.Context, id string) error
}
