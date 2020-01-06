package app

import (
	"context"
	"time"

	"github.com/mlhamel/survilleray/pkg/acquisition"
	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/mlhamel/survilleray/pkg/running"
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
	acquisition := acquisition.NewApp(s.cfg)
	vectorization := vectorization.NewApp(s.cfg)
	collection := NewCollectionApp(s.cfg)

	wrapper := running.Wrapper(s.cfg, closer, acquisition, vectorization, collection)
	periodic := running.Periodic(s.cfg, time.Minute*5, wrapper)

	return runnable.
		Signal(periodic).
		Run(ctx)
}

func closer(cfg *config.Config, ctx context.Context, runner runnable.Runnable) error {
	if err := runner.Run(ctx); err != nil {
		return err
	}

	cfg.Logger().Debug().Msg("Closing database connection")

	err := cfg.Database().Close()

	cfg.Logger().Debug().Msg("Done: Closing database connection")

	return err
}
