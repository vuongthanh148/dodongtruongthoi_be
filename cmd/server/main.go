package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vuongthanh148/dodongtruongthoi_be/internal/config"
	httpdelivery "github.com/vuongthanh148/dodongtruongthoi_be/internal/delivery/http"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/delivery/http/handler"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/infrastructure/database"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/infrastructure/storage"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/repository/postgres"
	"github.com/vuongthanh148/dodongtruongthoi_be/internal/usecase"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	dbPool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if dbPool == nil {
		log.Fatalf("database connection required; in-memory fallback is not supported")
	}
	defer database.Close(dbPool)
	log.Println("✓ PostgreSQL connected")

	uploader, err := storage.NewCloudinaryUploader(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		log.Printf("warning: cloudinary initialization failed: %v", err)
	}
	if uploader != nil && cfg.CloudinaryCloudName != "" {
		log.Println("✓ Cloudinary configured")
	} else {
		log.Println("ℹ Cloudinary not configured (image uploads disabled)")
	}

	healthUsecase := usecase.NewHealthUsecase(cfg.AppName)
	healthHandler := handler.NewHealthHandler(healthUsecase)

	platformUsecase, err := usecase.NewPlatformUsecase(usecase.PlatformUsecaseConfig{
		JWTSecret:        cfg.JWTSecret,
		CategoryRepo:     postgres.NewCategoryRepository(dbPool),
		ProductRepo:      postgres.NewProductRepository(dbPool),
		ProductImageRepo: postgres.NewProductImageRepository(dbPool),
		ProductSizeRepo:  postgres.NewProductSizeRepository(dbPool),
		CampaignRepo:     postgres.NewCampaignRepository(dbPool),
		ReviewRepo:       postgres.NewReviewRepository(dbPool),
		OrderRepo:        postgres.NewOrderRepository(dbPool),
		WishlistRepo:     postgres.NewWishlistRepository(dbPool),
		BannerRepo:       postgres.NewBannerRepository(dbPool),
		ContactRepo:      postgres.NewContactLinkRepository(dbPool),
		AdminUserRepo:    postgres.NewAdminUserRepository(dbPool),
		SettingsRepo:     postgres.NewSiteSettingsRepository(dbPool),
		CustomerPhotoRepo: postgres.NewCustomerPhotoRepository(dbPool),
		ImageUploader:    uploader,
	})
	if err != nil {
		log.Fatalf("failed to create platform usecase: %v", err)
	}
	log.Println("✓ Using PostgreSQL repositories")

	publicHandler := handler.NewPublicHandler(platformUsecase)
	adminHandler := handler.NewAdminHandler(platformUsecase)
	router := httpdelivery.NewRouter(cfg, healthHandler, publicHandler, adminHandler, platformUsecase)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("%s started on %s in %s mode", cfg.AppName, server.Addr, cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped")
}
