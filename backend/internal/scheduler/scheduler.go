// Package scheduler runs in-process background jobs on a daily cadence.
//
// There is no external cron here by design: the backend already runs as a
// single always-on process, so a goroutine started at app boot is enough.
// If a job ever needs to survive the process restarting mid-day, or run
// outside this container, move it out to an external trigger instead.
package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/ubcesports/memberships/internal/service"
	"go.uber.org/fx"
)

var Module = fx.Module("scheduler",
	fx.Invoke(registerExpiryNotificationJob),
)

// Hour of the day (America/Vancouver) the expiry notification job runs at.
const expiryNotificationHour = 9

func registerExpiryNotificationJob(lc fx.Lifecycle, membershipService *service.MembershipService) {
	ctx, cancel := context.WithCancel(context.Background())

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go runDaily(ctx, "expiry-notifications", expiryNotificationHour, func(jobCtx context.Context) {
				membershipService.RunExpiryNotifications(jobCtx)
			})
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})
}

// runDaily calls job once every day at hourOfDay in America/Vancouver time,
// until ctx is cancelled. It does not run immediately on startup, so
// redeploys don't re-trigger the day's job.
func runDaily(ctx context.Context, name string, hourOfDay int, job func(context.Context)) {
	for {
		wait := durationUntilNextRun(time.Now(), hourOfDay)
		timer := time.NewTimer(wait)

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			slog.Info("running scheduled job", "job", name)
			job(ctx)
		}
	}
}

func durationUntilNextRun(now time.Time, hourOfDay int) time.Duration {
	location, err := time.LoadLocation("America/Vancouver")
	if err != nil {
		location = time.UTC
	}

	localNow := now.In(location)
	next := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), hourOfDay, 0, 0, 0, location)
	if !next.After(localNow) {
		next = next.AddDate(0, 0, 1)
	}

	return next.Sub(localNow)
}
