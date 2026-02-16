package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

type PopularRepository struct {
	client *firestore.Client
}

func NewPopularRepository(client *firestore.Client) *PopularRepository {
	return &PopularRepository{client: client}
}

func (r *PopularRepository) Increment(ctx context.Context, animeID string) error {
	doc := r.client.Collection("anime_popular").Doc(animeID)

	_, err := doc.Set(ctx, map[string]interface{}{
		"watch_count":  firestore.Increment(1),
		"last_updated": time.Now(),
	}, firestore.MergeAll)

	return err
}
