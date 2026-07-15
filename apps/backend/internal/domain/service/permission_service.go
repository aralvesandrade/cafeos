package service

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

// defaultMatrix is the starting access-level configuration seeded for every
// new organization. Organization admins can then adjust it via the
// Permissions screen — this is only ever used as a seed value, never
// consulted directly at request time.
var defaultMatrix = map[entity.Module]map[entity.UserRole]entity.AccessLevel{
	// Dashboard also gates GET/PUT /alerts — write lets a role resolve or
	// dismiss an alert, not edit the dashboard itself.
	entity.ModuleDashboard: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessWrite,
		entity.RoleEngenheiro:        entity.AccessWrite,
		entity.RoleTecnico:           entity.AccessWrite,
		entity.RoleOperador:          entity.AccessWrite,
		entity.RoleFinanceiro:        entity.AccessWrite,
		entity.RoleConsultor:         entity.AccessRead,
		entity.RoleAuditor:           entity.AccessRead,
	},
	entity.ModuleFarms: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessWrite,
		entity.RoleEngenheiro:        entity.AccessWrite,
		entity.RoleTecnico:           entity.AccessWrite,
		entity.RoleOperador:          entity.AccessRead,
		entity.RoleFinanceiro:        entity.AccessRead,
		entity.RoleConsultor:         entity.AccessRead,
		entity.RoleAuditor:           entity.AccessRead,
	},
	entity.ModuleOperations: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessWrite,
		entity.RoleEngenheiro:        entity.AccessWrite,
		entity.RoleTecnico:           entity.AccessWrite,
		entity.RoleOperador:          entity.AccessWrite,
		entity.RoleFinanceiro:        entity.AccessRead,
		entity.RoleConsultor:         entity.AccessRead,
		entity.RoleAuditor:           entity.AccessRead,
	},
	entity.ModuleHarvests: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessWrite,
		entity.RoleEngenheiro:        entity.AccessWrite,
		entity.RoleTecnico:           entity.AccessRead,
		entity.RoleOperador:          entity.AccessRead,
		entity.RoleFinanceiro:        entity.AccessRead,
		entity.RoleConsultor:         entity.AccessRead,
		entity.RoleAuditor:           entity.AccessRead,
	},
	entity.ModuleResources: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessWrite,
		entity.RoleEngenheiro:        entity.AccessRead,
		entity.RoleTecnico:           entity.AccessWrite,
		entity.RoleOperador:          entity.AccessRead,
		entity.RoleFinanceiro:        entity.AccessRead,
		entity.RoleConsultor:         entity.AccessRead,
		entity.RoleAuditor:           entity.AccessRead,
	},
	entity.ModuleFinancial: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessRead,
		entity.RoleEngenheiro:        entity.AccessNone,
		entity.RoleTecnico:           entity.AccessNone,
		entity.RoleOperador:          entity.AccessNone,
		entity.RoleFinanceiro:        entity.AccessWrite,
		entity.RoleConsultor:         entity.AccessRead,
		entity.RoleAuditor:           entity.AccessRead,
	},
	// Users covers the org-scoped team roster (POST/GET/PUT/DELETE
	// /{organization_id}/users) — distinct from the cross-tenant
	// /admin/users screen, which stays platform_owner-only regardless.
	entity.ModuleUsers: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessNone,
		entity.RoleEngenheiro:        entity.AccessNone,
		entity.RoleTecnico:           entity.AccessNone,
		entity.RoleOperador:          entity.AccessNone,
		entity.RoleFinanceiro:        entity.AccessNone,
		entity.RoleConsultor:         entity.AccessNone,
		entity.RoleAuditor:           entity.AccessNone,
	},
	entity.ModulePermissions: {
		entity.RolePlatformOwner:     entity.AccessWrite,
		entity.RoleOrganizationAdmin: entity.AccessWrite,
		entity.RoleProprietario:      entity.AccessWrite,
		entity.RoleGerente:           entity.AccessNone,
		entity.RoleEngenheiro:        entity.AccessNone,
		entity.RoleTecnico:           entity.AccessNone,
		entity.RoleOperador:          entity.AccessNone,
		entity.RoleFinanceiro:        entity.AccessNone,
		entity.RoleConsultor:         entity.AccessNone,
		entity.RoleAuditor:           entity.AccessNone,
	},
}

