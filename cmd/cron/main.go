package main

import (
	stdlog "log"
	"os"

	"github.com/go-logr/stdr"
	"github.com/robfig/cron/v3"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/jobs"
)

func main() {
	c := config.NewConfig()
	logger := stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile))
	runner := cron.New(cron.WithLogger(logger))

	job := jobs.NewAcquisition(c)

	runner.AddJob("@every 1m", job)
	runner.Run()
}
