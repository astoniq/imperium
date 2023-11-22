package objecttype

import "context"

type Service interface {
	Create(ctx context.Context, spec CreateSpec) (*Spec, error)
}

type ObjectTypeService struct {
	Repository Repository
}

func NewService(repository Repository) *ObjectTypeService {
	return &ObjectTypeService{
		Repository: repository,
	}
}
