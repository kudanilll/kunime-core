package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/kudanilll/kunime-core/internal/config"
	"github.com/kudanilll/kunime-core/internal/domain"
	"github.com/kudanilll/kunime-core/internal/transport/http/handler"
	"github.com/kudanilll/kunime-core/internal/usecase/watch"
)

type watchStore struct{ event *domain.WatchEvent }

func (s *watchStore) Save(_ context.Context, event *domain.WatchEvent) error {
	s.event = event
	return nil
}

func (*watchStore) Increment(context.Context, string) error { return nil }

func TestWatchRouteContract(t *testing.T) {
	t.Setenv("CLIENT_API_KEY", "client-test-key")
	t.Setenv("ADMIN_API_KEY", "admin-test-key")
	if _, err := config.Load(); err != nil {
		t.Fatal(err)
	}
	app := fiber.New()
	store := &watchStore{}
	RegisterRoutes(app, &Handlers{Watch: handler.NewWatchHandler(watch.NewRecordWatchUsecase(store, store))})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/events/watch", strings.NewReader(`{"anime_id":"sample","episode":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Device-Id", "test-device")
	req.Header.Set("X-API-Key", "client-test-key")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result["status"] != "recorded" {
		t.Fatalf("response = %v, error = %v", result, err)
	}
	if store.event == nil || store.event.AnimeID != "sample" || store.event.Episode != 1 || store.event.DeviceID != "test-device" {
		t.Fatalf("unexpected watch event: %+v", store.event)
	}
}
