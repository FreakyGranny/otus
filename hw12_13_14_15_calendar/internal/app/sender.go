package app

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/queue"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

// SendWorker func for process notification.
func SendWorker(ctx context.Context, wg *sync.WaitGroup, ch <-chan amqp.Delivery) {
	var err error
	defer wg.Done()
	notification := &queue.Notification{}

	log.Info().Msg("send worker started...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		select {
		case <-ctx.Done():
			return
		case x := <-ch:
			if err := x.Ack(true); err != nil {
				log.Error().
					Err(err).
					Msg("error while acking message")
				continue
			}
			if err = json.Unmarshal(x.Body, notification); err != nil {
				log.Error().
					Err(err).
					Bytes("raw", x.Body).
					Msg("error while unmarshalling message")
				continue
			}
			log.Info().
				Int64("ID", notification.ID).
				Str("title", notification.Title).
				Int64("user", notification.OwnerID).
				Str("date", notification.StartDate.String()).
				Msg("")
		}
	}
}
