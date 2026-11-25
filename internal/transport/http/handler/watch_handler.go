package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kudanilll/kunime-core/internal/usecase/watch"
)

type WatchHandler struct {
	usecase *watch.RecordWatchUsecase
}

func NewWatchHandler(usecase *watch.RecordWatchUsecase) *WatchHandler {
	return &WatchHandler{usecase: usecase}
}

type watchRequest struct {
	AnimeID string `json:"anime_id"`
	Episode int    `json:"episode"`
}

func (h *WatchHandler) Record(c fiber.Ctx) error {
	var req watchRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	deviceID := c.Get("X-Device-Id")

	if err := h.usecase.Execute(
		c.Context(),
		req.AnimeID,
		req.Episode,
		deviceID,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "recorded",
	})
}
