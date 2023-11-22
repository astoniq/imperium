package service

type Service interface {
	Routes() ([]Route, error)
}
