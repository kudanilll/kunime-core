package app

import (
	"context"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type Repository interface {
	GetVersion(ctx context.Context) (*domain.AppVersion, error)
}

type GetVersionUsecase struct {
	repo Repository
}

func NewGetVersionUsecase(repo Repository) *GetVersionUsecase {
	return &GetVersionUsecase{repo: repo}
}

func (u *GetVersionUsecase) Execute(ctx context.Context) (*domain.AppVersion, error) {
	return u.repo.GetVersion(ctx)
}
