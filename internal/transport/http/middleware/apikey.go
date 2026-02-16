package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kudanilll/kunime-core/internal/config"
)

func APIKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing api key",
			})
		}

		if apiKey != config.APIKey() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid api key",
			})
		}

		return c.Next()
	}
}
