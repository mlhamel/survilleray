package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/mlhamel/survilleray/pkg/app"
	"github.com/mlhamel/survilleray/pkg/config"
)

func main() {
	cliApp := cli.NewApp()
	cfg := config.NewConfig()

	cliApp.Commands = []cli.Command{
		{
			Name: "acquire",
			Action: func(c *cli.Context) error {
				acquisitionApp := app.NewAcquisitionApp(cfg)
				return acquisitionApp.Run()
			},
		},
		{
			Name: "migrate",
			Action: func(c *cli.Context) error {
				migrateApp := app.NewMigrateApp(cfg)
				return migrateApp.Run()
			},
		},
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				serverApp := app.NewServerApp(cfg)
				return serverApp.Run()
			},
		},
		{
			Name: "vectorize",
			Action: func(c *cli.Context) error {
				vectorizeApp := app.NewVectorizeApp(cfg)
				return vectorizeApp.Run()
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
