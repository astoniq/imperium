package service

import "net/http"

type Route interface {
	GetPattern() string
	GetMethod() string
	GetHandler() http.Handler
}
