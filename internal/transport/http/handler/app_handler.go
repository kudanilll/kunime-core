package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kudanilll/kunime-core/internal/domain"
	"github.com/kudanilll/kunime-core/internal/usecase/app"
)

type AppHandler struct {
	getUC    *app.GetVersionUsecase
	updateUC *app.UpdateVersionUsecase
}

func NewAppHandler(
	getUC *app.GetVersionUsecase,
	updateUC *app.UpdateVersionUsecase,
) *AppHandler {
	return &AppHandler{
		getUC:    getUC,
		updateUC: updateUC,
	}
}

func (h *AppHandler) GetVersion(c fiber.Ctx) error {
	result, err := h.getUC.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get app version",
		})
	}

	return c.JSON(fiber.Map{
		"data": result,
	})
}

func (h *AppHandler) UpdateVersion(c fiber.Ctx) error {
	var payload domain.AppVersion

	if err := c.Bind().Body(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if payload.LatestVersion == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "latest_version is required",
		})
	}

	if err := h.updateUC.Execute(c.Context(), &payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update version",
		})
	}

	return c.JSON(fiber.Map{
		"message": "app version updated",
	})
}
