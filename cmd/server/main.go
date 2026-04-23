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
if cfg.AppEnv == "development" {
log.Printf("warning: database unavailable, continue without DB in development: %v", err)
} else {
log.Fatalf("failed to connect database: %v", err)
}
}
if dbPool != nil {
defer database.Close(dbPool)
}

healthUsecase := usecase.NewHealthUsecase(cfg.AppName)
healthHandler := handler.NewHealthHandler(healthUsecase)
router := httpdelivery.NewRouter(healthHandler)

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
