package service

import (
    "context"
	"{{ .ProjectName }}/internal/model"
	"{{ .ProjectName }}/internal/repository"
)

type {{ .StructName }}Service interface {
	Create{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error
	Get{{ .StructName }}(ctx context.Context, id string) (*model.{{ .StructName }}, error)
	Update{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error
	Delete{{ .StructName }}(ctx context.Context, id string) error
	List{{ .StructName }}s(ctx context.Context, filter map[string]interface{}) ([]*model.{{ .StructName }}, error)
	BatchGet{{ .StructName }}s(ctx context.Context, ids []string) ([]*model.{{ .StructName }}, error)
	BatchCreate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error
	BatchUpdate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error
	BatchDelete{{ .StructName }}s(ctx context.Context, ids []string) error
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}

// @wire:Service
func New{{ .StructName }}Service(
    service *Service,
    {{ .StructNameLowerFirst }}Repository repository.{{ .StructName }}Repository,
) {{ .StructName }}Service {
	return &{{ .StructNameLowerFirst }}Service{
		Service:        service,
		{{ .StructNameLowerFirst }}Repository: {{ .StructNameLowerFirst }}Repository,
	}
}

type {{ .StructNameLowerFirst }}Service struct {
	*Service
	{{ .StructNameLowerFirst }}Repository repository.{{ .StructName }}Repository
}

func (s *{{ .StructNameLowerFirst }}Service) Create{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.Create{{ .StructName }}(ctx, {{ .StructNameLowerFirst }})
}

func (s *{{ .StructNameLowerFirst }}Service) Get{{ .StructName }}(ctx context.Context, id string) (*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.Get{{ .StructName }}(ctx, id)
}

func (s *{{ .StructNameLowerFirst }}Service) Update{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.Update{{ .StructName }}(ctx, {{ .StructNameLowerFirst }})
}

func (s *{{ .StructNameLowerFirst }}Service) Delete{{ .StructName }}(ctx context.Context, id string) error {
	return s.{{ .StructNameLowerFirst }}Repository.Delete{{ .StructName }}(ctx, id)
}

func (s *{{ .StructNameLowerFirst }}Service) List{{ .StructName }}s(ctx context.Context, filter map[string]interface{}) ([]*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.List{{ .StructName }}s(ctx, filter)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchGet{{ .StructName }}s(ctx context.Context, ids []string) ([]*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.BatchGet{{ .StructName }}s(ctx, ids)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchCreate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.BatchCreate{{ .StructName }}s(ctx, {{ .StructNameLowerFirst }}s)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchUpdate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.BatchUpdate{{ .StructName }}s(ctx, {{ .StructNameLowerFirst }}s)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchDelete{{ .StructName }}s(ctx context.Context, ids []string) error {
	return s.{{ .StructNameLowerFirst }}Repository.BatchDelete{{ .StructName }}s(ctx, ids)
}

