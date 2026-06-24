package http

import (
	"context"
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
	ctx context.Context,
	page int,
) ([]recommendation.CompletedAnime, error) {
	scrapeBaseURL, err := c.getScrapeBaseURL(ctx)
	if err != nil {
		c.logger.Error("recommendation fetch failed to get scrape base url", zap.Error(err))
		return nil, err
	}

	url := fmt.Sprintf("%s/completed-anime/%d", c.baseURL, page)
	c.logger.Info("recommendation fetch completed anime request", zap.Int("page", page), zap.String("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

func (c *KunimeAPIClient) getScrapeBaseURL(ctx context.Context) (string, error) {
	c.mu.RLock()
	val := c.scrapeBaseURL
	c.mu.RUnlock()
	if val != "" {
		return val, nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check again after acquiring write lock
	if c.scrapeBaseURL != "" {
		return c.scrapeBaseURL, nil
	}

	scrapeBaseURL, err := c.fetchScrapeBaseURL(ctx)
	if err != nil {
		return "", err
	}

	c.scrapeBaseURL = scrapeBaseURL

	return scrapeBaseURL, nil
}

func (c *KunimeAPIClient) fetchScrapeBaseURL(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/scrape-base-url", c.baseURL)
	c.logger.Info("recommendation fetch scrape base url request", zap.String("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
