package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kudanilll/kunime-core/internal/domain"
	"github.com/kudanilll/kunime-core/internal/usecase/recommendation"
)

type RecommendationHandler struct {
	uc *recommendation.GetRecommendationUsecase
}

func NewRecommendationHandler(
	uc *recommendation.GetRecommendationUsecase,
) *RecommendationHandler {
	return &RecommendationHandler{uc: uc}
}

func (h *RecommendationHandler) Get(c *fiber.Ctx) error {
	items, err := h.uc.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch recommendations",
		})
	}

	if items == nil {
		items = make([]domain.RecommendationItem, 0)
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}
