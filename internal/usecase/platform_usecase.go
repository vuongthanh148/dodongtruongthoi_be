package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductQuery struct {
	Category string
	Sort     string
	Limit    int
	Offset   int
}

type CreateOrderRequest struct {
	Phone        string
	CustomerName *string
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

// PlatformUsecase handles all business logic by delegating to PostgreSQL repositories
type PlatformUsecase struct {
	// PostgreSQL repositories (REQUIRED)
	categoryRepo     domain.CategoryRepository
	productRepo      domain.ProductRepository
	productImageRepo domain.ProductImageRepository
	productSizeRepo  domain.ProductSizeRepository
	campaignRepo     domain.CampaignRepository
	reviewRepo       domain.ReviewRepository
	orderRepo        domain.OrderRepository
	wishlistRepo     domain.WishlistRepository
	bannerRepo       domain.BannerRepository
	contactRepo      domain.ContactLinkRepository
	adminUserRepo    domain.AdminUserRepository
	settingsRepo     domain.SiteSettingsRepository

	// Image uploader (optional - image uploads disabled if nil)
	imageUploader ImageUploader

	jwtSecret string
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
	ImageUploader    ImageUploader // optional
}

// NewPlatformUsecase creates a usecase that uses PostgreSQL repositories exclusively
// All 12 repositories MUST be provided; there is no in-memory fallback
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
		categoryRepo:     cfg.CategoryRepo,
		productRepo:      cfg.ProductRepo,
		productImageRepo: cfg.ProductImageRepo,
		productSizeRepo:  cfg.ProductSizeRepo,
		campaignRepo:     cfg.CampaignRepo,
		reviewRepo:       cfg.ReviewRepo,
		orderRepo:        cfg.OrderRepo,
		wishlistRepo:     cfg.WishlistRepo,
		bannerRepo:       cfg.BannerRepo,
		contactRepo:      cfg.ContactRepo,
		adminUserRepo:    cfg.AdminUserRepo,
		settingsRepo:     cfg.SettingsRepo,
		imageUploader:    cfg.ImageUploader,
		jwtSecret:        cfg.JWTSecret,
	}, nil
}

func (u *PlatformUsecase) GetPublicSettings(ctx context.Context) (map[string]string, error) {
	return u.settingsRepo.GetAll(ctx)
}

func (u *PlatformUsecase) GetAdminSettings(ctx context.Context) (map[string]string, error) {
	return u.settingsRepo.GetAll(ctx)
}

func (u *PlatformUsecase) UpdateSettings(ctx context.Context, in map[string]string) (map[string]string, error) {
	if err := u.settingsRepo.SetBulk(ctx, in); err != nil {
		return nil, err
	}
	return u.settingsRepo.GetAll(ctx)
}

func (u *PlatformUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx, false)
}

func (u *PlatformUsecase) GetCategory(ctx context.Context, id string) (domain.Category, bool, error) {
	return u.categoryRepo.Get(ctx, id, false)
}

func (u *PlatformUsecase) ListProducts(ctx context.Context, q ProductQuery, includeInactive bool) ([]ProductPublic, error) {
	products, err := u.productRepo.List(ctx, q.Category, includeInactive)
	if err != nil {
		return nil, err
	}

	items := make([]ProductPublic, 0, len(products))
	for _, p := range products {
		pub, err := u.buildProductPublic(ctx, p)
		if err != nil {
			return nil, err
		}
		items = append(items, pub)
	}

	switch q.Sort {
	case "price_asc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].Price < items[j].Price
		})
	case "price_desc":
		sort.Slice(items, func(i, j int) bool {
			return items[i].Price > items[j].Price
		})
	default:
		sort.Slice(items, func(i, j int) bool {
			if items[i].SortOrder == items[j].SortOrder {
				return items[i].CreatedAt.After(items[j].CreatedAt)
			}
			return items[i].SortOrder < items[j].SortOrder
		})
	}

	if q.Offset > 0 {
		if q.Offset >= len(items) {
			return []ProductPublic{}, nil
		}
		items = items[q.Offset:]
	}
	if q.Limit > 0 && q.Limit < len(items) {
		items = items[:q.Limit]
	}
	return items, nil
}