var allRoles = []entity.UserRole{
	entity.RolePlatformOwner, entity.RoleOrganizationAdmin, entity.RoleProprietario,
	entity.RoleGerente, entity.RoleEngenheiro, entity.RoleTecnico, entity.RoleOperador,
	entity.RoleFinanceiro, entity.RoleConsultor, entity.RoleAuditor,
}

func validRole(role entity.UserRole) bool {
	return slices.Contains(allRoles, role)
}

func validModule(module entity.Module) bool {
	return slices.Contains(entity.AllModules, module)
}

func validAccess(access entity.AccessLevel) bool {
	switch access {
	case entity.AccessNone, entity.AccessRead, entity.AccessWrite:
		return true
	default:
		return false
	}
}

type PermissionService struct {
	repo repository.PermissionRepository

	mu    sync.RWMutex
	cache map[string]map[entity.UserRole]map[entity.Module]entity.AccessLevel
}

func NewPermissionService(repo repository.PermissionRepository) *PermissionService {
	return &PermissionService{
		repo:  repo,
		cache: make(map[string]map[entity.UserRole]map[entity.Module]entity.AccessLevel),
	}
}

// SeedDefaults overwrites every cell with the default matrix — callers that
// must not clobber manual edits should use SeedDefaultsIfMissing instead.
func (s *PermissionService) SeedDefaults(organizationID string) error {
	for module, roles := range defaultMatrix {
		for role, access := range roles {
			permission := &entity.RolePermission{
				ID:             uuid.New().String(),
				OrganizationID: organizationID,
				Role:           role,
				Module:         module,
				Access:         access,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			if err := s.repo.Upsert(permission); err != nil {
				return err
			}
		}
	}
	s.invalidate(organizationID)
	return nil
}

// SeedDefaultsIfMissing seeds the default matrix only if the organization has
// no permission rows yet — used to backfill organizations created before
// this feature existed, without overwriting any manual configuration.
func (s *PermissionService) SeedDefaultsIfMissing(organizationID string) error {
	count, err := s.repo.CountByOrganization(organizationID)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.SeedDefaults(organizationID)
}

func (s *PermissionService) ListByOrganization(organizationID string) ([]*entity.RolePermission, error) {
	return s.repo.ListByOrganization(organizationID)
}

func (s *PermissionService) Update(organizationID string, role entity.UserRole, module entity.Module, access entity.AccessLevel) error {
	if !validRole(role) {
		return errors.New("invalid role")
	}
	if !validModule(module) {
		return errors.New("invalid module")
	}
	if !validAccess(access) {
		return errors.New("invalid access level")
	}

	permission := &entity.RolePermission{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Role:           role,
		Module:         module,
		Access:         access,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.Upsert(permission); err != nil {
		return err
	}
	s.invalidate(organizationID)
	return nil
}

// GetAccess falls back to defaultMatrix for orgs/cells not yet seeded, so a
// missing row never silently locks a role out.
func (s *PermissionService) GetAccess(organizationID string, role entity.UserRole, module entity.Module) (entity.AccessLevel, error) {
	s.mu.RLock()
	orgCache, ok := s.cache[organizationID]
	s.mu.RUnlock()

	if !ok {
		var err error
		orgCache, err = s.loadOrganization(organizationID)
		if err != nil {
			return entity.AccessNone, err
		}
	}

	if roleAccess, ok := orgCache[role]; ok {
		if access, ok := roleAccess[module]; ok {
			return access, nil
		}
	}

	if roleDefaults, ok := defaultMatrix[module]; ok {
		if access, ok := roleDefaults[role]; ok {
			return access, nil
		}
	}
	return entity.AccessNone, nil
}

func (s *PermissionService) loadOrganization(organizationID string) (map[entity.UserRole]map[entity.Module]entity.AccessLevel, error) {
	rows, err := s.repo.ListByOrganization(organizationID)
	if err != nil {
		return nil, err
	}

	orgCache := make(map[entity.UserRole]map[entity.Module]entity.AccessLevel)
	for _, row := range rows {
		if orgCache[row.Role] == nil {
			orgCache[row.Role] = make(map[entity.Module]entity.AccessLevel)
		}
		orgCache[row.Role][row.Module] = row.Access
	}

	s.mu.Lock()
	s.cache[organizationID] = orgCache
	s.mu.Unlock()

	return orgCache, nil
}

func (s *PermissionService) invalidate(organizationID string) {
	s.mu.Lock()
	delete(s.cache, organizationID)
	s.mu.Unlock()
}
