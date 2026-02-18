package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type GenreRepository struct {
	client *firestore.Client
}

func NewGenreRepository(client *firestore.Client) *GenreRepository {
	return &GenreRepository{client: client}
}

func (r *GenreRepository) GetAll(ctx context.Context) ([]domain.GenreAsset, error) {
	iter := r.client.
		Collection("genre_assets").
		Documents(ctx)

	defer iter.Stop()

	var result []domain.GenreAsset

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var g domain.GenreAsset
		if err := doc.DataTo(&g); err != nil {
			return nil, err
		}

		g.Slug = doc.Ref.ID
		result = append(result, g)
	}

	return result, nil
}