func (u *PlatformUsecase) GetProduct(ctx context.Context, id string, includeInactive bool) (ProductPublic, bool, error) {
	p, ok, err := u.productRepo.Get(ctx, id, includeInactive)
	if err != nil {
		return ProductPublic{}, false, err
	}
	if !ok {
		return ProductPublic{}, false, nil
	}
	pub, err := u.buildProductPublic(ctx, p)
	return pub, true, err
}

func (u *PlatformUsecase) ListReviews(ctx context.Context, productID string, includeUnapproved bool) ([]domain.Review, error) {
	reviews, err := u.reviewRepo.List(ctx, productID)
	if err != nil {
		return nil, err
	}

	out := make([]domain.Review, 0, len(reviews))
	for _, r := range reviews {
		if includeUnapproved || r.IsApproved {
			out = append(out, r)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	return out, nil
}

func (u *PlatformUsecase) CreateReview(ctx context.Context, productID, reviewerName string, rating int, body *string) (domain.Review, error) {
	if productID == "" {
		return domain.Review{}, errors.New("product_id is required")
	}
	if reviewerName == "" {
		return domain.Review{}, errors.New("reviewer_name is required")
	}
	if rating < 1 || rating > 5 {
		return domain.Review{}, errors.New("rating must be between 1 and 5")
	}

	product, ok, err := u.productRepo.Get(ctx, productID, false)
	if err != nil {
		return domain.Review{}, err
	}
	if !ok || !product.IsActive {
		return domain.Review{}, errors.New("product not found or inactive")
	}

	review := domain.Review{
		ID:           uuid.NewString(),
		ProductID:    productID,
		ReviewerName: reviewerName,
		Rating:       rating,
		Body:         body,
		IsApproved:   false,
		CreatedAt:    time.Now(),
	}

	return u.reviewRepo.Create(ctx, review)
}

func (u *PlatformUsecase) GetReviewByID(ctx context.Context, reviewID string) (domain.Review, bool, error) {
	reviews, err := u.reviewRepo.ListAll(ctx, nil)
	if err != nil {
		return domain.Review{}, false, err
	}

	for _, r := range reviews {
		if r.ID == reviewID {
			return r, true, nil
		}
	}
	return domain.Review{}, false, nil
}

func (u *PlatformUsecase) UpdateReviewApproval(ctx context.Context, reviewID string, isApproved bool) error {
	return u.reviewRepo.UpdateApproval(ctx, reviewID, isApproved)
}

func (u *PlatformUsecase) DeleteReview(ctx context.Context, reviewID string) error {
	return u.reviewRepo.Delete(ctx, reviewID)
}

func (u *PlatformUsecase) ListAllReviews(ctx context.Context, includeUnapproved bool) ([]domain.Review, error) {
	var approved *bool
	if !includeUnapproved {
		t := true
		approved = &t
	}
	reviews, err := u.reviewRepo.ListAll(ctx, approved)
	if err != nil {
		return nil, err
	}
	sort.Slice(reviews, func(i, j int) bool { return reviews[i].CreatedAt.After(reviews[j].CreatedAt) })
	return reviews, nil
}

func (u *PlatformUsecase) ListBanners(ctx context.Context) ([]domain.Banner, error) {
	return u.bannerRepo.List(ctx, false)
}

func (u *PlatformUsecase) ListContacts(ctx context.Context) ([]domain.ContactLink, error) {
	return u.contactRepo.List(ctx, false)
}

func (u *PlatformUsecase) SyncWishlist(ctx context.Context, phone string, productIDs []string) ([]domain.WishlistItem, error) {
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

func (u *PlatformUsecase) GetWishlist(ctx context.Context, phone string) ([]domain.WishlistItem, error) {
	return u.wishlistRepo.GetByPhone(ctx, phone)
}

func (u *PlatformUsecase) DeleteWishlistItem(ctx context.Context, phone, productID string) error {
	return u.wishlistRepo.Remove(ctx, phone, productID)
}

func (u *PlatformUsecase) CreateOrder(ctx context.Context, req CreateOrderRequest) (domain.Order, error) {
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
			public, err := u.buildProductPublic(ctx, product)
			if err != nil {
				return domain.Order{}, err
			}
			price = public.Price
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
		Note:         req.Note,
		Status:       "pending_confirm",
		TotalAmount:  total,
		Items:        items,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return u.orderRepo.Create(ctx, order)
}

func (u *PlatformUsecase) ListOrdersByPhone(ctx context.Context, phone string) ([]domain.Order, error) {
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
	}

	return orders, nil
}

func (u *PlatformUsecase) ListOrders(ctx context.Context, status *string) ([]domain.Order, error) {
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
	}

	return orders, nil
}

func (u *PlatformUsecase) GetOrder(ctx context.Context, id string) (domain.Order, bool, error) {
	ord, ok, err := u.orderRepo.Get(ctx, id)
	if err != nil || !ok {
		return ord, ok, err
	}

	items, getItemsErr := u.orderRepo.GetItems(ctx, id)
	if getItemsErr != nil {
		return domain.Order{}, false, getItemsErr
	}
	ord.Items = items

	return ord, true, nil
}

func (u *PlatformUsecase) UpdateOrderStatus(ctx context.Context, id, status, adminNote string) error {
	return u.orderRepo.UpdateStatus(ctx, id, status, &adminNote)
}

func (u *PlatformUsecase) Login(ctx context.Context, username, password string) (LoginResult, error) {
	user, ok, err := u.adminUserRepo.GetByUsername(ctx, username)
	if err != nil {
		return LoginResult{}, err
	}
	if !ok || !user.IsActive {
		return LoginResult{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return LoginResult{}, errors.New("invalid credentials")
	}

	// Update last login
	_ = u.adminUserRepo.UpdateLastLogin(ctx, user.ID)

	token, err := u.signAccessToken(username)
	if err != nil {
		return LoginResult{}, err
	}
	refreshToken, err := u.signRefreshToken(username)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token, RefreshToken: refreshToken}, nil
}

func (u *PlatformUsecase) Refresh(ctx context.Context, refreshToken string) (LoginResult, error) {
	username, err := u.verifyTokenByType(ctx, refreshToken, "refresh")
	if err != nil {
		return LoginResult{}, err
	}

	token, err := u.signAccessToken(username)
	if err != nil {
		return LoginResult{}, err
	}
	rotatedRefreshToken, err := u.signRefreshToken(username)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token, RefreshToken: rotatedRefreshToken}, nil
}

