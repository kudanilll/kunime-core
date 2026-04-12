package bootstrap

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/kudanilll/kunime-core/internal/config"
	firestoreInfra "github.com/kudanilll/kunime-core/internal/infra/firestore"
	httpInfra "github.com/kudanilll/kunime-core/internal/infra/http"
	"github.com/kudanilll/kunime-core/internal/job"
	httpTransport "github.com/kudanilll/kunime-core/internal/transport/http"
	"github.com/kudanilll/kunime-core/internal/transport/http/handler"
	appUsecase "github.com/kudanilll/kunime-core/internal/usecase/app"
	bannerUsecase "github.com/kudanilll/kunime-core/internal/usecase/banner"
	genreUsecase "github.com/kudanilll/kunime-core/internal/usecase/genre"
	"github.com/kudanilll/kunime-core/internal/usecase/popular"
	"github.com/kudanilll/kunime-core/internal/usecase/recommendation"
	"github.com/kudanilll/kunime-core/internal/usecase/watch"
)

type App struct {
	Fiber     *fiber.App
	Logger    *zap.Logger
	Firestore *firebase.App
}

func NewApp() (*App, error) {
	logger, err := zap.NewProduction()
	ctx := context.Background()
	cfg, err := config.Load()

	if err != nil {
		return nil, err
	}

	// Firebase init
	credentialsJSON, err := decodeFirebaseCredentials(cfg.FirebaseCredentialsBase64)
	if err != nil {
		return nil, err
	}

	fbApp, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		logger.Fatal("failed to initialize firebase app", zap.Error(err))
		return nil, err
	}

	firestoreClient, err := firestoreInfra.NewClient(ctx, fbApp)
	if err != nil {
		return nil, err
	}

	// Cron job init
	c := cron.New()
	cleanupJob := job.NewCleanupWatchEventJob(firestoreClient, logger)
	c.AddFunc("@every 1h", func() {
		if err := cleanupJob.Run(ctx); err != nil {
			logger.Error("cleanup job failed", zap.Error(err))
		}
	})

	recRepo := firestoreInfra.NewRecommendationRepository(firestoreClient)
	apiClient := httpInfra.NewKunimeAPIClient(
		cfg.KunimeAPIBaseURL,
		cfg.KunimeAPIKey,
	)

	recUC := recommendation.NewBuildRecommendationUsecase(apiClient, recRepo)

	c.AddFunc("@every 1h", func() {
		if err := recUC.Execute(ctx); err != nil {
			logger.Error("recommendation cron failed", zap.Error(err))
		}
	})

	c.Start()

	popularRepo := firestoreInfra.NewPopularRepository(firestoreClient)
	popularUC := popular.NewGetPopularUsecase(popularRepo)
	popularHandler := handler.NewPopularHandler(popularUC)

	watchRepo := firestoreInfra.NewWatchRepository(firestoreClient)
	watchUC := watch.NewRecordWatchUsecase(watchRepo, popularRepo)
	watchHandler := handler.NewWatchHandler(watchUC)

	recReadUC := recommendation.NewGetRecommendationUsecase(recRepo)
	recHandler := handler.NewRecommendationHandler(recReadUC)

	appRepo := firestoreInfra.NewAppRepository(firestoreClient)
	appGetUC := appUsecase.NewGetVersionUsecase(appRepo)
	appUpdateUC := appUsecase.NewUpdateVersionUsecase(appRepo)
	appHandler := handler.NewAppHandler(appGetUC, appUpdateUC)

	bannerRepo := firestoreInfra.NewBannerRepository(firestoreClient)
	bannerUC := bannerUsecase.NewGetActiveUsecase(bannerRepo)
	bannerHandler := handler.NewBannerHandler(bannerUC)

	genreRepo := firestoreInfra.NewGenreRepository(firestoreClient)
	genreUC := genreUsecase.NewGetAllUsecase(genreRepo)
	genreHandler := handler.NewGenreHandler(genreUC)

	handlers := &httpTransport.Handlers{
		Watch:          watchHandler,
		Popular:        popularHandler,
		Recommendation: recHandler,
		App:            appHandler,
		Banner:         bannerHandler,
		Genre:          genreHandler,
	}

	if err := recUC.Execute(ctx); err != nil {
		logger.Error("initial recommendation build failed", zap.Error(err))
	} else {
		logger.Info("initial recommendation built")
	}

	// Fiber init
	app := fiber.New(fiber.Config{
		AppName: "kunime-core",
	})

	// Register routes
	httpTransport.RegisterRoutes(app, handlers)

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

func decodeFirebaseCredentials(value string) ([]byte, error) {
	if value == "" {
		return nil, errors.New("FIREBASE_CREDENTIALS_BASE64 is required")
	}

	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid FIREBASE_CREDENTIALS_BASE64: %w", err)
	}

	return decoded, nil
}
