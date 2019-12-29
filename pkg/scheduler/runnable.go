package scheduler

import (
	"context"
	"time"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/vectorization"
	"github.com/pior/runnable"
	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	cfg *config.Config
}

func NewScheduler(cfg *config.Config) Scheduler {
	return Scheduler{cfg}
}

func (scheduler *Scheduler) Run(ctx context.Context) error {
	log.Info().Msg("Running Survilleray Scheduler")
	acquisitionApp := acquisition.NewApp(scheduler.cfg)
	vectorizeApp := vectorization.NewApp(scheduler.cfg)

	group := runnable.Group(
		vectorizeApp,
		acquisitionApp,
	)

	return runnable.
		Signal(Periodic(time.Second*15, group)).
		Run(ctx)
}
