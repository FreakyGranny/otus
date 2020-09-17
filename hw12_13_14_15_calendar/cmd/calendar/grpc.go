package main

import (
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/app"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/helper"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// NewGRPCcmd return grpc command.
func NewGRPCcmd() *cobra.Command {
	return &cobra.Command{
		Use:   "grpc",
		Short: "run grpc",
		Long:  "starts GRPC server",
		Run:   GRPC,
	}
}

// GRPC starts http GRPC server.
func GRPC(cmd *cobra.Command, args []string) {
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
	storage, err := helper.CreateStorage(
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
	defer storage.Close()

	lsn, err := net.Listen("tcp", net.JoinHostPort(config.GRPC.Host, strconv.Itoa(config.GRPC.Port)))
	if err != nil {
		log.Error().Err(err).Msg("")

		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(internalgrpc.LoggingInterceptor))
	service := internalgrpc.New(app.New(storage))
	internalgrpc.RegisterEventsServer(server, service)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)

		server.GracefulStop()
	}()

	log.Info().Msgf("Starting server on %s", lsn.Addr().String())

	if err := server.Serve(lsn); err != nil {
		log.Error().Err(err).
			Msg("failed to start http server")

		return
	}
}
