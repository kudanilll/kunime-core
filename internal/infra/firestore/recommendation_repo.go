package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kudanilll/kunime-core/internal/domain"
	"google.golang.org/api/iterator"
)

type RecommendationItem struct {
	AnimeID   string    `firestore:"anime_id"`
	Title     string    `firestore:"title"`
	Rating    float64   `firestore:"rating"`
	Image     string    `firestore:"image"`
	Source    string    `firestore:"source"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

type RecommendationRepository struct {
	client *firestore.Client
}

func NewRecommendationRepository(client *firestore.Client) *RecommendationRepository {
	return &RecommendationRepository{client: client}
}

func (r *RecommendationRepository) ReplaceAll(
	ctx context.Context,
	items []RecommendationItem,
) error {
	batch := r.client.Batch()

	col := r.client.Collection("anime_recommendation")

	docs, err := col.Documents(ctx).GetAll()
	if err != nil {
		return err
	}
	for _, d := range docs {
		batch.Delete(d.Ref)
	}

	for _, item := range items {
		ref := col.Doc(item.AnimeID)
		batch.Set(ref, item)
	}

	_, err = batch.Commit(ctx)
	return err
}

func (r *RecommendationRepository) GetAll(
	ctx context.Context,
) ([]domain.RecommendationItem, error) {
	iter := r.client.
		Collection("anime_recommendation").
		OrderBy("rating", firestore.Desc).
		Documents(ctx)

	defer iter.Stop()

	var result []domain.RecommendationItem
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var raw RecommendationItem
		if err := doc.DataTo(&raw); err != nil {
			return nil, err
		}

		result = append(result, domain.RecommendationItem{
			AnimeID:   raw.AnimeID,
			Title:     raw.Title,
			Rating:    raw.Rating,
			Image:     raw.Image,
			Source:    raw.Source,
			UpdatedAt: raw.UpdatedAt,
		})
	}

	return result, nil
}
