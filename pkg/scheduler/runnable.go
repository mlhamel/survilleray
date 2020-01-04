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

type group struct {
	cfg *config.Config
}

func (g *group) Run(ctx context.Context) error {
	acquisitionApp := acquisition.NewApp(g.cfg)
	collectionApp := app.NewCollectionApp(g.cfg)
	vectorizeApp := vectorization.NewApp(g.cfg)

	if err := acquisitionApp.Run(ctx); err != nil {
		return err
	}

	if err := collectionApp.Run(ctx); err != nil {
		return err
	}

	if err := vectorizeApp.Run(ctx); err != nil {
		return err
	}

	return nil
}

func NewScheduler(cfg *config.Config) Scheduler {
	return Scheduler{cfg}
}

func (s *Scheduler) Run(ctx context.Context) error {
	return runnable.
		Signal(Periodic(s.cfg, time.Second*25, &group{s.cfg})).
		Run(ctx)
}