func (u *PlatformUsecase) VerifyToken(ctx context.Context, token string) error {
	_, err := u.verifyTokenByType(ctx, token, "access")
	return err
}

func (u *PlatformUsecase) verifyTokenByType(ctx context.Context, token, expectedType string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", errors.New("invalid token")
	}

	payloadRaw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", errors.New("invalid token")
	}
	sigRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", errors.New("invalid token")
	}

	mac := hmac.New(sha256.New, []byte(u.jwtSecret))
	_, _ = mac.Write(payloadRaw)
	expected := mac.Sum(nil)
	if !hmac.Equal(sigRaw, expected) {
		return "", errors.New("invalid token")
	}

	var payload struct {
		Username string `json:"u"`
		Type     string `json:"t"`
		Exp      int64  `json:"exp"`
	}
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		return "", errors.New("invalid token")
	}
	if payload.Type != expectedType {
		return "", errors.New("invalid token type")
	}
	if payload.Exp < time.Now().Unix() {
		return "", errors.New("token expired")
	}

	// Verify username exists
	_, ok, err := u.adminUserRepo.GetByUsername(ctx, payload.Username)
	if err != nil {
		return "", errors.New("invalid token user")
	}
	if !ok {
		return "", errors.New("invalid token user")
	}
	return payload.Username, nil
}

func (u *PlatformUsecase) CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error) {
	if p.ID == "" || p.Title == "" || p.CategoryID == "" {
		return domain.Product{}, errors.New("id, title, and category_id are required")
	}

	// Verify category exists
	_, ok, err := u.categoryRepo.Get(ctx, p.CategoryID, true)
	if err != nil {
		return domain.Product{}, err
	}
	if !ok {
		return domain.Product{}, errors.New("category not found")
	}

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.IsActive == false && p.ID != "" {
		// Only set to false if explicitly false, otherwise default to true
	} else {
		p.IsActive = true
	}

	return u.productRepo.Create(ctx, p)
}

