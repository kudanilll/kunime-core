package job

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
)

type CleanupWatchEventJob struct {
	client *firestore.Client
	logger *zap.Logger
}

func NewCleanupWatchEventJob(
	client *firestore.Client,
	logger *zap.Logger,
) *CleanupWatchEventJob {
	return &CleanupWatchEventJob{
		client: client,
		logger: logger,
	}
}

func (j *CleanupWatchEventJob) Run(ctx context.Context) error {
	now := time.Now()

	query := j.client.
		Collection("watch_event").
		Where("expires_at", "<=", now).
		Limit(500)

	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	if len(docs) == 0 {
		j.logger.Info("cleanup_watch_event: no expired documents")
		return nil
	}

	batch := j.client.Batch()

	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}

	if _, err := batch.Commit(ctx); err != nil {
		return err
	}

	j.logger.Info(
		"cleanup_watch_event: deleted expired documents",
		zap.Int("count", len(docs)),
	)

	return nil
}
