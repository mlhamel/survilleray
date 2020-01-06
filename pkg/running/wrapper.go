package running

import (
	"context"

	"github.com/mlhamel/survilleray/pkg/config"
	"github.com/pior/runnable"
)

func Wrapper(cfg *config.Config, wrapper func(*config.Config, context.Context, runnable.Runnable) error, runners ...runnable.Runnable) runnable.Runnable {
	return &wrapperRunnable{cfg, wrapper, runners}
}

type wrapperRunnable struct {
	cfg     *config.Config
	wrapper func(*config.Config, context.Context, runnable.Runnable) error
	runners []runnable.Runnable
}

func (w *wrapperRunnable) Run(ctx context.Context) error {
	for i := range w.runners {
		if err := w.wrapper(w.cfg, ctx, w.runners[i]); err != nil {
			return err
		}
	}
	return nil
}
