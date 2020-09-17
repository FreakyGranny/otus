package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
)

// Config confuguration struct.
type Config struct {
	HTTP struct {
		Host string `config:"http-host,required"`
		Port int    `config:"http-port,required"`
	}
	GRPC struct {
		Host string `config:"grpc-host,required"`
		Port int    `config:"grpc-port,required"`
	}
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
