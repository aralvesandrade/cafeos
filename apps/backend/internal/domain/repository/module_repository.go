package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type ModuleRepository interface {
	List() ([]*entity.Module, error)
	GetByKey(key entity.ModuleKey) (*entity.Module, error)
	Upsert(module *entity.Module) error
}
