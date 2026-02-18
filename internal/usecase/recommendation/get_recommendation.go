package recommendation

import (
	"context"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type ReadRepository interface {
	GetAll(ctx context.Context) ([]domain.RecommendationItem, error)
}

type GetRecommendationUsecase struct {
	repo ReadRepository
}

func NewGetRecommendationUsecase(repo ReadRepository) *GetRecommendationUsecase {
	return &GetRecommendationUsecase{repo: repo}
}

func (u *GetRecommendationUsecase) Execute(
	ctx context.Context,
) ([]domain.RecommendationItem, error) {
	return u.repo.GetAll(ctx)
}
