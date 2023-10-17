package main

import (
	"go.uber.org/fx"
	"trust-verse-backend/app"
	"trust-verse-backend/app/module/user"
	"trust-verse-backend/app/router"
)

func main() {
	fx.New(
		fx.Provide(app.NewConfig),
		fx.Provide(app.NewApp),
		fx.Provide(router.NewRouter),
		user.NewUserModule,
		fx.Invoke(app.Start),
		//fx.WithLogger(fxzerolog.Init()),
	).Run()
}
