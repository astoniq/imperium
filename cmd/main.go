package main

import (
	"fmt"
	"github.com/astoniq/imperium/pkg/config"
	"github.com/astoniq/imperium/pkg/database"
	"github.com/astoniq/imperium/pkg/objecttype"
	"github.com/astoniq/imperium/pkg/service"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {

	cfg := config.NewConfig()
	db, err := database.NewDatabase(cfg)

	if err != nil {
		log.Fatal().Err(err).Msg(
			"init: could not initialize and connect to the configured datastore. Shutting down.")
	}

	objectTypeRepository, err := objecttype.NewRepository(db)
	if err != nil {
		log.Fatal().Err(err).Msg("init: could not initialize ObjectTypeRepository")
	}

	objectTypeService := objecttype.NewService(objectTypeRepository)

	services := []service.Service{
		objectTypeService,
	}

	routes := make([]service.Route, 0)
	for _, svc := range services {
		svcRoutes, err := svc.Routes()
		if err != nil {
			log.Fatal().Err(err).Msg("init: could not setup routes for service")
		}
		routes = append(routes, svcRoutes...)
	}

	router, err := service.NewRouter(cfg, "", routes)
	if err != nil {
		log.Fatal().Err(err).Msg("init: could not initialize service router")
	}

	log.Info().Msgf("init: listening on port %d", cfg.GetPort())
	shutdownErr := http.ListenAndServe(fmt.Sprintf(":%d", cfg.GetPort()), router)
	log.Fatal().Err(shutdownErr).Msg("shutdown")
}
