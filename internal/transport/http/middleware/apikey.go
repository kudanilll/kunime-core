package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func RequireAPIKey(getExpected func() string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-Key")
		expected := getExpected()

		if key == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing api key",
			})
		}

		if key != expected {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid api key",
			})
		}

		return c.Next()
	}
}
