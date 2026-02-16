package popular

import "context"

type Repository interface {
	GetTop(ctx context.Context, limit int) ([]Item, error)
}

type Item struct {
	AnimeID    string `json:"anime_id"`
	WatchCount int64  `json:"watch_count"`
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
) ([]Item, error) {
	return u.repo.GetTop(ctx, limit)
}
