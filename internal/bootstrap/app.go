package bootstrap

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/gofiber/fiber/v3"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	Fiber      *fiber.App
	Logger     *zap.Logger
	Firestore  *firebase.App
	FireClient *firestore.Client
	Cron       *cron.Cron
}

func NewApp() (*App, error) {
	ctx := context.Background()
	cfg, err := config.Load()

	if err != nil {
		return nil, err
	}

	logger, err := newLogger(cfg.Env)
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
		logger.Error("failed to initialize firebase app", zap.Error(err))
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
		// Fix M-06: Use context with timeout for cron jobs
		jobCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := cleanupJob.Run(jobCtx); err != nil {
			logger.Error("cleanup job failed", zap.Error(err))
		}
	})

	recRepo := firestoreInfra.NewRecommendationRepository(firestoreClient)
	apiClient := httpInfra.NewKunimeAPIClient(
		cfg.KunimeAPIBaseURL,
		cfg.KunimeAPIKey,
		logger,
	)

	recUC := recommendation.NewBuildRecommendationUsecase(apiClient, recRepo, logger)

	c.AddFunc("@every 1h", func() {
		// Fix M-06: Use context with timeout for cron jobs
		jobCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		if err := recUC.Execute(jobCtx); err != nil {
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
	recHandler := handler.NewRecommendationHandler(recReadUC, recUC, logger)

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
		Fiber:      app,
		Logger:     logger,
		Firestore:  fbApp,
		FireClient: firestoreClient,
		Cron:       c,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", config.AppPort())

	// Start server in goroutine
	go func() {
		a.Logger.Info("starting server", zap.String("addr", addr))
		if err := a.Fiber.Listen(addr); err != nil {
			a.Logger.Error("server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Logger.Info("shutting down gracefully...")

	// Stop cron scheduler
	if a.Cron != nil {
		a.Cron.Stop()
	}

	// Shutdown Fiber with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.Fiber.ShutdownWithContext(ctx); err != nil {
		a.Logger.Error("fiber shutdown error", zap.Error(err))
	}

	// Close Firestore
	if a.FireClient != nil {
		if err := a.FireClient.Close(); err != nil {
			a.Logger.Error("firestore close error", zap.Error(err))
		}
	}

	// Sync logger
	if a.Logger != nil {
		_ = a.Logger.Sync()
	}

	return nil
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

func newLogger(env string) (*zap.Logger, error) {
	if env == "development" {
		cfg := zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.LevelKey = "level"
		cfg.EncoderConfig.MessageKey = "msg"
		cfg.EncoderConfig.CallerKey = ""
		cfg.EncoderConfig.StacktraceKey = ""
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		cfg.DisableStacktrace = true
		return cfg.Build()
	}

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.DisableStacktrace = false
	return cfg.Build()
}
