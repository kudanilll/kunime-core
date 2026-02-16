package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kudanilll/kunime-core/internal/transport/http/handler"
	"github.com/kudanilll/kunime-core/internal/transport/http/middleware"
)

type Handlers struct {
	Watch   *handler.WatchHandler
	Popular *handler.PopularHandler
}

func RegisterRoutes(app *fiber.App, h *Handlers) {
	// Public routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Protected API routes
	api := app.Group("/api/v1", middleware.APIKey())

	api.Post("/events/watch", h.Watch.Record)

	api.Get("/anime/popular", h.Popular.Get)

	api.Get("/recommendations", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "recommendation endpoint (not implemented yet)",
		})
	})

	api.Get("/banners", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "banner endpoint (not implemented yet)",
		})
	})

	api.Get("/genres", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "genre endpoint (not implemented yet)",
		})
	})

	api.Get("/app/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "app version endpoint (not implemented yet)",
		})
	})
}
