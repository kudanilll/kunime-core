package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kudanilll/kunime-core/internal/usecase/recommendation"
)

type KunimeAPIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewKunimeAPIClient(baseURL, apiKey string) *KunimeAPIClient {
	return &KunimeAPIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *KunimeAPIClient) FetchCompleted(
	page int,
) ([]recommendation.CompletedAnime, error) {
	url := fmt.Sprintf("%s/completed-anime/%d", c.baseURL, page)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("kunime-api returned %d", resp.StatusCode)
	}

	var raw struct {
		Data []struct {
			Title    string  `json:"title"`
			Score    float64 `json:"score"`
			Image    string  `json:"image"`
			Endpoint string  `json:"endpoint"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var result []recommendation.CompletedAnime
	for _, a := range raw.Data {
		result = append(result, recommendation.CompletedAnime{
			AnimeID: extractAnimeSlug(a.Endpoint),
			Title:   a.Title,
			Rating:  a.Score,
			Image:   a.Image,
		})
	}

	return result, nil
}

func extractAnimeSlug(endpoint string) string {
	u := strings.TrimSuffix(endpoint, "/")
	parts := strings.Split(u, "/")
	return parts[len(parts)-1]
}
