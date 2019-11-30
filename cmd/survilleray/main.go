package main

import (
	"log"
	"os"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/app"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/runtime"
	"github.com/urfave/cli"
)

func main() {
	cliApp := cli.NewApp()
	cfg := config.NewConfig()
	context := runtime.NewContext(cfg, nil)

	cliApp.Commands = []cli.Command{
		{
			Name: "acquire",
			Action: func(c *cli.Context) error {
				acquisitionApp := acquisition.NewApp(context)
				return acquisitionApp.Run()
			},
		},
		{
			Name: "migrate",
			Action: func(c *cli.Context) error {
				migrateApp := app.NewMigrateApp(context)
				return migrateApp.Run()
			},
		},
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				serverApp := app.NewServerApp(context)
				return serverApp.Run()
			},
		},
		{
			Name: "vectorize",
			Action: func(c *cli.Context) error {
				vectorizeApp := app.NewVectorizeApp(context)
				return vectorizeApp.Run()
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
