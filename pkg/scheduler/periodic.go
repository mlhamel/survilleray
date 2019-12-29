package scheduler

import (
	"context"
	"time"

	"github.com/pior/runnable"
	"github.com/rs/zerolog/log"
)

func Periodic(duration time.Duration, runner runnable.Runnable) runnable.Runnable {
	log.Info().Str("duration", fmtDuration(duration)).Msg("Initializing periodic runner")
	return &periodicRunnable{duration, runner}
}

type periodicRunnable struct {
	duration time.Duration
	runner   runnable.Runnable
}

func (periodic *periodicRunnable) Run(ctx context.Context) error {
	log.Info().Msg("Running periodic runner")

	errs := make(chan error, 1)
	done := make(chan bool)

	defer close(errs)
	defer close(done)

	go periodic.runForever(ctx, errs, done)

	select {
	case <-ctx.Done():
		log.Info().Msg("Finishing running periodic runner")
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
		log.Info().Msg("Running again...")
		select {
		case <-done:
			log.Info().Msg("Finishing running forever")
			return
		case <-ticker.C:
			log.Info().Msg("Waking up, running underlying runnable")

			if err := periodic.runner.Run(ctx); err != nil {
				log.Error().Err(err).Msg("Error while running the runnable")
				errs <- err
				return
			}
		}
	}
}
