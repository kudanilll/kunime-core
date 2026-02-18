package banner

import (
	"context"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type Repository interface {
	GetActive(ctx context.Context) ([]domain.Banner, error)
}

type GetActiveUsecase struct {
	repo Repository
}

func NewGetActiveUsecase(repo Repository) *GetActiveUsecase {
	return &GetActiveUsecase{repo: repo}
}

func (u *GetActiveUsecase) Execute(ctx context.Context) ([]domain.Banner, error) {
	return u.repo.GetActive(ctx)
}
