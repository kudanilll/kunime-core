package bootstrap

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/kudanilll/kunime-core/internal/config"
	firestoreInfra "github.com/kudanilll/kunime-core/internal/infra/firestore"
	"github.com/kudanilll/kunime-core/internal/job"
	httpTransport "github.com/kudanilll/kunime-core/internal/transport/http"
	"github.com/kudanilll/kunime-core/internal/transport/http/handler"
	"github.com/kudanilll/kunime-core/internal/usecase/popular"
	"github.com/kudanilll/kunime-core/internal/usecase/watch"
)

type App struct {
	Fiber     *fiber.App
	Logger    *zap.Logger
	Firestore *firebase.App
}

func NewApp() (*App, error) {
	_, err := config.Load()
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Firebase init
	fbApp, err := firebase.NewApp(ctx, nil)
	if err != nil {
		logger.Fatal("failed to initialize firebase app", zap.Error(err))
		return nil, err
	}

	firestoreClient, err := firestoreInfra.NewClient(ctx, fbApp)
	if err != nil {
		return nil, err
	}

	popularRepo := firestoreInfra.NewPopularRepository(firestoreClient)
	popularUC := popular.NewGetPopularUsecase(popularRepo)
	popularHandler := handler.NewPopularHandler(popularUC)

	watchRepo := firestoreInfra.NewWatchRepository(firestoreClient)
	watchUC := watch.NewRecordWatchUsecase(watchRepo, popularRepo)
	watchHandler := handler.NewWatchHandler(watchUC)

	handlers := &httpTransport.Handlers{
		Watch:   watchHandler,
		Popular: popularHandler,
	}

	// Fiber init
	app := fiber.New(fiber.Config{
		AppName: "kunime-core",
	})

	// Register routes
	httpTransport.RegisterRoutes(app, handlers)

	// Cron job init
	c := cron.New()
	cleanupJob := job.NewCleanupWatchEventJob(firestoreClient, logger)
	c.AddFunc("@every 1h", func() {
		ctx := context.Background()
		if err := cleanupJob.Run(ctx); err != nil {
			logger.Error("cleanup job failed", zap.Error(err))
		}
	})
	c.Start()

	logger.Info("application initialized")

	return &App{
		Fiber:     app,
		Logger:    logger,
		Firestore: fbApp,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", config.AppPort())
	a.Logger.Info("starting server", zap.String("addr", addr))
	return a.Fiber.Listen(addr)
}
