package firestore

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/kudanilll/kunime-core/internal/domain"
)

type AppRepository struct {
	client *firestore.Client
}

func NewAppRepository(client *firestore.Client) *AppRepository {
	return &AppRepository{client: client}
}

func (r *AppRepository) GetVersion(ctx context.Context) (*domain.AppVersion, error) {
	doc, err := r.client.
		Collection("app_config").
		Doc("mobile").
		Get(ctx)
	if err != nil {
		return nil, err
	}

	var result domain.AppVersion
	if err := doc.DataTo(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
