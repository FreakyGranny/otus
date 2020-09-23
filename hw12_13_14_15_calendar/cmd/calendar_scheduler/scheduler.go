package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/helper"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/logger"
	rmq "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/queue/rabbit"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// NewSchedulerCmd return scheduler command.
func NewSchedulerCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "calendar_scheduler",
		Short: "run scheduler",
		Long:  "starts scheduler",
		Run:   Scheduler,
	}
}

// Scheduler starts Scheduler.
func Scheduler(cmd *cobra.Command, args []string) {
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
	wakeupInterval, err := time.ParseDuration(config.Scheduler.Interval)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to parse scheduler interval")
	}
	s, err := helper.CreateStorage(
		config.DB.Type,
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize storage")
	}
	defer s.Close()

	p := rmq.NewProducer(
		helper.BuildRmqAddr(config.Rmq.Host, config.Rmq.Port, config.Rmq.User, config.Rmq.Password),
		config.Rmq.ExchangeName,
		config.Rmq.ExchangeType,
		config.Rmq.Queue,
		config.Rmq.BindingKey,
	)
	app := app.NewSchedulerApp(s, p, wakeupInterval, config.Scheduler.CleanupDays)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)

		cancel()
	}()

	if err = app.Watch(ctx); err != nil {
		log.Error().Err(err).Msg("unable to start scheduler")

		return
	}
}
