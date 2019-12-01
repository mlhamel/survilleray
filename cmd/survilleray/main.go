package main

import (
	"log"
	"os"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/app"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/vectorization"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cfg := config.NewConfig()

	cliApp.Commands = []cli.Command{
		{
			Name: "acquire",
			Action: func(context *cli.Context) error {
				acquisitionApp := acquisition.NewApp(cfg)
				return acquisitionApp.Run()
			},
		},
		{
			Name: "migrate",
			Action: func(context *cli.Context) error {
				migrateApp := app.NewMigrateApp(cfg)
				return migrateApp.Run()
			},
		},
		{
			Name: "server",
			Action: func(context *cli.Context) error {
				serverApp := app.NewServerApp(cfg)
				return serverApp.Run()
			},
		},
		{
			Name: "vectorize",
			Action: func(context *cli.Context) error {
				vectorizeApp := vectorization.NewApp(cfg)
				return vectorizeApp.Run()
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
