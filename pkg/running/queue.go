package running

import (
	"context"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/pior/runnable"
)

func Queue(cfg *config.Config, runners ...runnable.Runnable) runnable.Runnable {
	return &groupRunnable{cfg, runners}
}

type groupRunnable struct {
	cfg     *config.Config
	runners []runnable.Runnable
}

func (g *groupRunnable) Run(ctx context.Context) error {
	for i := range g.runners {
		if err := g.runners[i].Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
