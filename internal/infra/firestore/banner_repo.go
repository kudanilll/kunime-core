package firestore

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type BannerRepository struct {
	client *firestore.Client
}

func NewBannerRepository(client *firestore.Client) *BannerRepository {
	return &BannerRepository{client: client}
}

func (r *BannerRepository) GetActive(ctx context.Context) ([]domain.Banner, error) {
	now := time.Now()

	iter := r.client.
		Collection("ads_banners").
		Where("is_active", "==", true).
		Where("start_at", "<=", now).
		OrderBy("start_at", firestore.Desc).
		Documents(ctx)

	defer iter.Stop()

	var result []domain.Banner

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var b domain.Banner
		if err := doc.DataTo(&b); err != nil {
			return nil, err
		}

		// manual end_at filter
		if now.After(b.EndAt) {
			continue
		}

		b.ID = doc.Ref.ID
		result = append(result, b)
	}

	return result, nil
}
