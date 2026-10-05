package usecase

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ProductUsecase struct {
	productRepo      domain.ProductRepository
	productImageRepo domain.ProductImageRepository
	productSizeRepo  domain.ProductSizeRepository
	productSKURepo   domain.ProductSKURepository
	reviewRepo       domain.ReviewRepository
	campaignRepo     domain.CampaignRepository
	categoryRepo     domain.CategoryRepository
	imageUploader    ImageUploader
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	productImageRepo domain.ProductImageRepository,
	productSizeRepo domain.ProductSizeRepository,
	productSKURepo domain.ProductSKURepository,
	reviewRepo domain.ReviewRepository,
	campaignRepo domain.CampaignRepository,
	categoryRepo domain.CategoryRepository,
	imageUploader ImageUploader,
) *ProductUsecase {
	return &ProductUsecase{
		productRepo:      productRepo,
		productImageRepo: productImageRepo,
		productSizeRepo:  productSizeRepo,
		productSKURepo:   productSKURepo,
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
	p.IsActive = true

	return u.productRepo.Create(ctx, p)
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, id string, updates domain.Product, isActive *bool, requiresSize *bool) (domain.Product, error) {
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
	// When SKUs exist their min price is the product's base price; ignore the incoming value.
	skus, err := u.productSKURepo.ListByProduct(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	if minPrice, hasSKUs := minSKUPrice(skus); hasSKUs {
		existing.BasePrice = minPrice
	} else if updates.BasePrice > 0 {
		existing.BasePrice = updates.BasePrice
	}
	if updates.Description != nil && *updates.Description != "" {
		existing.Description = updates.Description
	}
	if updates.Meaning != nil && *updates.Meaning != "" {
		existing.Meaning = updates.Meaning
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
	if len(updates.VariantOptions) > 0 {
		existing.VariantOptions = updates.VariantOptions
	}
	if len(updates.DefaultVariant) > 0 {
		existing.DefaultVariant = updates.DefaultVariant
	}
	if requiresSize != nil {
		existing.RequiresSize = *requiresSize
	}
	if isActive != nil {
		existing.IsActive = *isActive
	}
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
			sizes[i].ID = uuid.New().String()
		}
		created, err := u.productSizeRepo.Create(ctx, sizes[i])
		if err != nil {
			return nil, err
		}
		result = append(result, created)
	}

	return result, nil
}

func (u *ProductUsecase) ListProductImages(ctx context.Context, productID string) ([]domain.ProductImage, error) {
	return u.productImageRepo.ListByProduct(ctx, productID)
}

// AttachImage attaches an existing library image to a product with optional variant attrs.
func (u *ProductUsecase) AttachImage(ctx context.Context, productID string, imageID string, attrs []domain.VariantAttr) (domain.ProductImage, error) {
	_, ok, err := u.productRepo.Get(ctx, productID, true)
	if err != nil {
		return domain.ProductImage{}, err
	}
	if !ok {
		return domain.ProductImage{}, errors.New("product not found")
	}

	existing, err := u.productImageRepo.ListByProduct(ctx, productID)
	if err != nil {
		return domain.ProductImage{}, err
	}

	img := domain.ProductImage{
		ID:        fmt.Sprintf("%s-pi-%d", productID, len(existing)+1),
		ProductID: productID,
		ImageID:   imageID,
		SortOrder: len(existing),
		Attrs:     attrs,
		CreatedAt: time.Now(),
	}
	return u.productImageRepo.Attach(ctx, img)
}

// SetProductImageAttrs replaces all variant attrs on a product_image.
func (u *ProductUsecase) SetProductImageAttrs(ctx context.Context, productImageID string, attrs []domain.VariantAttr) error {
	return u.productImageRepo.SetAttrs(ctx, productImageID, attrs)
}

// DetachProductImage removes the product→image link (does not delete from library).
func (u *ProductUsecase) DetachProductImage(ctx context.Context, productID, productImageID string) error {
	img, ok, err := u.productImageRepo.Get(ctx, productImageID)
	if err != nil {
		return err
	}
	if !ok || img.ProductID != productID {
		return errors.New("product image not found")
	}
	return u.productImageRepo.Detach(ctx, productImageID)
}

func (u *ProductUsecase) buildProductPublic(ctx context.Context, p domain.Product) (ProductPublic, error) {
	price := p.BasePrice
	var discountPrice *int64

	rule, hasRule, err := u.activeCampaignDiscount(ctx, p.ID, p.BasePrice)
	if err != nil {
		return ProductPublic{}, err
	}
	if hasRule {
		dp := rule.apply(p.BasePrice)
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
	if images == nil {
		images = []domain.ProductImage{}
	}

	// Get SKUs for this product
	skus, err := u.productSKURepo.ListByProduct(ctx, p.ID)
	if err != nil {
		return ProductPublic{}, err
	}
	if skus == nil {
		skus = []domain.ProductSKU{}
	}

	// Read the list price of the lowest SKU before any campaign rule is applied.
	listMin, hasSKUMin := minSKUPrice(skus)

	// The same campaign rule applies to every price the storefront reads.
	if hasRule {
		for i := range sizes {
			if sizes[i].Price > 0 {
				sizes[i].Price = rule.apply(sizes[i].Price)
			}
		}
		for i := range skus {
			if skus[i].Price > 0 {
				skus[i].Price = rule.apply(skus[i].Price)
			}
		}
	}

	// With SKUs, the list price is the lowest SKU price; a campaign sale price is applied on top of it.
	if hasSKUMin {
		price = listMin
		discountPrice = nil
		if hasRule {
			sale := rule.apply(listMin)
			discountPrice = &sale
		}
	}

	rating, count := aggregateRating(reviews)
	return ProductPublic{
		Product:       p,
		Price:         price,
		DiscountPrice: discountPrice,
		Sizes:         sizes,
		SKUs:          skus,
		Images:        images,
		Rating:        rating,
		ReviewCount:   count,
	}, nil
}

func (u *ProductUsecase) GetSKUs(ctx context.Context, productID string) ([]domain.ProductSKU, error) {
	return u.productSKURepo.ListByProduct(ctx, productID)
}

// SetSKUs replaces a product's SKUs and, when any are given, sets base_price to the min SKU price.
// The SKU table and product table are separate repos, so this is two sequential writes, not one transaction.
func (u *ProductUsecase) SetSKUs(ctx context.Context, productID string, skus []domain.ProductSKU) error {
	product, ok, err := u.productRepo.Get(ctx, productID, true)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%w: product not found", domain.ErrNotFound)
	}

	if err := u.productSKURepo.SetByProduct(ctx, productID, skus); err != nil {
		return err
	}

	minPrice, hasSKUs := minSKUPrice(skus)
	if !hasSKUs {
		return nil
	}
	product.BasePrice = minPrice
	_, err = u.productRepo.Update(ctx, productID, product)
	return err
}

// minSKUPrice returns the lowest SKU price and whether there is at least one SKU.
func minSKUPrice(skus []domain.ProductSKU) (int64, bool) {
	if len(skus) == 0 {
		return 0, false
	}
	minPrice := skus[0].Price
	for _, s := range skus[1:] {
		if s.Price < minPrice {
			minPrice = s.Price
		}
	}
	return minPrice, true
}

// campaignDiscount is one active campaign's discount rule.
type campaignDiscount struct {
	kind  string
	value int64
}

// apply returns price after this discount, never below zero.
func (d campaignDiscount) apply(price int64) int64 {
	var out int64
	switch d.kind {
	case "percentage":
		out = price - (price * d.value / 100)
	case "fixed_amount":
		out = price - d.value
	default:
		out = price
	}
	if out < 0 {
		return 0
	}
	return out
}

// activeCampaignDiscount returns the active campaign rule that gives the lowest price
// for basePrice. found is false when no active, in-window campaign includes the product.
func (u *ProductUsecase) activeCampaignDiscount(ctx context.Context, productID string, basePrice int64) (campaignDiscount, bool, error) {
	now := time.Now()
	var best campaignDiscount
	bestPrice := int64(0)
	found := false

	campaigns, err := u.campaignRepo.List(ctx, true)
	if err != nil {
		return campaignDiscount{}, false, err
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
			return campaignDiscount{}, false, err
		}

		inCampaign := false
		for _, pid := range productIDs {
			if pid == productID {
				inCampaign = true
				break
			}
		}
		if !inCampaign {
			continue
		}

		rule := campaignDiscount{kind: campaign.DiscountType, value: campaign.DiscountValue}
		candidate := rule.apply(basePrice)
		if !found || candidate < bestPrice {
			best = rule
			bestPrice = candidate
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
