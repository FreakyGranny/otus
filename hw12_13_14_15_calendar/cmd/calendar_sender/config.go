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
	Rmq struct {
		Host     string `config:"rmq-host"`
		Port     int    `config:"rmq-port"`
		User     string `config:"rmq-user"`
		Password string `config:"rmq-password"`
		Queue    string `config:"rmq-queue"`
		Tag      string `config:"rmq-tag"`
	}
	Sender struct {
		Threads int `config:"sender-threads"`
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
