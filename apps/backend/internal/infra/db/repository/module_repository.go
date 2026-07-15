package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{db: db}
}

func (r *ModuleRepository) List() ([]*entity.Module, error) {
	var modules []*entity.Module
	err := r.db.Order("\"order\" ASC").Find(&modules).Error
	return modules, err
}

func (r *ModuleRepository) GetByKey(key entity.ModuleKey) (*entity.Module, error) {
	var module entity.Module
	err := r.db.Where("key = ?", key).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (r *ModuleRepository) Upsert(module *entity.Module) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "order", "updated_at"}),
	}).Create(module).Error
}
