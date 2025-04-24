package service

import (
	"{{ .ProjectName }}/internal/model"
	"{{ .ProjectName }}/internal/repository"
)

type {{ .StructName }}Service interface {
	Create{{ .StructName }}( {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error
	Get{{ .StructName }}( id string) (*model.{{ .StructName }}, error)
	Update{{ .StructName }}( {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error
	Delete{{ .StructName }}( id string) error
	List{{ .StructName }}s( filter map[string]interface{}) ([]*model.{{ .StructName }}, error)
	BatchGet{{ .StructName }}s( ids []string) ([]*model.{{ .StructName }}, error)
	BatchCreate{{ .StructName }}s( {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error
	BatchUpdate{{ .StructName }}s( {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error
	BatchDelete{{ .StructName }}s( ids []string) error
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

func (s *{{ .StructNameLowerFirst }}Service) Create{{ .StructName }}( {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.Create{{ .StructName }}( {{ .StructNameLowerFirst }})
}

func (s *{{ .StructNameLowerFirst }}Service) Get{{ .StructName }}( id string) (*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.Get{{ .StructName }}( id)
}

func (s *{{ .StructNameLowerFirst }}Service) Update{{ .StructName }}( {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.Update{{ .StructName }}( {{ .StructNameLowerFirst }})
}

func (s *{{ .StructNameLowerFirst }}Service) Delete{{ .StructName }}( id string) error {
	return s.{{ .StructNameLowerFirst }}Repository.Delete{{ .StructName }}( id)
}

func (s *{{ .StructNameLowerFirst }}Service) List{{ .StructName }}s( filter map[string]interface{}) ([]*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.List{{ .StructName }}s( filter)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchGet{{ .StructName }}s( ids []string) ([]*model.{{ .StructName }}, error) {
	return s.{{ .StructNameLowerFirst }}Repository.BatchGet{{ .StructName }}s( ids)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchCreate{{ .StructName }}s( {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.BatchCreate{{ .StructName }}s( {{ .StructNameLowerFirst }}s)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchUpdate{{ .StructName }}s( {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error {
	return s.{{ .StructNameLowerFirst }}Repository.BatchUpdate{{ .StructName }}s( {{ .StructNameLowerFirst }}s)
}

func (s *{{ .StructNameLowerFirst }}Service) BatchDelete{{ .StructName }}s( ids []string) error {
	return s.{{ .StructNameLowerFirst }}Repository.BatchDelete{{ .StructName }}s( ids)
}

