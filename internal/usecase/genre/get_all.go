package genre

import (
	"context"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.GenreAsset, error)
}

type GetAllUsecase struct {
	repo Repository
}

func NewGetAllUsecase(repo Repository) *GetAllUsecase {
	return &GetAllUsecase{repo: repo}
}

func (u *GetAllUsecase) Execute(ctx context.Context) ([]domain.GenreAsset, error) {
	return u.repo.GetAll(ctx)
}
