package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUploader(cloudName, apiKey, apiSecret string) (*CloudinaryUploader, error) {
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return &CloudinaryUploader{}, nil
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	return &CloudinaryUploader{cld: cld}, nil
}

func (c *CloudinaryUploader) UploadImage(ctx context.Context, file interface{}, filename string, folder string) (string, error) {
	if c.cld == nil {
		return "", fmt.Errorf("cloudinary not configured")
	}

	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	result, err := c.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload error: %w", err)
	}

	if result == nil || result.SecureURL == "" {
		return "", fmt.Errorf("cloudinary returned empty URL")
	}

	return result.SecureURL, nil
}

func (c *CloudinaryUploader) UploadImageFromReader(ctx context.Context, reader io.Reader, filename string, folder string) (string, error) {
	if c.cld == nil {
		return "", fmt.Errorf("cloudinary not configured")
	}

	uploadParams := uploader.UploadParams{
		Folder: folder,
	}

	result, err := c.cld.Upload.Upload(ctx, reader, uploadParams)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

func (c *CloudinaryUploader) DeleteImage(ctx context.Context, publicID string) error {
	if c.cld == nil {
		return fmt.Errorf("cloudinary not configured")
	}

	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}
