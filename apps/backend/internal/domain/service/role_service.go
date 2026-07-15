package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

var (
	ErrRoleIsSystem  = errors.New("system roles cannot be edited or deleted")
	ErrRoleInUse     = errors.New("role is assigned to at least one user")
	ErrRoleKeyExists = errors.New("a role with this key already exists in the organization")
	ErrRoleNotFound  = errors.New("role not found")
	ErrForbiddenOrg  = errors.New("role does not belong to this organization")
)

type RoleService struct {
	repo     repository.RoleRepository
	userRepo repository.UserRepository
}

func NewRoleService(repo repository.RoleRepository, userRepo repository.UserRepository) *RoleService {
	return &RoleService{repo: repo, userRepo: userRepo}
}

// SeedSystemRolesIfMissing creates the two global system roles
// (platform_owner, organization_admin) once, at API boot. Safe to call on
// every boot.
func (s *RoleService) SeedSystemRolesIfMissing() error {
	existing, err := s.repo.ListForOrganization("")
	if err != nil {
		return err
	}
	haveKey := make(map[string]bool, len(existing))
	for _, r := range existing {
		if r.OrganizationID == nil {
			haveKey[r.Key] = true
		}
	}

	names := map[string]string{
		entity.SystemRolePlatformOwner:     "Platform Owner",
		entity.SystemRoleOrganizationAdmin: "Admin da Organização",
	}
	for _, key := range entity.SystemRoleKeys {
		if haveKey[key] {
			continue
		}
		role := &entity.Role{
			ID:        uuid.New().String(),
			Key:       key,
			Name:      names[key],
			IsSystem:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.repo.Create(role); err != nil {
			return err
		}
	}
	return nil
}

// SeedDefaultsIfMissing backfills the starter-kit roles (proprietario,
// gerente_agricola, ...) for an organization that has none yet — called at
// organization creation and at API boot for pre-existing organizations.
func (s *RoleService) SeedDefaultsIfMissing(organizationID string) error {
	count, err := s.repo.CountByOrganization(organizationID)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	for _, key := range entity.DefaultOrgRoleKeys {
		role := &entity.Role{
			ID:             uuid.New().String(),
			OrganizationID: &organizationID,
			Key:            key,
			Name:           entity.DefaultOrgRoleNames[key],
			IsSystem:       false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.repo.Create(role); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoleService) ListForOrganization(organizationID string) ([]*entity.Role, error) {
	return s.repo.ListForOrganization(organizationID)
}

func (s *RoleService) GetByID(id string) (*entity.Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// FindByKey resolves a role key within an organization (org-scoped role or
// global system role) — used to translate the legacy "role" string on user
// create/update requests into a RoleID.
func (s *RoleService) FindByKey(organizationID, key string) (*entity.Role, error) {
	role, err := s.repo.FindByKey(organizationID, key)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *RoleService) Create(organizationID, key, name string) (*entity.Role, error) {
	if key == "" || name == "" {
		return nil, errors.New("key and name are required")
	}
	if existing, err := s.repo.FindByKey(organizationID, key); err == nil && existing != nil {
		return nil, ErrRoleKeyExists
	}
	role := &entity.Role{
		ID:             uuid.New().String(),
		OrganizationID: &organizationID,
		Key:            key,
		Name:           name,
		IsSystem:       false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Update(organizationID, id, name string) (*entity.Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	if role.IsSystem {
		return nil, ErrRoleIsSystem
	}
	if role.OrganizationID == nil || *role.OrganizationID != organizationID {
		return nil, ErrForbiddenOrg
	}
	if name != "" {
		role.Name = name
	}
	role.UpdatedAt = time.Now()
	if err := s.repo.Update(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Delete(organizationID, id string) error {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return ErrRoleNotFound
	}
	if role.IsSystem {
		return ErrRoleIsSystem
	}
	if role.OrganizationID == nil || *role.OrganizationID != organizationID {
		return ErrForbiddenOrg
	}
	count, err := s.userRepo.CountByRole(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleInUse
	}
	return s.repo.Delete(id)
}
