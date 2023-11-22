package service

import (
	"fmt"
	"github.com/astoniq/imperium/pkg/config"
	"github.com/gorilla/mux"
)

func NewRouter(config config.Config, pathPrefix string, routes []Route) (*mux.Router, error) {
	router := mux.NewRouter()

	//Setup routes
	for _, route := range routes {
		routePattern := fmt.Sprintf("%s%s", pathPrefix, route.GetPattern())
		router.Handle(routePattern, route.GetHandler()).Methods(route.GetMethod())
	}

	return router, nil
}