func (u *PlatformUsecase) UpdateProduct(ctx context.Context, id string, updates domain.Product) (domain.Product, error) {
	if id == "" {
		return domain.Product{}, errors.New("product id is required")
	}

	// Get existing product
	existing, ok, err := u.productRepo.Get(ctx, id, true)
	if err != nil {
		return domain.Product{}, err
	}
	if !ok {
		return domain.Product{}, errors.New("product not found")
	}

	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Subtitle != nil && *updates.Subtitle != "" {
		existing.Subtitle = updates.Subtitle
	}
	if updates.CategoryID != "" {
		// Verify category exists
		_, ok, err := u.categoryRepo.Get(ctx, updates.CategoryID, true)
		if err != nil {
			return domain.Product{}, err
		}
		if !ok {
			return domain.Product{}, errors.New("category not found")
		}
		existing.CategoryID = updates.CategoryID
	}
	if updates.Badge != nil && *updates.Badge != "" {
		existing.Badge = updates.Badge
	}
	if updates.BasePrice > 0 {
		existing.BasePrice = updates.BasePrice
	}
	if updates.Description != nil && *updates.Description != "" {
		existing.Description = updates.Description
	}
	if updates.Meaning != nil && *updates.Meaning != "" {
		existing.Meaning = updates.Meaning
	}
	if updates.DefaultBG != "" {
		existing.DefaultBG = updates.DefaultBG
	}
	if updates.DefaultFrame != "" {
		existing.DefaultFrame = updates.DefaultFrame
	}
	if len(updates.BGTones) > 0 {
		existing.BGTones = updates.BGTones
	}
	if len(updates.Frames) > 0 {
		existing.Frames = updates.Frames
	}
	if len(updates.ZodiacIDs) > 0 {
		existing.ZodiacIDs = updates.ZodiacIDs
	}
	if len(updates.PurposePlace) > 0 {
		existing.PurposePlace = updates.PurposePlace
	}
	if len(updates.PurposeUse) > 0 {
		existing.PurposeUse = updates.PurposeUse
	}
	if len(updates.PurposeAvoid) > 0 {
		existing.PurposeAvoid = updates.PurposeAvoid
	}
	if updates.Specs != nil && len(updates.Specs) > 0 {
		existing.Specs = updates.Specs
	}
	existing.RequiresBGTone = updates.RequiresBGTone
	existing.RequiresFrame = updates.RequiresFrame
	existing.RequiresSize = updates.RequiresSize
	existing.IsActive = updates.IsActive
	existing.SortOrder = updates.SortOrder
	existing.UpdatedAt = time.Now()

	return u.productRepo.Update(ctx, id, existing)
}

func (u *PlatformUsecase) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product id is required")
	}
	return u.productRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) SetProductSizes(ctx context.Context, productID string, sizes []domain.ProductSize) ([]domain.ProductSize, error) {
	if productID == "" {
		return nil, errors.New("product id is required")
	}

	// Verify product exists
	_, ok, err := u.productRepo.Get(ctx, productID, true)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("product not found")
	}

	// Delete existing sizes
	if err := u.productSizeRepo.DeleteByProduct(ctx, productID); err != nil {
		return nil, err
	}

	// Create new sizes
	if len(sizes) == 0 {
		return []domain.ProductSize{}, nil
	}

	result := make([]domain.ProductSize, 0, len(sizes))
	for i := range sizes {
		sizes[i].ProductID = productID
		if sizes[i].ID == "" {
			sizes[i].ID = fmt.Sprintf("%s-size-%d", productID, i)
		}
		created, err := u.productSizeRepo.Create(ctx, sizes[i])
		if err != nil {
			return nil, err
		}
		result = append(result, created)
	}

	return result, nil
}

