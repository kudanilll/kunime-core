package app

import (
	"context"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type UpdateRepository interface {
	UpdateVersion(ctx context.Context, data *domain.AppVersion) error
}

type UpdateVersionUsecase struct {
	repo UpdateRepository
}

func NewUpdateVersionUsecase(repo UpdateRepository) *UpdateVersionUsecase {
	return &UpdateVersionUsecase{repo: repo}
}

func (u *UpdateVersionUsecase) Execute(
	ctx context.Context,
	data *domain.AppVersion,
) error {
	return u.repo.UpdateVersion(ctx, data)
}
