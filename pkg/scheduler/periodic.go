package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/pior/runnable"
)

func Periodic(duration time.Duration, runner runnable.Runnable) runnable.Runnable {
	log.Printf("Initializing periodic runner for %s", fmtDuration(duration))
	return &periodicRunnable{duration, runner}
}

type periodicRunnable struct {
	duration time.Duration
	runner   runnable.Runnable
}

func (periodic *periodicRunnable) Run(ctx context.Context) error {
	log.Printf("Running periodic runner")

	errs := make(chan error, 1)
	done := make(chan bool)

	defer close(errs)
	defer close(done)

	go periodic.runForever(ctx, errs, done)

	select {
	case <-ctx.Done():
		log.Printf("Finishing running periodic runner")
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
		log.Printf("Running again...")
		select {
		case <-done:
			log.Printf("Finishing running forever")
			return
		case <-ticker.C:
			log.Printf("Waking up, running underlying runnable")

			if err := periodic.runner.Run(ctx); err != nil {
				errs <- err
				return
			}
		}
	}
	log.Printf("Finishing")
}
