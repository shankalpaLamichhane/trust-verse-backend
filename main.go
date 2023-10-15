package main

import (
	"go.uber.org/fx"
	"trust-verse-backend/app"
	"trust-verse-backend/app/router"
)

func main() {
	fx.New(
		fx.Provide(app.NewConfig),
		fx.Provide(app.NewApp),
		fx.Invoke(app.Start),
		fx.Provide(router.NewRouter),
		//fx.WithLogger(fxzerolog.Init()),
	).Run()
}
