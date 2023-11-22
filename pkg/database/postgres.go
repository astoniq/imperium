package database

import (
	"context"
	"github.com/astoniq/imperium/pkg/config"
)

type Postgres struct {
	SQL
	Config config.PostgresConfig
}

func NewPostgres(config config.PostgresConfig) *Postgres {
	return &Postgres{
		Config: config,
	}
}

func (ds *Postgres) Type() string {
	return TypePostgres
}

func (ds *Postgres) Connect(ctx context.Context) error {
	return nil
}

func (ds *Postgres) Migrate(ctx context.Context) error {
	return nil
}

func (ds *Postgres) Ping(ct context.Context) error {
	return nil
}
