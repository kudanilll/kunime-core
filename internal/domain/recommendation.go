package domain

import "time"

type RecommendationItem struct {
	AnimeID   string    `json:"anime_id"`
	Title     string    `json:"title"`
	Rating    float64   `json:"rating"`
	Image     string    `json:"image"`
	Source    string    `json:"source"`
	UpdatedAt time.Time `json:"updated_at"`
}
