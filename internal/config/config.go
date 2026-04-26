package config

import (
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

cfg := &Config{
AppName:              getEnv("APP_NAME", "app"),
AppEnv:              getEnv("APP_ENV", "development"),
Port:                getEnv("PORT", "8080"),
DatabaseURL:         os.Getenv("DATABASE_URL"),
		JWTSecret:            getEnv("JWT_SECRET", "dev-secret"),
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
