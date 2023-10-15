package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"time"
	"trust-verse-backend/app/router"
)

func NewApp(cfg *Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ServerHeader: cfg.App.Name,
		AppName:      cfg.App.Name,
		//ErrorHandler: response.ErrorHandler,
		IdleTimeout: cfg.App.Timeout * time.Second,
	})
	return app
}

func Start(lifecycle fx.Lifecycle, fiber *fiber.App, router *router.Router) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				router.Register()
				port := "8080"
				go func() {
					if err := fiber.Listen(port); err != nil {
						//log.Error().Err(err).Msg("Something went wrong starting the server !")
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().Msg("Shutting down the app...")
				if err := fiber.Shutdown(); err != nil {
					//log.Panic().Err(err).Msg("")
				}
				//log.Info().Msg("Running clean up tasks")
				//log.Info().Msgf("%s was sucessfully shutdown.", cfg.App.Name)
				return nil
			},
		},
	)
}
