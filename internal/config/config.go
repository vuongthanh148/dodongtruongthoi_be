package config

import (
"os"

"github.com/joho/godotenv"
)

type Config struct {
AppName     string
AppEnv      string
Port        string
DatabaseURL string
}

func Load() (*Config, error) {
_ = godotenv.Load()

cfg := &Config{
AppName:     getEnv("APP_NAME", "app"),
AppEnv:      getEnv("APP_ENV", "development"),
Port:        getEnv("PORT", "8080"),
DatabaseURL: os.Getenv("DATABASE_URL"),
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
