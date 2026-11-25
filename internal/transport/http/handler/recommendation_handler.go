package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kudanilll/kunime-core/internal/domain"
	"github.com/kudanilll/kunime-core/internal/usecase/recommendation"
	"go.uber.org/zap"
)

type RecommendationHandler struct {
	getUC   *recommendation.GetRecommendationUsecase
	buildUC *recommendation.BuildRecommendationUsecase
	logger  *zap.Logger
}

func NewRecommendationHandler(
	getUC *recommendation.GetRecommendationUsecase,
	buildUC *recommendation.BuildRecommendationUsecase,
	logger *zap.Logger,
) *RecommendationHandler {
	return &RecommendationHandler{
		getUC:   getUC,
		buildUC: buildUC,
		logger:  logger,
	}
}

func (h *RecommendationHandler) Get(c fiber.Ctx) error {
	items, err := h.getUC.Execute(c.Context())
	if err != nil {
		h.logger.Error("recommendation get failed", zap.Error(err))
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

func (h *RecommendationHandler) Refresh(c fiber.Ctx) error {
	h.logger.Info("recommendation manual refresh requested")
	if err := h.buildUC.Execute(c.Context()); err != nil {
		h.logger.Error("recommendation manual refresh failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to refresh recommendations",
		})
	}

	h.logger.Info("recommendation manual refresh succeeded")
	return c.JSON(fiber.Map{
		"message": "recommendations refreshed",
	})
}
