package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

// Config confuguration struct.
type Config struct {
	Logger struct {
		Level string `config:"logger-level,required"`
	}
	DB struct {
		Type     string `config:"db-type,required"`
		User     string `config:"db-user"`
		Password string `config:"db-password"`
		Host     string `config:"db-host"`
		Port     int    `config:"db-port"`
		Name     string `config:"db-name"`
		SSLMode  bool   `config:"db-sslmode"`
	}
	Rmq struct {
		Host         string `config:"rmq-host"`
		Port         int    `config:"rmq-port"`
		User         string `config:"rmq-user"`
		Password     string `config:"rmq-password"`
		ExchangeName string `config:"rmq-exchangename"`
		ExchangeType string `config:"rmq-exchangetype"`
		Queue        string `config:"rmq-queue"`
		BindingKey   string `config:"rmq-bindingkey"`
	}
	Scheduler struct {
		Interval    string `config:"scheduler-interval"`
		CleanupDays int    `config:"scheduler-cleanupdays"`
	}
}

// NewConfig returns configuratio.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	loader := confita.NewLoader(
		file.NewBackend(path),
		env.NewBackend(),
	)
	if err := loader.Load(context.Background(), cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
