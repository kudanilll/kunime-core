package recommendation

import (
	"context"
	"sort"
	"time"

	"github.com/kudanilll/kunime-core/internal/infra/firestore"
)

type APIClient interface {
	FetchCompleted(page int) ([]CompletedAnime, error)
}

type CompletedAnime struct {
	AnimeID string
	Title   string
	Rating  float64
	Image   string
}

type Repository interface {
	ReplaceAll(ctx context.Context, items []firestore.RecommendationItem) error
}

type BuildRecommendationUsecase struct {
	api  APIClient
	repo Repository
}

func NewBuildRecommendationUsecase(api APIClient, repo Repository) *BuildRecommendationUsecase {
	return &BuildRecommendationUsecase{api: api, repo: repo}
}

func (u *BuildRecommendationUsecase) Execute(ctx context.Context) error {
	var all []CompletedAnime

	for page := 1; page <= 10; page++ {
		items, err := u.api.FetchCompleted(page)
		if err != nil {
			return err
		}
		all = append(all, items...)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Rating > all[j].Rating
	})

	top := all
	if len(top) > 5 {
		top = top[:5]
	}

	var result []firestore.RecommendationItem
	now := time.Now()

	for _, a := range top {
		result = append(result, firestore.RecommendationItem{
			AnimeID:   a.AnimeID,
			Title:     a.Title,
			Rating:    a.Rating,
			Image:     a.Image,
			Source:    "completed_top_rated",
			UpdatedAt: now,
		})
	}

	return u.repo.ReplaceAll(ctx, result)
}
