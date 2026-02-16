package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/kudanilll/kunime-core/internal/usecase/popular"
	"google.golang.org/api/iterator"
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

func (r *PopularRepository) GetTop(
	ctx context.Context,
	limit int,
) ([]popular.Item, error) {
	if limit <= 0 {
		limit = 10
	}

	iter := r.client.
		Collection("anime_popular").
		OrderBy("watch_count", firestore.Desc).
		Limit(limit).
		Documents(ctx)

	defer iter.Stop()

	var result []popular.Item
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var data struct {
			WatchCount int64 `firestore:"watch_count"`
		}

		if err := doc.DataTo(&data); err != nil {
			return nil, err
		}

		result = append(result, popular.Item{
			AnimeID:    doc.Ref.ID,
			WatchCount: data.WatchCount,
		})
	}

	return result, nil
}
