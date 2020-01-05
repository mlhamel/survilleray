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
	collection := NewCollectionApp(s.cfg)
	vectorization := vectorization.NewApp(s.cfg)

	queue := running.Queue(s.cfg, acquisition, collection, vectorization)
	periodic := running.Periodic(s.cfg, time.Second, queue)

	return runnable.
		Signal(periodic).
		Run(ctx)
}