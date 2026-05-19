package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type ImageUsecase struct {
	imageRepo     domain.ImageRepository
	imageUploader ImageUploader
}

func NewImageUsecase(imageRepo domain.ImageRepository, imageUploader ImageUploader) *ImageUsecase {
	return &ImageUsecase{imageRepo: imageRepo, imageUploader: imageUploader}
}

func (u *ImageUsecase) ListImages(ctx context.Context) ([]domain.Image, error) {
	return u.imageRepo.List(ctx)
}

func (u *ImageUsecase) GetImage(ctx context.Context, id string) (domain.Image, bool, error) {
	return u.imageRepo.Get(ctx, id)
}

// UploadImage uploads a file to Cloudinary and persists it to the image library.
func (u *ImageUsecase) UploadImage(ctx context.Context, file interface{}, filename string, name string) (domain.Image, error) {
	if u.imageUploader == nil {
		return domain.Image{}, errors.New("image uploads disabled; cloudinary not configured")
	}

	url, err := u.imageUploader.UploadImage(ctx, file, filename, "library")
	if err != nil {
		return domain.Image{}, err
	}

	// Extract cloudinary public_id from URL (after /upload/ and before extension)
	publicID := extractCloudinaryPublicID(url)

	img := domain.Image{
		ID:                 uuid.NewString(),
		Name:               name,
		URL:                url,
		CloudinaryPublicID: publicID,
		CreatedAt:          time.Now(),
	}
	return u.imageRepo.Create(ctx, img)
}

func (u *ImageUsecase) RenameImage(ctx context.Context, id string, name string) (domain.Image, error) {
	return u.imageRepo.Update(ctx, id, name)
}

// DeleteImage deletes from Cloudinary + cascades product_images via FK.
func (u *ImageUsecase) DeleteImage(ctx context.Context, id string) error {
	img, ok, err := u.imageRepo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("image not found")
	}

	if u.imageUploader != nil && img.CloudinaryPublicID != "" {
		_ = u.imageUploader.DeleteImage(ctx, img.CloudinaryPublicID)
	}

	return u.imageRepo.Delete(ctx, id)
}

func extractCloudinaryPublicID(url string) string {
	const marker = "/upload/"
	idx := -1
	for i := 0; i < len(url)-len(marker); i++ {
		if url[i:i+len(marker)] == marker {
			idx = i + len(marker)
			break
		}
	}
	if idx < 0 {
		return url
	}
	// Skip version segment like "v1234567890/"
	rest := url[idx:]
	if len(rest) > 1 && rest[0] == 'v' {
		for i := 1; i < len(rest); i++ {
			if rest[i] == '/' {
				rest = rest[i+1:]
				break
			}
		}
	}
	// Strip extension
	for i := len(rest) - 1; i >= 0; i-- {
		if rest[i] == '.' {
			return rest[:i]
		}
		if rest[i] == '/' {
			break
		}
	}
	return rest
}