func (u *PlatformUsecase) CreateCategory(ctx context.Context, c domain.Category) (domain.Category, error) {
	if c.ID == "" || c.Name == "" {
		return domain.Category{}, errors.New("id and name are required")
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.IsActive = true

	return u.categoryRepo.Create(ctx, c)
}

func (u *PlatformUsecase) UpdateCategory(ctx context.Context, id string, updates domain.Category) (domain.Category, error) {
	if id == "" {
		return domain.Category{}, errors.New("category id is required")
	}

	// Get existing category
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
	if updates.Slug != "" {
		existing.Slug = updates.Slug
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

func (u *PlatformUsecase) DeleteCategory(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category id is required")
	}
	return u.categoryRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) CreateCampaign(ctx context.Context, c domain.Campaign) (domain.Campaign, error) {
	if c.ID == "" || c.Name == "" {
		return domain.Campaign{}, errors.New("id and name are required")
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if c.IsActive == false && c.ID != "" {
		// Keep explicit false
	} else {
		c.IsActive = true
	}

	return u.campaignRepo.Create(ctx, c)
}

func (u *PlatformUsecase) UpdateCampaign(ctx context.Context, id string, updates domain.Campaign) (domain.Campaign, error) {
	if id == "" {
		return domain.Campaign{}, errors.New("campaign id is required")
	}

	// Get existing campaign
	existing, ok, err := u.campaignRepo.Get(ctx, id)
	if err != nil {
		return domain.Campaign{}, err
	}
	if !ok {
		return domain.Campaign{}, errors.New("campaign not found")
	}

	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.Description != nil {
		existing.Description = updates.Description
	}
	if updates.DiscountType != "" {
		existing.DiscountType = updates.DiscountType
	}
	if updates.DiscountValue > 0 {
		existing.DiscountValue = updates.DiscountValue
	}
	if !updates.StartsAt.IsZero() {
		existing.StartsAt = updates.StartsAt
	}
	if !updates.EndsAt.IsZero() {
		existing.EndsAt = updates.EndsAt
	}
	existing.IsActive = updates.IsActive
	existing.UpdatedAt = time.Now()

	return u.campaignRepo.Update(ctx, id, existing)
}

func (u *PlatformUsecase) DeleteCampaign(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("campaign id is required")
	}
	return u.campaignRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) SetCampaignProducts(ctx context.Context, campaignID string, productIDs []string) error {
	if campaignID == "" {
		return errors.New("campaign id is required")
	}
	return u.campaignRepo.SetProducts(ctx, campaignID, productIDs)
}

func (u *PlatformUsecase) CreateBanner(ctx context.Context, b domain.Banner) (domain.Banner, error) {
	if b.ID == "" || b.Title == nil || *b.Title == "" {
		return domain.Banner{}, errors.New("id and title are required")
	}

	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	b.IsActive = true

	return u.bannerRepo.Create(ctx, b)
}

func (u *PlatformUsecase) UpdateBanner(ctx context.Context, id string, updates domain.Banner) (domain.Banner, error) {
	if id == "" {
		return domain.Banner{}, errors.New("banner id is required")
	}

	// Get existing banner
	existing, ok, err := u.bannerRepo.Get(ctx, id)
	if err != nil {
		return domain.Banner{}, err
	}
	if !ok {
		return domain.Banner{}, errors.New("banner not found")
	}

	if updates.Title != nil && *updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Subtitle != nil && *updates.Subtitle != "" {
		existing.Subtitle = updates.Subtitle
	}
	if updates.ImageURL != nil && *updates.ImageURL != "" {
		existing.ImageURL = updates.ImageURL
	}
	if updates.LinkURL != nil && *updates.LinkURL != "" {
		existing.LinkURL = updates.LinkURL
	}
	existing.SortOrder = updates.SortOrder
	existing.IsActive = updates.IsActive
	existing.UpdatedAt = time.Now()

	return u.bannerRepo.Update(ctx, id, existing)
}

func (u *PlatformUsecase) DeleteBanner(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("banner id is required")
	}
	return u.bannerRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) CreateContact(ctx context.Context, cl domain.ContactLink) (domain.ContactLink, error) {
	if cl.ID == "" || cl.Platform == "" || cl.URL == "" {
		return domain.ContactLink{}, errors.New("id, platform, and url are required")
	}

	now := time.Now()
	cl.CreatedAt = now
	cl.IsActive = true

	return u.contactRepo.Create(ctx, cl)
}

func (u *PlatformUsecase) UpdateContact(ctx context.Context, id string, updates domain.ContactLink) (domain.ContactLink, error) {
	if id == "" {
		return domain.ContactLink{}, errors.New("contact id is required")
	}

	// Get existing contact
	existing, ok, err := u.contactRepo.Get(ctx, id)
	if err != nil {
		return domain.ContactLink{}, err
	}
	if !ok {
		return domain.ContactLink{}, errors.New("contact not found")
	}

	if updates.Platform != "" {
		existing.Platform = updates.Platform
	}
	if updates.Label != "" {
		existing.Label = updates.Label
	}
	if updates.URL != "" {
		existing.URL = updates.URL
	}
	existing.SortOrder = updates.SortOrder
	existing.IsActive = updates.IsActive

	return u.contactRepo.Update(ctx, id, existing)
}

func (u *PlatformUsecase) DeleteContact(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("contact id is required")
	}
	return u.contactRepo.Delete(ctx, id)
}

