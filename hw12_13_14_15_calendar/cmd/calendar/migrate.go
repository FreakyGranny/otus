package main

import (
	"database/sql"

	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/helper"
	"github.com/FreakyGranny/otus/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/FreakyGranny/otus/hw12_13_14_15_calendar/migrations"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const usageText = `Supported commands are:
  - up - runs all available migrations.
  - down - reverts last migration.
  - version - prints current db version.
`

// NewMigrateCmd returns migrate cmd struct.
func NewMigrateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "run migrations",
		Long:  usageText,
		Run:   Migrate,
	}
}

// Migrate command for run migrations.
func Migrate(cmd *cobra.Command, args []string) {
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
	if config.DB.Type != "postgres" {
		log.Fatal().
			Msg("Migrations supported only postgress")
	}

	dsn := helper.BuildDsn(
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
		config.DB.SSLMode,
	)
	err = goose.SetDialect("postgres")
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect")
	}
	defer db.Close()

	if err := goose.Run(args[0], db, "./"); err != nil {
		log.Error().Err(err).Msg("goose run")
	}
}
