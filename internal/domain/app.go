package domain

type AppVersion struct {
	LatestVersion       string `json:"latest_version" firestore:"latest_version"`
	MinSupportedVersion string `json:"min_supported_version" firestore:"min_supported_version"`
	ForceUpdate         bool   `json:"force_update" firestore:"force_update"`
	UpdateURL           string `json:"update_url" firestore:"update_url"`
}