func (u *PlatformUsecase) ListAllProducts(ctx context.Context) ([]ProductPublic, error) {
	products, err := u.productRepo.List(ctx, "", true)
	if err != nil {
		return nil, err
	}

	items := make([]ProductPublic, 0, len(products))
	for _, p := range products {
		pub, err := u.buildProductPublic(ctx, p)
		if err != nil {
			return nil, err
		}
		items = append(items, pub)
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].SortOrder == items[j].SortOrder {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].SortOrder < items[j].SortOrder
	})
	return items, nil
}

func (u *PlatformUsecase) ListAllCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx, true)
}

func (u *PlatformUsecase) ListAllCampaigns(ctx context.Context) ([]domain.Campaign, error) {
	return u.campaignRepo.List(ctx, true)
}

func (u *PlatformUsecase) ListAllBanners(ctx context.Context) ([]domain.Banner, error) {
	return u.bannerRepo.List(ctx, true)
}

func (u *PlatformUsecase) ListAllContacts(ctx context.Context) ([]domain.ContactLink, error) {
	return u.contactRepo.List(ctx, true)
}

func (u *PlatformUsecase) AddProductImage(ctx context.Context, productID string, bgTone, frame *string, url, altText string) (domain.ProductImage, error) {
	// Verify product exists
	_, ok, err := u.productRepo.Get(ctx, productID, true)
	if err != nil {
		return domain.ProductImage{}, err
	}
	if !ok {
		return domain.ProductImage{}, errors.New("product not found")
	}

	// Get current images to determine sort order
	images, err := u.productImageRepo.ListByProduct(ctx, productID)
	if err != nil {
		return domain.ProductImage{}, err
	}

	alt := altText
	img := domain.ProductImage{
		ID:        fmt.Sprintf("%s-img-%d", productID, len(images)+1),
		ProductID: productID,
		BGTone:    bgTone,
		Frame:     frame,
		URL:       url,
		AltText:   &alt,
		SortOrder: len(images),
		CreatedAt: time.Now(),
	}

	return u.productImageRepo.Create(ctx, img)
}

// UploadProductImage uploads a file to Cloudinary and saves image metadata to database
func (u *PlatformUsecase) UploadProductImage(ctx context.Context, productID string, file interface{}, filename string, bgTone *string, frame *string) (domain.ProductImage, error) {
	if u.imageUploader == nil {
		return domain.ProductImage{}, errors.New("image uploads disabled; cloudinary not configured")
	}

	// Verify product exists
	_, ok, err := u.productRepo.Get(ctx, productID, true)
	if err != nil {
		return domain.ProductImage{}, fmt.Errorf("failed to get product: %w", err)
	}
	if !ok {
		return domain.ProductImage{}, errors.New("product not found")
	}

	// Upload to Cloudinary
	url, err := u.imageUploader.UploadImage(ctx, file, filename, "products")
	if err != nil {
		return domain.ProductImage{}, fmt.Errorf("cloudinary upload failed: %w", err)
	}

	// Get current images to determine sort order
	images, err := u.productImageRepo.ListByProduct(ctx, productID)
	if err != nil {
		return domain.ProductImage{}, fmt.Errorf("failed to list images: %w", err)
	}

	// Create image record in database
	alt := filename
	img := domain.ProductImage{
		ProductID: productID,
		BGTone:    bgTone,
		Frame:     frame,
		URL:       url,
		AltText:   &alt,
		SortOrder: len(images),
		CreatedAt: time.Now(),
	}

	result, err := u.productImageRepo.Create(ctx, img)
	if err != nil {
		return domain.ProductImage{}, fmt.Errorf("failed to create image record: %w", err)
	}

	return result, nil
}

