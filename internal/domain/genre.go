package domain

type GenreAsset struct {
	Slug     string `json:"slug"`
	ImageURL string `json:"image_url" firestore:"image_url"`
}
