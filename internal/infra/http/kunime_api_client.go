package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/kudanilll/kunime-core/internal/usecase/recommendation"
)

type KunimeAPIClient struct {
	baseURL       string
	apiKey        string
	client        *http.Client
	mu            sync.RWMutex
	scrapeBaseURL string
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
	scrapeBaseURL, err := c.getScrapeBaseURL()
	if err != nil {
		return nil, err
	}

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
			Image:   buildScrapeImageURL(scrapeBaseURL, a.Image),
		})
	}

	return result, nil
}

func (c *KunimeAPIClient) getScrapeBaseURL() (string, error) {
	c.mu.RLock()
	if c.scrapeBaseURL != "" {
		defer c.mu.RUnlock()
		return c.scrapeBaseURL, nil
	}
	c.mu.RUnlock()

	scrapeBaseURL, err := c.fetchScrapeBaseURL()
	if err != nil {
		return "", err
	}

	c.mu.Lock()
	c.scrapeBaseURL = scrapeBaseURL
	c.mu.Unlock()

	return scrapeBaseURL, nil
}

func (c *KunimeAPIClient) fetchScrapeBaseURL() (string, error) {
	url := fmt.Sprintf("%s/scrape-base-url", c.baseURL)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scrape-base-url returned %d", resp.StatusCode)
	}

	var raw struct {
		ScrapeBaseURL string `json:"scrape_base_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return "", err
	}

	return strings.TrimRight(raw.ScrapeBaseURL, "/"), nil
}

func buildScrapeImageURL(scrapeBaseURL, imageURL string) string {
	if scrapeBaseURL == "" || imageURL == "" {
		return imageURL
	}

	imagePath := imageURL
	if parsed, err := url.Parse(imageURL); err == nil && parsed.Path != "" {
		imagePath = parsed.Path
	}

	if !strings.HasPrefix(imagePath, "/") {
		imagePath = "/" + strings.TrimLeft(imagePath, "/")
	}

	return strings.TrimRight(scrapeBaseURL, "/") + imagePath
}

func extractAnimeSlug(endpoint string) string {
	u := strings.TrimSuffix(endpoint, "/")
	parts := strings.Split(u, "/")
	return parts[len(parts)-1]
}
