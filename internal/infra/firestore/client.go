package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

func NewClient(ctx context.Context, app *firebase.App) (*firestore.Client, error) {
	return app.Firestore(ctx)
}
