package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductUsecase struct {
	productRepo      domain.ProductRepository
	productImageRepo domain.ProductImageRepository
	productSizeRepo  domain.ProductSizeRepository
	reviewRepo       domain.ReviewRepository
	campaignRepo     domain.CampaignRepository
	categoryRepo     domain.CategoryRepository
	imageUploader    ImageUploader
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	productImageRepo domain.ProductImageRepository,
	productSizeRepo domain.ProductSizeRepository,
	reviewRepo domain.ReviewRepository,
	campaignRepo domain.CampaignRepository,
	categoryRepo domain.CategoryRepository,
	imageUploader ImageUploader,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo:      productRepo,
		productImageRepo: productImageRepo,
		productSizeRepo:  productSizeRepo,
		reviewRepo:       reviewRepo,
		campaignRepo:     campaignRepo,
		categoryRepo:     categoryRepo,
		imageUploader:    imageUploader,
	}
}

func (u *ProductUsecase) ListProducts(ctx context.Context, q domain.ProductQuery, includeInactive bool) ([]ProductPublic, error) {
	products, err := u.productRepo.List(ctx, q, includeInactive)
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

func (u *ProductUsecase) GetProduct(ctx context.Context, id string, includeInactive bool) (ProductPublic, bool, error) {
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

func (u *ProductUsecase) ListAllProducts(ctx context.Context) ([]ProductPublic, error) {
	products, err := u.productRepo.List(ctx, domain.ProductQuery{}, true)
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

func (u *ProductUsecase) CreateProduct(ctx context.Context, p domain.Product) (domain.Product, error) {
	if strings.TrimSpace(p.ID) == "" {
		return domain.Product{}, fmt.Errorf("%w: product id is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(p.Title) == "" {
		return domain.Product{}, fmt.Errorf("%w: product title is required", domain.ErrInvalidInput)
	}
	if strings.TrimSpace(p.CategoryID) == "" {
		return domain.Product{}, fmt.Errorf("%w: category_id is required", domain.ErrInvalidInput)
	}
	if p.BasePrice < 0 {
		return domain.Product{}, fmt.Errorf("%w: base_price must not be negative", domain.ErrInvalidInput)
	}

	// Verify category exists
	_, ok, err := u.categoryRepo.Get(ctx, p.CategoryID, true)
	if err != nil {
		return domain.Product{}, err
	}
	if !ok {
		return domain.Product{}, fmt.Errorf("%w: category not found", domain.ErrNotFound)
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

func (u *ProductUsecase) UpdateProduct(ctx context.Context, id string, updates domain.Product) (domain.Product, error) {
	if strings.TrimSpace(id) == "" {
		return domain.Product{}, fmt.Errorf("%w: product id is required", domain.ErrInvalidInput)
	}

	// Get existing product
	existing, ok, err := u.productRepo.Get(ctx, id, true)
	if err != nil {
		return domain.Product{}, err
	}
	if !ok {
		return domain.Product{}, fmt.Errorf("%w: product not found", domain.ErrNotFound)
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

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("product id is required")
	}
	return u.productRepo.Delete(ctx, id)
}

func (u *ProductUsecase) SetProductSizes(ctx context.Context, productID string, sizes []domain.ProductSize) ([]domain.ProductSize, error) {
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

func (u *ProductUsecase) AddProductImage(ctx context.Context, productID string, bgTone, frame *string, url, altText string) (domain.ProductImage, error) {
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
func (u *ProductUsecase) UploadProductImage(ctx context.Context, productID string, file interface{}, filename string, bgTone *string, frame *string) (domain.ProductImage, error) {
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

func (u *ProductUsecase) DeleteProductImage(ctx context.Context, productID, imageID string) error {
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
		parts := split(img.URL, "/upload/")
		if len(parts) == 2 {
			publicID := parts[1]
			if idx := firstIndex(publicID, "/"); idx > 0 {
				publicID = publicID[idx+1:]
			}
			// Attempt deletion but don't fail if it errors
			_ = u.imageUploader.DeleteImage(ctx, publicID)
		}
	}

	return u.productImageRepo.Delete(ctx, imageID)
}

func (u *ProductUsecase) buildProductPublic(ctx context.Context, p domain.Product) (ProductPublic, error) {
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

func (u *ProductUsecase) calculateDiscountPrice(ctx context.Context, productID string, basePrice int64) (int64, bool, error) {
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

// Helper functions for string operations
func split(s, sep string) []string {
	result := make([]string, 0)
	idx := 0
	for {
		pos := findString(s[idx:], sep)
		if pos == -1 {
			result = append(result, s[idx:])
			break
		}
		result = append(result, s[idx:idx+pos])
		idx += pos + len(sep)
	}
	return result
}

func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

func firstIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}
