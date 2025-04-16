package repository

import (
    "context"
	"{{ .ProjectName }}/internal/model"
)

type {{ .StructName }}Repository interface {
	Create{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error
	Get{{ .StructName }}(ctx context.Context, id string) (*model.{{ .StructName }}, error)
	Update{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error
	Delete{{ .StructName }}(ctx context.Context, id string) error
	List{{ .StructName }}s(ctx context.Context, filter map[string]interface{}) ([]*model.{{ .StructName }}, error)
	BatchGet{{ .StructName }}s(ctx context.Context, ids []string) ([]*model.{{ .StructName }}, error)
	BatchCreate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error
	BatchUpdate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error
	BatchDelete{{ .StructName }}s(ctx context.Context, ids []string) error
}

// @wire:Repository
func New{{ .StructName }}Repository(
	repository *Repository,
) {{ .StructName }}Repository {
	return &{{ .StructNameLowerFirst }}Repository{
		Repository: repository,
	}
}

type {{ .StructNameLowerFirst }}Repository struct {
	*Repository
}

func (r *{{ .StructNameLowerFirst }}Repository) Create{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error {
	if err := r.Repository.db.WithContext(ctx).Create({{ .StructNameLowerFirst }}).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .StructNameLowerFirst }}Repository) Get{{ .StructName }}(ctx context.Context, id string) (*model.{{ .StructName }}, error) {
	var {{ .StructNameLowerFirst }} model.{{ .StructName }}
	if err := r.Repository.db.WithContext(ctx).Where("id = ?", id).First(&{{ .StructNameLowerFirst }}).Error; err != nil {
		return nil, err
	}
	return &{{ .StructNameLowerFirst }}, nil
}

func (r *{{ .StructNameLowerFirst }}Repository) Update{{ .StructName }}(ctx context.Context, {{ .StructNameLowerFirst }} *model.{{ .StructName }}) error {
	if err := r.Repository.db.WithContext(ctx).Save({{ .StructNameLowerFirst }}).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .StructNameLowerFirst }}Repository) Delete{{ .StructName }}(ctx context.Context, id string) error {
	if err := r.Repository.db.WithContext(ctx).Where("id = ?", id).Delete(&model.{{ .StructName }}{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *{{ .StructNameLowerFirst }}Repository) List{{ .StructName }}s(ctx context.Context, filter map[string]interface{}) ([]*model.{{ .StructName }}, error) {
	var {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}
	query := r.Repository.db.WithContext(ctx).Model(&model.{{ .StructName }}{})
	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}
	if err := query.Find(&{{ .StructNameLowerFirst }}s).Error; err != nil {
		return nil, err
	}
	return {{ .StructNameLowerFirst }}s, nil
}

func (r *{{ .StructNameLowerFirst }}Repository) BatchGet{{ .StructName }}s(ctx context.Context, ids []string) ([]*model.{{ .StructName }}, error) {
	var {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}
	if err := r.Repository.db.WithContext(ctx).Where("id IN ?", ids).Find(&{{ .StructNameLowerFirst }}s).Error; err != nil {
		return nil, err
	}
	return {{ .StructNameLowerFirst }}s, nil
}

func (r *{{ .StructNameLowerFirst }}Repository) BatchCreate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error {
	return r.Repository.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches({{ .StructNameLowerFirst }}s, 100).Error
	})
}

func (r *{{ .StructNameLowerFirst }}Repository) BatchUpdate{{ .StructName }}s(ctx context.Context, {{ .StructNameLowerFirst }}s []*model.{{ .StructName }}) error {
	return r.Repository.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range {{ .StructNameLowerFirst }}s {
			if err := tx.Save(item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *{{ .StructNameLowerFirst }}Repository) BatchDelete{{ .StructName }}s(ctx context.Context, ids []string) error {
	return r.Repository.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.{{ .StructName }}{}).Error
}
