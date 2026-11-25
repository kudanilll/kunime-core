package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kudanilll/kunime-core/internal/domain"
	"github.com/kudanilll/kunime-core/internal/usecase/banner"
)

type BannerHandler struct {
	uc *banner.GetActiveUsecase
}

func NewBannerHandler(uc *banner.GetActiveUsecase) *BannerHandler {
	return &BannerHandler{uc: uc}
}

func (h *BannerHandler) Get(c fiber.Ctx) error {
	data, err := h.uc.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch banners",
		})
	}

	if data == nil {
		data = []domain.Banner{}
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
