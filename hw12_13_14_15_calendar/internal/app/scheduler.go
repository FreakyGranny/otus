package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/rs/zerolog/log"
)

const jsonContentType = "application/json"

// SchedulerApp scheduler logic.
type SchedulerApp struct {
	storage     storage.Storage
	prod        queue.Producer
	interval    time.Duration
	cleanupDays int
}

// NewSchedulerApp returns new scheduler app.
func NewSchedulerApp(s storage.Storage, q queue.Producer, i time.Duration, cd int) *SchedulerApp {
	return &SchedulerApp{
		storage:     s,
		prod:        q,
		interval:    i,
		cleanupDays: cd,
	}
}

// Watch for events for notification.
func (a *SchedulerApp) Watch(ctx context.Context) error {
	var err error
	if err = a.prod.Start(); err != nil {
		return err
	}
	log.Info().Msg("scheduler started...")

	ticker := time.NewTicker(a.interval)
	defer ticker.Stop()

	n := queue.Notification{}
	m := queue.Message{
		ContentType: jsonContentType,
	}
	for {
		select {
		case <-ticker.C:
			log.Debug().
				Msg("check events for notification")
			events, err := a.storage.GetEventForNotification(ctx, a.interval)
			if err != nil {
				log.Error().
					Err(err).
					Msg("error while getting events for notification")
				continue
			}
			for _, e := range events {
				n.ID = e.ID
				n.Title = e.Title
				n.StartDate = e.StartDate
				n.OwnerID = e.OwnerID
				m.Body, err = json.Marshal(n)
				if err != nil {
					log.Error().
						Err(err).
						Msg("error while marshalling notification")
					continue
				}
				if err = a.prod.Publish(ctx, m); err != nil {
					log.Error().
						Err(err).
						Msg("error while publish message")
				}
			}
			if err := a.storage.DeleteOldEvents(ctx, a.cleanupDays); err != nil {
				log.Error().
					Err(err).
					Msg("error while cleaning old events")
			}
		case <-ctx.Done():
			a.prod.Stop()

			return nil
		}
	}
}
