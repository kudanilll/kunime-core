package domain

import "time"

type WatchEvent struct {
	AnimeID   string    `json:"anime_id"`
	Episode   int       `json:"episode"`
	DeviceID  string    `json:"device_id"`
	WatchedAt time.Time `json:"watched_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
