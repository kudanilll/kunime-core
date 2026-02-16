package firestore

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type WatchRepository struct {
	client *firestore.Client
}

func NewWatchRepository(client *firestore.Client) *WatchRepository {
	return &WatchRepository{client: client}
}

func (r *WatchRepository) Save(ctx context.Context, event *domain.WatchEvent) error {
	_, _, err := r.client.
		Collection("watch_event").
		Add(ctx, map[string]interface{}{
			"anime_id":   event.AnimeID,
			"episode":    event.Episode,
			"device_id":  event.DeviceID,
			"watched_at": event.WatchedAt,
			"expires_at": event.ExpiresAt,
		})

	return err
}
