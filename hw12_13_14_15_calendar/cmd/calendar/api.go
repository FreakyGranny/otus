package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createStorage(config *Config) (app.Storage, error) {
	switch strings.ToLower(config.DB.Type) {
	case "memory":
		return memorystorage.New(), nil
	case "postgres":
		dsn := sqlstorage.BuildDsn(
			config.DB.Host,
			config.DB.Port,
			config.DB.User,
			config.DB.Password,
			config.DB.Name,
			config.DB.SSLMode,
		)

		return sqlstorage.New(dsn), nil
	default:
		return nil, errors.New("unsupported storage type")
	}
}

// NewAPICmd return api command.
func NewAPICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "run api",
		Long:  "starts API server",
		Run:   API,
	}
}

// API starts http API server.
func API(cmd *cobra.Command, args []string) {
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
	storage, err := createStorage(config)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize storage")
	}
	defer storage.Close()

	server := internalhttp.NewServer(
		net.JoinHostPort(config.HTTP.Host, strconv.Itoa(config.HTTP.Port)),
		app.New(storage),
	)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Err(err).
				Msg("failed to stop http server")
		}
	}()

	if err := server.Start(); err != nil {
		log.Err(err).
			Msg("failed to start http server")
		os.Exit(1)
	}
}