func (u *PlatformUsecase) DeleteProductImage(ctx context.Context, productID, imageID string) error {
	// Verify product exists
	_, ok, err := u.productRepo.Get(ctx, productID, true)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("product not found")
	}

	// Verify image exists and belongs to product
	img, ok, err := u.productImageRepo.Get(ctx, imageID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("image not found")
	}
	if img.ProductID != productID {
		return errors.New("image not found")
	}

	// Delete from Cloudinary if uploader configured
	if u.imageUploader != nil {
		// Extract public_id from Cloudinary URL
		parts := strings.Split(img.URL, "/upload/")
		if len(parts) == 2 {
			publicID := strings.TrimPrefix(parts[1], "v")
			// Skip version if present
			if idx := strings.Index(publicID, "/"); idx > 0 {
				publicID = publicID[idx+1:]
			}
			// Attempt deletion but don't fail if it errors
			_ = u.imageUploader.DeleteImage(ctx, publicID)
		}
	}

	return u.productImageRepo.Delete(ctx, imageID)
}

func (u *PlatformUsecase) buildProductPublic(ctx context.Context, p domain.Product) (ProductPublic, error) {
	price := p.BasePrice
	var discountPrice *int64

	if dp, ok, err := u.calculateDiscountPrice(ctx, p.ID, p.BasePrice); err == nil && ok {
		price = dp
		discountPrice = &dp
	}

	// Get reviews for this product
	reviews, err := u.reviewRepo.List(ctx, p.ID)
	if err != nil {
		return ProductPublic{}, err
	}

	// Get sizes for this product
	sizes, err := u.productSizeRepo.ListByProduct(ctx, p.ID)
	if err != nil {
		return ProductPublic{}, err
	}

	// Get images for this product
	images, err := u.productImageRepo.ListByProduct(ctx, p.ID)
	if err != nil {
		return ProductPublic{}, err
	}
	// Ensure images is never nil for JSON marshaling
	if images == nil {
		images = []domain.ProductImage{}
	}

	rating, count := aggregateRating(reviews)
	return ProductPublic{
		Product:       p,
		Price:         price,
		DiscountPrice: discountPrice,
		Sizes:         sizes,
		Images:        images,
		Rating:        rating,
		ReviewCount:   count,
	}, nil
}

func (u *PlatformUsecase) calculateDiscountPrice(ctx context.Context, productID string, basePrice int64) (int64, bool, error) {
	now := time.Now()
	best := int64(0)
	found := false

	campaigns, err := u.campaignRepo.List(ctx, true)
	if err != nil {
		return 0, false, err
	}

	for _, campaign := range campaigns {
		if !campaign.IsActive {
			continue
		}
		if now.Before(campaign.StartsAt) || now.After(campaign.EndsAt) {
			continue
		}

		productIDs, err := u.campaignRepo.GetProductIDs(ctx, campaign.ID)
		if err != nil {
			return 0, false, err
		}

		found := false
		for _, pid := range productIDs {
			if pid == productID {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		candidate := basePrice
		switch campaign.DiscountType {
		case "percentage":
			candidate = basePrice - (basePrice * campaign.DiscountValue / 100)
		case "fixed_amount":
			candidate = basePrice - campaign.DiscountValue
		}
		if candidate < 0 {
			candidate = 0
		}
		if !found || candidate < best {
			best = candidate
			found = true
		}
	}
	return best, found, nil
}

func (u *PlatformUsecase) signAccessToken(username string) (string, error) {
	return u.signToken(username, "access", 24*time.Hour)
}

func (u *PlatformUsecase) signRefreshToken(username string) (string, error) {
	return u.signToken(username, "refresh", 30*24*time.Hour)
}

func (u *PlatformUsecase) signToken(username, tokenType string, ttl time.Duration) (string, error) {
	payload := struct {
		Username string `json:"u"`
		Type     string `json:"t"`
		Exp      int64  `json:"exp"`
	}{
		Username: username,
		Type:     tokenType,
		Exp:      time.Now().Add(ttl).Unix(),
	}
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(u.jwtSecret))
	_, _ = mac.Write(payloadRaw)
	sig := mac.Sum(nil)

	return fmt.Sprintf("%s.%s",
		base64.RawURLEncoding.EncodeToString(payloadRaw),
		base64.RawURLEncoding.EncodeToString(sig),
	), nil
}

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
