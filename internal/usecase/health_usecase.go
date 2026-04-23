package usecase

import (
"context"

"github.com/vuongthanh148/dodongtruongthoi_be/internal/domain"
)

type HealthUsecase struct {
appName string
}

func NewHealthUsecase(appName string) *HealthUsecase {
return &HealthUsecase{appName: appName}
}

func (u *HealthUsecase) GetStatus(ctx context.Context) domain.HealthStatus {
_ = ctx

return domain.HealthStatus{
Success: true,
Message: "server is running",
App:     u.appName,
}
}
