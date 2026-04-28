package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type CategoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *CategoryUsecase) ListCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx, false)
}

func (u *CategoryUsecase) GetCategory(ctx context.Context, id string) (domain.Category, bool, error) {
	return u.categoryRepo.Get(ctx, id, false)
}

func (u *CategoryUsecase) ListAllCategories(ctx context.Context) ([]domain.Category, error) {
	return u.categoryRepo.List(ctx, true)
}

func (u *CategoryUsecase) CreateCategory(ctx context.Context, c domain.Category) (domain.Category, error) {
	if c.ID == "" || c.Name == "" {
		return domain.Category{}, errors.New("id and name are required")
	}

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.IsActive = true

	return u.categoryRepo.Create(ctx, c)
}

func (u *CategoryUsecase) UpdateCategory(ctx context.Context, id string, updates domain.Category) (domain.Category, error) {
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

func (u *CategoryUsecase) DeleteCategory(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("category id is required")
	}
	return u.categoryRepo.Delete(ctx, id)
}
