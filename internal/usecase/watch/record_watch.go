package watch

import (
	"context"
	"errors"
	"time"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type WatchRepository interface {
	Save(ctx context.Context, event *domain.WatchEvent) error
}

type PopularRepository interface {
	Increment(ctx context.Context, animeID string) error
}

type RecordWatchUsecase struct {
	watchRepo   WatchRepository
	popularRepo PopularRepository
}

func NewRecordWatchUsecase(
	watchRepo WatchRepository,
	popularRepo PopularRepository,
) *RecordWatchUsecase {
	return &RecordWatchUsecase{
		watchRepo:   watchRepo,
		popularRepo: popularRepo,
	}
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
		ExpiresAt: now.Add(14 * 24 * time.Hour),
	}

	// Save watch event
	if err := u.watchRepo.Save(ctx, event); err != nil {
		return err
	}

	// Update popular count
	if err := u.popularRepo.Increment(ctx, animeID); err != nil {
		return err
	}

	return nil
}
