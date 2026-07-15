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
	ErrRoleKeyExists = errors.New("a role with this key already exists")
	ErrRoleNotFound  = errors.New("role not found")
)

// RoleService manages the global roles catalog. Roles are shared by every
// organization — like Module — so nothing here is organization-scoped; what
// varies per organization is access (see PermissionService), not the set of
// available roles.
type RoleService struct {
	repo     repository.RoleRepository
	userRepo repository.UserRepository
}

func NewRoleService(repo repository.RoleRepository, userRepo repository.UserRepository) *RoleService {
	return &RoleService{repo: repo, userRepo: userRepo}
}

// SeedDefaultsIfMissing populates the roles table with the two system roles
// plus the eight starter-kit roles on first boot. Safe to call on every
// boot — a non-empty table means seeding already happened.
func (s *RoleService) SeedDefaultsIfMissing() error {
	count, err := s.repo.Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	names := map[string]string{
		entity.SystemRolePlatformOwner:     "Platform Owner",
		entity.SystemRoleOrganizationAdmin: "Admin da Organização",
	}
	for _, key := range entity.SystemRoleKeys {
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
	for _, key := range entity.DefaultRoleKeys {
		role := &entity.Role{
			ID:        uuid.New().String(),
			Key:       key,
			Name:      entity.DefaultRoleNames[key],
			IsSystem:  false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := s.repo.Create(role); err != nil {
			return err
		}
	}
	return nil
}

func (s *RoleService) List() ([]*entity.Role, error) {
	return s.repo.List()
}

func (s *RoleService) GetByID(id string) (*entity.Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// FindByKey resolves a role key — used to translate the legacy "role"
// string on user create/update requests into a RoleID.
func (s *RoleService) FindByKey(key string) (*entity.Role, error) {
	role, err := s.repo.FindByKey(key)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

func (s *RoleService) Create(key, name string) (*entity.Role, error) {
	if key == "" || name == "" {
		return nil, errors.New("key and name are required")
	}
	if existing, err := s.repo.FindByKey(key); err == nil && existing != nil {
		return nil, ErrRoleKeyExists
	}
	role := &entity.Role{
		ID:        uuid.New().String(),
		Key:       key,
		Name:      name,
		IsSystem:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) Update(id, name string) (*entity.Role, error) {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrRoleNotFound
	}
	if role.IsSystem {
		return nil, ErrRoleIsSystem
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

func (s *RoleService) Delete(id string) error {
	role, err := s.repo.GetByID(id)
	if err != nil {
		return ErrRoleNotFound
	}
	if role.IsSystem {
		return ErrRoleIsSystem
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
