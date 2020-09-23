package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/helper"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/logger"
	rmq "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewSenderCmd return sender command.
func NewSenderCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "calendar_sender",
		Short: "run sender",
		Long:  "starts sender",
		Run:   Sender,
	}
}

// Sender starts sender process.
func Sender(cmd *cobra.Command, args []string) {
	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("unable initialize config")
	}
	if err := logger.SetLogLevel(config.Logger.Level); err != nil {
		log.Fatal().
			Err(err).
			Msg("unable to initialize logger")
	}

	c := rmq.NewConsumer(
		config.Rmq.Tag,
		helper.BuildRmqAddr(config.Rmq.Host, config.Rmq.Port, config.Rmq.User, config.Rmq.Password),
		config.Rmq.Queue,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)

		c.Stop()
		log.Info().Msg("stopping")
		cancel()
	}()
	if err := c.Handle(ctx, app.SendWorker, config.Sender.Threads); err != nil {
		log.Err(err).
			Msg("failed to start sender")

		return
	}
}
