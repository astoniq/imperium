package objecttype

import (
	"context"
	"fmt"
	"github.com/astoniq/imperium/pkg/database"
	"github.com/pkg/errors"
)

type Repository interface {
	Create(ctx context.Context, objectType Model) (int64, error)
	GetById(ctx context.Context, id int64) (Model, error)
	GetByTypeId(ctx context.Context, typeId string) (Model, error)
	ListAll(ctx context.Context) ([]Model, error)
	UpdateByTypeId(ctx context.Context, typeId string, objectType Model) error
	DeleteByTypeId(ctx context.Context, typeId string) error
}

func NewRepository(db database.Database) (Repository, error) {
	switch db.Type() {
	case database.TypePostgres:
		postgres, ok := db.(*database.Postgres)
		if !ok {
			return nil, errors.New(fmt.Sprintf("invalid %s database config", database.TypePostgres))
		}
		return NewPostgresRepository(postgres), nil
	default:
		return nil, errors.New(fmt.Sprintf("unsupported database type %s specified", db.Type()))
	}
}
