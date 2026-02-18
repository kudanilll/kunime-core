package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kudanilll/kunime-core/internal/usecase/app"
)

type AppHandler struct {
	uc *app.GetVersionUsecase
}

func NewAppHandler(uc *app.GetVersionUsecase) *AppHandler {
	return &AppHandler{uc: uc}
}

func (h *AppHandler) GetVersion(c *fiber.Ctx) error {
	result, err := h.uc.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get app version",
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}
