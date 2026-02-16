package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/kudanilll/kunime-core/internal/usecase/popular"
)

type PopularHandler struct {
	usecase *popular.GetPopularUsecase
}

func NewPopularHandler(uc *popular.GetPopularUsecase) *PopularHandler {
	return &PopularHandler{usecase: uc}
}

func (h *PopularHandler) Get(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	items, err := h.usecase.Execute(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch popular anime",
		})
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}
