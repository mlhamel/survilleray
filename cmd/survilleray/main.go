package main

import (
	"context"
	"os"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/app"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/scheduler"
	"github.com/mlhamel/survilleray/pkg/vectorization"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cfg := config.NewConfig()
	ctx := context.Background()

	defer cfg.Database().Close()

	cliApp.Commands = []cli.Command{
		{
			Name: "acquire",
			Action: func(*cli.Context) error {
				acquisitionApp := acquisition.NewApp(cfg)
				return acquisitionApp.Run(ctx)
			},
		},
		{
			Name: "collect",
			Action: func(*cli.Context) error {
				collectionApp := app.NewCollectionApp(cfg)
				return collectionApp.Run(ctx)
			},
		},
		{
			Name: "migrate",
			Action: func(*cli.Context) error {
				migrateApp := app.NewMigrateApp(cfg)
				return migrateApp.Run(ctx)
			},
		},
		{
			Name: "server",
			Action: func(*cli.Context) error {
				serverApp := app.NewServerApp(cfg)
				return serverApp.Run(ctx)
			},
		},
		{
			Name: "vectorize",
			Action: func(*cli.Context) error {
				vectorizeApp := vectorization.NewApp(cfg)
				return vectorizeApp.Run(ctx)
			},
		},
		{
			Name: "schedule",
			Action: func(*cli.Context) error {
				scheduler := scheduler.NewScheduler(cfg)
				return scheduler.Run(ctx)
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		cfg.Logger().Error().Err(err)
	}
}
