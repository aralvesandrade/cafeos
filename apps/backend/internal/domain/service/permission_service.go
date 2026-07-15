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
// new organization, keyed by role key (see entity.SystemRoleKeys and
// entity.DefaultOrgRoleKeys). Organization admins can then adjust it via the
// Permissions screen — this is only ever used as a seed value, never
// consulted directly at request time.
var defaultMatrix = map[entity.ModuleKey]map[string]entity.AccessLevel{
	// Dashboard also gates GET/PUT /alerts — write lets a role resolve or
	// dismiss an alert, not edit the dashboard itself.
	entity.ModuleDashboard: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessWrite,
		"gerente_agricola":                 entity.AccessWrite,
		"engenheiro_agronomo":              entity.AccessWrite,
		"tecnico_agricola":                 entity.AccessWrite,
		"operador_campo":                   entity.AccessWrite,
		"financeiro":                       entity.AccessWrite,
		"consultor_externo":                entity.AccessRead,
		"auditor":                          entity.AccessRead,
	},
	// Only organization_admin (besides platform_owner) registers farms by
	// default — everyone else can view but not create/edit/delete. Also
	// gates the farm↔user links (Farm.Producers), since those are set
	// through the same create/update endpoints.
	entity.ModuleFarms: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessRead,
		"gerente_agricola":                 entity.AccessRead,
		"engenheiro_agronomo":              entity.AccessRead,
		"tecnico_agricola":                 entity.AccessRead,
		"operador_campo":                   entity.AccessRead,
		"financeiro":                       entity.AccessRead,
		"consultor_externo":                entity.AccessRead,
		"auditor":                          entity.AccessRead,
	},
	entity.ModuleOperations: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessWrite,
		"gerente_agricola":                 entity.AccessWrite,
		"engenheiro_agronomo":              entity.AccessWrite,
		"tecnico_agricola":                 entity.AccessWrite,
		"operador_campo":                   entity.AccessWrite,
		"financeiro":                       entity.AccessRead,
		"consultor_externo":                entity.AccessRead,
		"auditor":                          entity.AccessRead,
	},
	entity.ModuleHarvests: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessWrite,
		"gerente_agricola":                 entity.AccessWrite,
		"engenheiro_agronomo":              entity.AccessWrite,
		"tecnico_agricola":                 entity.AccessRead,
		"operador_campo":                   entity.AccessRead,
		"financeiro":                       entity.AccessRead,
		"consultor_externo":                entity.AccessRead,
		"auditor":                          entity.AccessRead,
	},
	entity.ModuleResources: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessWrite,
		"gerente_agricola":                 entity.AccessWrite,
		"engenheiro_agronomo":              entity.AccessRead,
		"tecnico_agricola":                 entity.AccessWrite,
		"operador_campo":                   entity.AccessRead,
		"financeiro":                       entity.AccessRead,
		"consultor_externo":                entity.AccessRead,
		"auditor":                          entity.AccessRead,
	},
	entity.ModuleFinancial: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessWrite,
		"gerente_agricola":                 entity.AccessRead,
		"engenheiro_agronomo":              entity.AccessNone,
		"tecnico_agricola":                 entity.AccessNone,
		"operador_campo":                   entity.AccessNone,
		"financeiro":                       entity.AccessWrite,
		"consultor_externo":                entity.AccessRead,
		"auditor":                          entity.AccessRead,
	},
	// Users covers the org-scoped team roster (POST/GET/PUT/DELETE
	// /{organization_id}/users) — distinct from the cross-tenant
	// /admin/users screen, which stays platform_owner-only regardless.
	// Only organization_admin registers new users by default; proprietario
	// can still view the roster but not create/edit/delete.
	entity.ModuleUsers: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessRead,
		"gerente_agricola":                 entity.AccessNone,
		"engenheiro_agronomo":              entity.AccessNone,
		"tecnico_agricola":                 entity.AccessNone,
		"operador_campo":                   entity.AccessNone,
		"financeiro":                       entity.AccessNone,
		"consultor_externo":                entity.AccessNone,
		"auditor":                          entity.AccessNone,
	},
	entity.ModulePermissions: {
		entity.SystemRolePlatformOwner:     entity.AccessWrite,
		entity.SystemRoleOrganizationAdmin: entity.AccessWrite,
		"proprietario":                     entity.AccessWrite,
		"gerente_agricola":                 entity.AccessNone,
		"engenheiro_agronomo":              entity.AccessNone,
		"tecnico_agricola":                 entity.AccessNone,
		"operador_campo":                   entity.AccessNone,
		"financeiro":                       entity.AccessNone,
		"consultor_externo":                entity.AccessNone,
		"auditor":                          entity.AccessNone,
	},
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
	repo     repository.PermissionRepository
	roleRepo repository.RoleRepository

	mu    sync.RWMutex
	cache map[string]map[string]map[entity.ModuleKey]entity.AccessLevel
}

