package scheduler

import (
	"context"
	"time"

	"github.com/pior/runnable"
)

func Periodic(duration time.Duration, runner runnable.Runnable) runnable.Runnable {
	return &periodicRunnable{duration, runner}
}

type periodicRunnable struct {
	duration time.Duration
	runner   runnable.Runnable
}

func (periodic *periodicRunnable) Run(ctx context.Context) error {
	errs := make(chan error, 1)
	done := make(chan bool)

	defer close(errs)
	defer close(done)

	go periodic.runForever(ctx, errs, done)

	select {
	case <-ctx.Done():
		done <- true
		return nil
	case err := <-errs:
		return err
	}
}

func (periodic *periodicRunnable) runForever(ctx context.Context, errs chan<- error, done chan bool) {
	ticker := time.NewTicker(periodic.duration)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := periodic.runner.Run(ctx); err != nil {
				errs <- err
				return
			}
		}
	}
}
