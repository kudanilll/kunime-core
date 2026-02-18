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
	const (
		TargetRecommendation = 5
		MaxPage              = 10
	)

	var candidates []firestore.RecommendationItem

	for page := 1; page <= MaxPage; page++ {
		completed, err := u.api.FetchCompleted(page)
		if err != nil {
			return err
		}

		for _, anime := range completed {
			// Filter rating
			if anime.Rating <= 0 || anime.Rating > 10 {
				continue
			}

			candidates = append(candidates, firestore.RecommendationItem{
				AnimeID:   anime.AnimeID,
				Title:     anime.Title,
				Rating:    anime.Rating,
				Image:     anime.Image,
				Source:    "completed_top_rated",
				UpdatedAt: time.Now(),
			})
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Sort by rating
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Rating > candidates[j].Rating
	})

	if len(candidates) > TargetRecommendation {
		candidates = candidates[:TargetRecommendation]
	}

	return u.repo.ReplaceAll(ctx, candidates)
}
