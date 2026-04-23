package domain

type HealthStatus struct {
Success bool   `json:"success"`
Message string `json:"message"`
App     string `json:"app"`
}
