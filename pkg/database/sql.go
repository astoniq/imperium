package database

import "github.com/rs/zerolog/log"

type SQL struct {
}

type SQLRepository struct {
	DB *SQL
}

func NewSQLRepository(db *SQL) SQLRepository {
	if db == nil {
		log.Fatal().Msg("init: cannot initialize SQLRepository with a nil db parameter")
	}

	return SQLRepository{
		DB: db,
	}
}
