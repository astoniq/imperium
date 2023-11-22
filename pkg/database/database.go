package database

import (
	"context"
	"github.com/astoniq/imperium/pkg/config"
	"github.com/pkg/errors"
	"time"
)

const (
	TypePostgres = "postgres"
	TypeSQLite   = "sqlite"
)

type Database interface {
	Type() string
	Connect(ctx context.Context) error
	Migrate(ctx context.Context) error
	Ping(ctx context.Context) error
}

func NewDatabase(cfg config.Config) (Database, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	if cfg.GetDatastore().Postgres.Hostname != "" {
		db := NewPostgres(*cfg.GetDatastore().Postgres)
		err := db.Connect(ctx)
		if err != nil {
			return nil, err
		}

		err = db.Migrate(ctx)
		if err != nil {
			return nil, err
		}
		return db, nil
	}

	return nil, errors.New("invalid database configuration provided")
}
