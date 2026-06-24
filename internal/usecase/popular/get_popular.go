package popular

import (
	"context"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type Repository interface {
	GetTop(ctx context.Context, limit int) ([]domain.AnimePopular, error)
}

type GetPopularUsecase struct {
	repo Repository
}

func NewGetPopularUsecase(repo Repository) *GetPopularUsecase {
	return &GetPopularUsecase{repo: repo}
}

func (u *GetPopularUsecase) Execute(
	ctx context.Context,
	limit int,
) ([]domain.AnimePopular, error) {
	return u.repo.GetTop(ctx, limit)
}
