package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// SignupService handles the public, unauthenticated self-service
// registration flow: a new "usuário principal" (proprietario) created under
// a single fixed platform organization (defaultOrgSlug) — the request never
// carries an organization_id, so a public caller can never target an
// arbitrary tenant. Farm registration happens later, once authenticated.
type SignupService struct {
	userRepo         repository.UserRepository
	organizationRepo repository.OrganizationRepository
	planRepo         repository.PlanRepository
	roleSvc          *RoleService
	defaultOrgSlug   string
}

func NewSignupService(userRepo repository.UserRepository, organizationRepo repository.OrganizationRepository, planRepo repository.PlanRepository, roleSvc *RoleService, defaultOrgSlug string) *SignupService {
	return &SignupService{
		userRepo:         userRepo,
		organizationRepo: organizationRepo,
		planRepo:         planRepo,
		roleSvc:          roleSvc,
		defaultOrgSlug:   defaultOrgSlug,
	}
}

type RegisterPrincipalInput struct {
	Name     string
	Email    string
	Password string
	PlanSlug string
}

// RegisterPrincipal creates the user every public signup produces. The new
// user is always an independent principal (ManagedByUserID nil, role
// "proprietario") — unlike UserHandler.CreateForOrg, there is no "first user
// of the org becomes admin, rest become sub-users" rule here, since every
// public signup is meant to be its own independent owner sharing the same
// default platform organization. Farm registration is a separate,
// authenticated step (POST /api/v1/{organization_id}/farms).
func (s *SignupService) RegisterPrincipal(input RegisterPrincipalInput) (*entity.User, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("name, email and password are required")
	}

	org, err := s.organizationRepo.GetBySlug(s.defaultOrgSlug)
	if err != nil {
		return nil, errors.New("default organization not configured")
	}

	role, err := s.roleSvc.FindByKey(entity.RoleKeyProprietario)
	if err != nil {
		return nil, errors.New("proprietario role not found")
	}

	var planID *string
	if input.PlanSlug != "" {
		plan, err := s.planRepo.GetBySlug(input.PlanSlug)
		if err != nil {
			return nil, errors.New("plan not found")
		}
		planID = &plan.ID
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	now := time.Now()
	user := &entity.User{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		Name:           input.Name,
		Email:          input.Email,
		PasswordHash:   string(hash),
		RoleID:         role.ID,
		PlanID:         planID,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
