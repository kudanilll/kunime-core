package watch

import (
	"context"
	"errors"
	"time"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, event *domain.WatchEvent) error
}

type RecordWatchUsecase struct {
	repo Repository
}

func NewRecordWatchUsecase(repo Repository) *RecordWatchUsecase {
	return &RecordWatchUsecase{repo: repo}
}

func (u *RecordWatchUsecase) Execute(
	ctx context.Context,
	animeID string,
	episode int,
	deviceID string,
) error {
	if animeID == "" {
		return errors.New("anime_id is required")
	}
	if episode <= 0 {
		return errors.New("episode must be greater than 0")
	}
	if deviceID == "" {
		return errors.New("device_id is required")
	}

	now := time.Now()

	event := &domain.WatchEvent{
		AnimeID:   animeID,
		Episode:   episode,
		DeviceID:  deviceID,
		WatchedAt: now,
		ExpiresAt: now.Add(14 * 24 * time.Hour), // TTL 14 hari
	}

	return u.repo.Save(ctx, event)
}
