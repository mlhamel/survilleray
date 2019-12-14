package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/vectorization"
	"github.com/pior/runnable"
)

type Scheduler struct {
	cfg *config.Config
}

func NewScheduler(cfg *config.Config) Scheduler {
	return Scheduler{cfg}
}

func (scheduler *Scheduler) Run(ctx context.Context) error {
	log.Printf("Running Survilleray Scheduler")
	acquisitionApp := acquisition.NewApp(scheduler.cfg)
	vectorizeApp := vectorization.NewApp(scheduler.cfg)

	group := runnable.Group(
		vectorizeApp,
		acquisitionApp,
	)

	return runnable.
		Signal(Periodic(time.Minute, group)).
		Run(ctx)
}
