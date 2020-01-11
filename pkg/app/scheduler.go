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

const DURATION = 15 * time.Second

type Scheduler struct {
	cfg     *config.Config
	timeout time.Duration
}

func NewScheduler(cfg *config.Config) Scheduler {
	return Scheduler{cfg, DURATION}
}

func (s *Scheduler) Run(ctx context.Context) error {
	acquisition := acquisition.NewApp(s.cfg)
	vectorization := vectorization.NewApp(s.cfg)
	collection := NewCollectionApp(s.cfg)

	wrapper := running.Queue(s.cfg, acquisition, vectorization, collection)
	periodic := running.Periodic(s.cfg, s.timeout, wrapper)

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
