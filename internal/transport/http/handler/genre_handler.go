package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/kudanilll/kunime-core/internal/domain"
	"github.com/kudanilll/kunime-core/internal/usecase/genre"
)

type GenreHandler struct {
	uc *genre.GetAllUsecase
}

func NewGenreHandler(uc *genre.GetAllUsecase) *GenreHandler {
	return &GenreHandler{uc: uc}
}

func (h *GenreHandler) Get(c fiber.Ctx) error {
	data, err := h.uc.Execute(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch genres",
		})
	}

	if data == nil {
		data = []domain.GenreAsset{}
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}
