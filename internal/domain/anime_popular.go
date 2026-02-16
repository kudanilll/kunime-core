package domain

import "time"

type AnimePopular struct {
	AnimeID     string
	WatchCount  int64
	LastUpdated time.Time
}
