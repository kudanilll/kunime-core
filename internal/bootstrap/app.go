package bootstrap

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/kudanilll/kunime-core/internal/config"
)

type App struct {
	Fiber  *fiber.App
	Logger *zap.Logger
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

	app := fiber.New(fiber.Config{
		AppName: "kunime-core",
	})

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	logger.Info("application initialized")

	return &App{
		Fiber:  app,
		Logger: logger,
	}, nil
}

func (a *App) Run() error {
	addr := fmt.Sprintf(":%s", config.AppPort())
	a.Logger.Info("starting server", zap.String("addr", addr))
	return a.Fiber.Listen(addr)
}
