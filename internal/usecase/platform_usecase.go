package usecase

import (
	"context"
	"errors"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CreateOrderRequest struct {
	Phone        string
	CustomerName *string
	Address      *string
	Note         *string
	Items        []CreateOrderItem
}

type CreateOrderItem struct {
	ProductID       string
	SizeCode        *string
	SizeLabel       *string
	BGTone          *string
	BGToneLabel     *string
	Frame           *string
	FrameLabel      *string
	Quantity        int
	UnitPrice       int64
	VariantImageURL *string
}

type LoginResult struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ProductPublic struct {
	domain.Product
	Price         int64                 `json:"price"`
	DiscountPrice *int64                `json:"discount_price,omitempty"`
	Sizes         []domain.ProductSize  `json:"sizes"`
	Images        []domain.ProductImage `json:"images"`
	Rating        float64               `json:"rating"`
	ReviewCount   int                   `json:"review_count"`
}

// ImageUploader interface for uploading images
type ImageUploader interface {
	UploadImage(ctx context.Context, file interface{}, filename string, folder string) (url string, err error)
	DeleteImage(ctx context.Context, publicID string) error
}

// PlatformUsecase is an aggregator that embeds all domain-scoped usecases
// It provides a unified interface for handlers while keeping domain logic separated
type PlatformUsecase struct {
	*AuthUsecase
	*ProductUsecase
	*CategoryUsecase
	*OrderUsecase
	*CampaignUsecase
	*ReviewUsecase
	*BannerUsecase
	*CustomerPhotoUsecase
	*SettingsUsecase
	*ContactUsecase
	*WishlistUsecase
}

type PlatformUsecaseConfig struct {
	JWTSecret        string
	CategoryRepo     domain.CategoryRepository
	ProductRepo      domain.ProductRepository
	ProductImageRepo domain.ProductImageRepository
	ProductSizeRepo  domain.ProductSizeRepository
	CampaignRepo     domain.CampaignRepository
	ReviewRepo       domain.ReviewRepository
	OrderRepo        domain.OrderRepository
	WishlistRepo     domain.WishlistRepository
	BannerRepo       domain.BannerRepository
	ContactRepo      domain.ContactLinkRepository
	AdminUserRepo    domain.AdminUserRepository
	SettingsRepo     domain.SiteSettingsRepository
	CustomerPhotoRepo domain.CustomerPhotoRepository // optional - feature may be added later
	ImageUploader    ImageUploader                 // optional
}

// NewPlatformUsecase creates a usecase that uses PostgreSQL repositories exclusively
// All repositories MUST be provided; there is no in-memory fallback
func NewPlatformUsecase(cfg PlatformUsecaseConfig) (*PlatformUsecase, error) {
	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "dev-secret"
	}

	// Validate all repositories are provided
	if cfg.CategoryRepo == nil || cfg.ProductRepo == nil || cfg.OrderRepo == nil ||
		cfg.ReviewRepo == nil || cfg.WishlistRepo == nil || cfg.BannerRepo == nil ||
		cfg.ContactRepo == nil || cfg.AdminUserRepo == nil || cfg.SettingsRepo == nil ||
		cfg.ProductImageRepo == nil || cfg.ProductSizeRepo == nil || cfg.CampaignRepo == nil {
		return nil, errors.New("all repositories must be provided; in-memory fallback is not supported")
	}

	return &PlatformUsecase{
		AuthUsecase: NewAuthUsecase(
			cfg.AdminUserRepo,
			cfg.JWTSecret,
		),
		ProductUsecase: NewProductUsecase(
			cfg.ProductRepo,
			cfg.ProductImageRepo,
			cfg.ProductSizeRepo,
			cfg.ReviewRepo,
			cfg.CampaignRepo,
			cfg.CategoryRepo,
			cfg.ImageUploader,
		),
		CategoryUsecase: NewCategoryUsecase(
			cfg.CategoryRepo,
		),
		OrderUsecase: NewOrderUsecase(
			cfg.OrderRepo,
			cfg.ProductRepo,
		),
		CampaignUsecase: NewCampaignUsecase(
			cfg.CampaignRepo,
		),
		ReviewUsecase: NewReviewUsecase(
			cfg.ReviewRepo,
			cfg.ProductRepo,
		),
		BannerUsecase: NewBannerUsecase(
			cfg.BannerRepo,
		),
		CustomerPhotoUsecase: NewCustomerPhotoUsecase(
			cfg.CustomerPhotoRepo,
			cfg.ImageUploader,
		),
		SettingsUsecase: NewSettingsUsecase(
			cfg.SettingsRepo,
		),
		ContactUsecase: NewContactUsecase(
			cfg.ContactRepo,
		),
		WishlistUsecase: NewWishlistUsecase(
			cfg.WishlistRepo,
		),
	}, nil
}

// aggregateRating computes average rating and count from approved reviews
func aggregateRating(rows []domain.Review) (float64, int) {
	if len(rows) == 0 {
		return 0, 0
	}
	var total int
	var count int
	for _, r := range rows {
		if !r.IsApproved {
			continue
		}
		total += r.Rating
		count++
	}
	if count == 0 {
		return 0, 0
	}
	return float64(total) / float64(count), count
}

// The following helper method sorts products by campaign discount application
// (formerly part of platform_usecase's calculateDiscountPrice)
type discountInfo struct {
	productID string
	discount  int64
	applied   bool
}
