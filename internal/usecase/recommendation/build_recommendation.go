package recommendation

import (
	"context"
	"sort"
	"time"

	"github.com/kudanilll/kunime-core/internal/infra/firestore"
	"go.uber.org/zap"
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
	api    APIClient
	repo   Repository
	logger *zap.Logger
}

func NewBuildRecommendationUsecase(api APIClient, repo Repository, logger *zap.Logger) *BuildRecommendationUsecase {
	return &BuildRecommendationUsecase{api: api, repo: repo, logger: logger}
}
func (u *BuildRecommendationUsecase) Execute(ctx context.Context) error {
	const (
		TargetRecommendation = 5
		MaxPage              = 10
	)

	u.logger.Info("recommendation build started", zap.Int("target_recommendation", TargetRecommendation), zap.Int("max_page", MaxPage))

	var candidates []firestore.RecommendationItem

	for page := 1; page <= MaxPage; page++ {
		completed, err := u.api.FetchCompleted(page)
		if err != nil {
			u.logger.Error("recommendation build failed while fetching completed anime", zap.Int("page", page), zap.Error(err))
			return err
		}

		u.logger.Info("recommendation build fetched page", zap.Int("page", page), zap.Int("items", len(completed)))

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
		u.logger.Warn("recommendation build finished with no candidates")
		return nil
	}

	// Sort by rating
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Rating > candidates[j].Rating
	})

	if len(candidates) > TargetRecommendation {
		candidates = candidates[:TargetRecommendation]
	}

	if err := u.repo.ReplaceAll(ctx, candidates); err != nil {
		u.logger.Error("recommendation build failed to replace firestore data", zap.Int("final_count", len(candidates)), zap.Error(err))
		return err
	}

	u.logger.Info("recommendation build finished", zap.Int("final_count", len(candidates)))
	return nil
}
