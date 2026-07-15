package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
)

type ModuleService struct {
	repo repository.ModuleRepository
}

func NewModuleService(repo repository.ModuleRepository) *ModuleService {
	return &ModuleService{repo: repo}
}

// SeedDefaultsIfMissing populates the modules table with the fixed set of
// application modules on first boot. Safe to call on every boot — Upsert is
// idempotent and never overwrites a name/order an admin already changed,
// since callers only invoke this when the table is empty.
func (s *ModuleService) SeedDefaultsIfMissing() error {
	existing, err := s.repo.List()
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}
	for _, key := range entity.AllModules {
		module := &entity.Module{
			Key:       key,
			Name:      entity.DefaultModuleNames[key],
			Order:     entity.DefaultModuleOrder[key],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.repo.Upsert(module); err != nil {
			return err
		}
	}
	return nil
}

func (s *ModuleService) List() ([]*entity.Module, error) {
	return s.repo.List()
}

// UpdateMeta renames/reorders a module. The key itself is fixed — it is
// wired directly into route declarations in router.go.
func (s *ModuleService) UpdateMeta(key entity.ModuleKey, name string, order int) (*entity.Module, error) {
	module, err := s.repo.GetByKey(key)
	if err != nil {
		return nil, errors.New("module not found")
	}
	if name != "" {
		module.Name = name
	}
	module.Order = order
	module.UpdatedAt = time.Now()
	if err := s.repo.Upsert(module); err != nil {
		return nil, err
	}
	return module, nil
}
