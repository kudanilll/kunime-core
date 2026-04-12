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
	"go.uber.org/zap"
)

type KunimeAPIClient struct {
	baseURL       string
	apiKey        string
	client        *http.Client
	logger        *zap.Logger
	mu            sync.RWMutex
	scrapeBaseURL string
}

func NewKunimeAPIClient(baseURL, apiKey string, logger *zap.Logger) *KunimeAPIClient {
	return &KunimeAPIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		logger:  logger,
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
		c.logger.Error("recommendation fetch failed to get scrape base url", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf("%s/completed-anime/%d", c.baseURL, page)
	c.logger.Info("recommendation fetch completed anime request", zap.Int("page", page), zap.String("url", url))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("recommendation fetch failed to create request", zap.Int("page", page), zap.Error(err))
		return nil, err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("recommendation fetch request failed", zap.Int("page", page), zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("kunime-api returned %d", resp.StatusCode)
		c.logger.Error("recommendation fetch returned non-200", zap.Int("page", page), zap.Int("status_code", resp.StatusCode), zap.Error(err))
		return nil, err
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
		c.logger.Error("recommendation fetch failed to decode response", zap.Int("page", page), zap.Error(err))
		return nil, err
	}

	c.logger.Info("recommendation fetch completed anime success", zap.Int("page", page), zap.Int("items", len(raw.Data)))

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
	c.logger.Info("recommendation fetch scrape base url request", zap.String("url", url))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.logger.Error("recommendation fetch scrape base url failed to create request", zap.Error(err))
		return "", err
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("recommendation fetch scrape base url request failed", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("scrape-base-url returned %d", resp.StatusCode)
		c.logger.Error("recommendation fetch scrape base url returned non-200", zap.Int("status_code", resp.StatusCode), zap.Error(err))
		return "", err
	}

	var raw struct {
		ScrapeBaseURL string `json:"scrape_base_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		c.logger.Error("recommendation fetch scrape base url failed to decode response", zap.Error(err))
		return "", err
	}

	result := strings.TrimRight(raw.ScrapeBaseURL, "/")
	c.logger.Info("recommendation fetch scrape base url success", zap.String("scrape_base_url", result))

	return result, nil
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
