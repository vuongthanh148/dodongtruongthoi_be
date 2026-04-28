package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
AppName              string
AppEnv               string
Port                 string
DatabaseURL          string
	JWTSecret            string
	CloudinaryCloudName  string
	CloudinaryAPIKey     string
	CloudinaryAPISecret  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" || len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be set and at least 32 characters")
	}

	cfg := &Config{
		AppName:              getEnv("APP_NAME", "app"),
		AppEnv:              getEnv("APP_ENV", "development"),
		Port:                getEnv("PORT", "8080"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		JWTSecret:            jwtSecret,
		CloudinaryCloudName:  os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:     os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret:  os.Getenv("CLOUDINARY_API_SECRET"),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
value := os.Getenv(key)
if value == "" {
return fallback
}
return value
}
