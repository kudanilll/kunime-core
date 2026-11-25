package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/kudanilll/kunime-core/internal/config"

	"github.com/kudanilll/kunime-core/internal/transport/http/handler"
	"github.com/kudanilll/kunime-core/internal/transport/http/middleware"
)

type Handlers struct {
	Watch          *handler.WatchHandler
	Popular        *handler.PopularHandler
	Recommendation *handler.RecommendationHandler
	App            *handler.AppHandler
	Banner         *handler.BannerHandler
	Genre          *handler.GenreHandler
}

func RegisterRoutes(app *fiber.App, h *Handlers) {
	// Health
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Admin routes
	admin := app.Group("/api/v1/admin",
		middleware.RequireAPIKey(config.AdminAPIKey),
	)

	admin.Post("/app/version", h.App.UpdateVersion)
	admin.Post("/recommendations/refresh", h.Recommendation.Refresh)

	// Client routes
	api := app.Group("/api/v1",
		middleware.RequireAPIKey(config.ClientAPIKey),
	)

	api.Post("/events/watch", h.Watch.Record)
	api.Get("/anime/popular", h.Popular.Get)
	api.Get("/recommendations", h.Recommendation.Get)
	api.Get("/banners", h.Banner.Get)
	api.Get("/genres", h.Genre.Get)
	api.Get("/app/version", h.App.GetVersion)
}
