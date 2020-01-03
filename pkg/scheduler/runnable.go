package scheduler

import (
	"context"
	"time"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/app"
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

func (s *Scheduler) Run(ctx context.Context) error {
	acquisitionApp := acquisition.NewApp(s.cfg)
	collectionApp := app.NewCollectionApp(s.cfg)
	vectorizeApp := vectorization.NewApp(s.cfg)

	group := runnable.Group(
		vectorizeApp,
		collectionApp,
		acquisitionApp,
	)

	return runnable.
		Signal(Periodic(s.cfg, time.Minute, group)).
		Run(ctx)
}
