package domain

import "time"

type Banner struct {
	ID          string    `json:"id"`
	Title       string    `json:"title" firestore:"title"`
	ImageURL    string    `json:"image_url" firestore:"image_url"`
	RedirectURL string    `json:"redirect_url" firestore:"redirect_url"`
	IsActive    bool      `json:"-" firestore:"is_active"`
	StartAt     time.Time `json:"-" firestore:"start_at"`
	EndAt       time.Time `json:"-" firestore:"end_at"`
	Priority    int       `json:"-" firestore:"priority"`
	CreatedAt   time.Time `json:"-" firestore:"created_at"`
	UpdatedAt   time.Time `json:"-" firestore:"updated_at"`
}