func NewPermissionService(repo repository.PermissionRepository, roleRepo repository.RoleRepository) *PermissionService {
	return &PermissionService{
		repo:     repo,
		roleRepo: roleRepo,
		cache:    make(map[string]map[string]map[entity.ModuleKey]entity.AccessLevel),
	}
}

// SeedDefaults overwrites every cell with the default matrix — callers that
// must not clobber manual edits should use SeedDefaultsIfMissing instead.
// It assumes the organization's roles have already been seeded (see
// RoleService.SeedDefaultsIfMissing), since it needs each role's ID.
func (s *PermissionService) SeedDefaults(organizationID string) error {
	roles, err := s.roleRepo.List()
	if err != nil {
		return err
	}
	roleIDByKey := make(map[string]string, len(roles))
	for _, role := range roles {
		roleIDByKey[role.Key] = role.ID
	}

	for module, roleAccess := range defaultMatrix {
		for roleKey, access := range roleAccess {
			roleID, ok := roleIDByKey[roleKey]
			if !ok {
				continue
			}
			permission := &entity.RolePermission{
				ID:             uuid.New().String(),
				OrganizationID: organizationID,
				RoleID:         roleID,
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

func (s *PermissionService) Update(organizationID, roleID string, module entity.ModuleKey, access entity.AccessLevel) error {
	if _, err := s.roleRepo.GetByID(roleID); err != nil {
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
		RoleID:         roleID,
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

func validModule(module entity.ModuleKey) bool {
	return slices.Contains(entity.AllModules, module)
}

// GetAccess falls back to defaultMatrix for orgs/cells not yet seeded, so a
// missing row never silently locks a role out. roleKey is the role's Key
// (what's carried in the JWT / request context), not its ID.
func (s *PermissionService) GetAccess(organizationID, roleKey string, module entity.ModuleKey) (entity.AccessLevel, error) {
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

	if roleAccess, ok := orgCache[roleKey]; ok {
		if access, ok := roleAccess[module]; ok {
			return access, nil
		}
	}

	if roleDefaults, ok := defaultMatrix[module]; ok {
		if access, ok := roleDefaults[roleKey]; ok {
			return access, nil
		}
	}
	return entity.AccessNone, nil
}

func (s *PermissionService) loadOrganization(organizationID string) (map[string]map[entity.ModuleKey]entity.AccessLevel, error) {
	rows, err := s.repo.ListByOrganization(organizationID)
	if err != nil {
		return nil, err
	}
	roles, err := s.roleRepo.List()
	if err != nil {
		return nil, err
	}
	roleKeyByID := make(map[string]string, len(roles))
	for _, role := range roles {
		roleKeyByID[role.ID] = role.Key
	}

	orgCache := make(map[string]map[entity.ModuleKey]entity.AccessLevel)
	for _, row := range rows {
		roleKey, ok := roleKeyByID[row.RoleID]
		if !ok {
			continue
		}
		if orgCache[roleKey] == nil {
			orgCache[roleKey] = make(map[entity.ModuleKey]entity.AccessLevel)
		}
		orgCache[roleKey][row.Module] = row.Access
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
