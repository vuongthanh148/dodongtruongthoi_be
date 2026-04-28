package usecase

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ReviewUsecase struct {
	reviewRepo  domain.ReviewRepository
	productRepo domain.ProductRepository
}

func NewReviewUsecase(
	reviewRepo domain.ReviewRepository,
	productRepo domain.ProductRepository,
) *ReviewUsecase {
	return &ReviewUsecase{
		reviewRepo:  reviewRepo,
		productRepo: productRepo,
	}
}

func (u *ReviewUsecase) ListReviews(ctx context.Context, productID string, includeUnapproved bool) ([]domain.Review, error) {
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

func (u *ReviewUsecase) CreateReview(ctx context.Context, productID, reviewerName string, rating int, body *string) (domain.Review, error) {
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

func (u *ReviewUsecase) GetReviewByID(ctx context.Context, reviewID string) (domain.Review, bool, error) {
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

func (u *ReviewUsecase) UpdateReviewApproval(ctx context.Context, reviewID string, isApproved bool) error {
	return u.reviewRepo.UpdateApproval(ctx, reviewID, isApproved)
}

func (u *ReviewUsecase) DeleteReview(ctx context.Context, reviewID string) error {
	return u.reviewRepo.Delete(ctx, reviewID)
}

func (u *ReviewUsecase) ListAllReviews(ctx context.Context, approved *bool) ([]domain.Review, error) {
	reviews, err := u.reviewRepo.ListAll(ctx, approved)
	if err != nil {
		return nil, err
	}
	sort.Slice(reviews, func(i, j int) bool { return reviews[i].CreatedAt.After(reviews[j].CreatedAt) })
	return reviews, nil
}
