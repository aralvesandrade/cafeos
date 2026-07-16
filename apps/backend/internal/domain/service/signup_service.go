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
// registration flow: a new "usuário principal" (proprietario) plus their
// first Farm, both created under a single fixed platform organization
// (defaultOrgSlug) — the request never carries an organization_id, so a
// public caller can never target an arbitrary tenant.
type SignupService struct {
	transactor       repository.Transactor
	organizationRepo repository.OrganizationRepository
	planRepo         repository.PlanRepository
	roleSvc          *RoleService
	defaultOrgSlug   string
}

func NewSignupService(transactor repository.Transactor, organizationRepo repository.OrganizationRepository, planRepo repository.PlanRepository, roleSvc *RoleService, defaultOrgSlug string) *SignupService {
	return &SignupService{
		transactor:       transactor,
		organizationRepo: organizationRepo,
		planRepo:         planRepo,
		roleSvc:          roleSvc,
		defaultOrgSlug:   defaultOrgSlug,
	}
}

type RegisterPrincipalInput struct {
	Name        string
	Email       string
	Password    string
	Phone       string
	FarmName    string
	State       string
	City        string
	MainCrop    string
	TotalAreaHA float64
	PlanSlug    string
}

// RegisterPrincipal creates the user and farm every public signup produces.
// The new user is always an independent principal (ManagedByUserID nil,
// role "proprietario") — unlike UserHandler.CreateForOrg, there is no
// "first user of the org becomes admin, rest become sub-users" rule here,
// since every public signup is meant to be its own independent owner
// sharing the same default platform organization.
func (s *SignupService) RegisterPrincipal(input RegisterPrincipalInput) (*entity.User, *entity.Farm, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, nil, errors.New("name, email and password are required")
	}
	if input.FarmName == "" {
		return nil, nil, errors.New("farm name is required")
	}

	org, err := s.organizationRepo.GetBySlug(s.defaultOrgSlug)
	if err != nil {
		return nil, nil, errors.New("default organization not configured")
	}

	role, err := s.roleSvc.FindByKey(entity.RoleKeyProprietario)
	if err != nil {
		return nil, nil, errors.New("proprietario role not found")
	}

	var planID *string
	if input.PlanSlug != "" {
		plan, err := s.planRepo.GetBySlug(input.PlanSlug)
		if err != nil {
			return nil, nil, errors.New("plan not found")
		}
		planID = &plan.ID
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, errors.New("failed to hash password")
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

	farm := &entity.Farm{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		Name:           input.FarmName,
		Phone:          input.Phone,
		State:          input.State,
		City:           input.City,
		MainCrop:       input.MainCrop,
		TotalAreaHA:    input.TotalAreaHA,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	producer := &entity.Producer{
		ID:             uuid.New().String(),
		OrganizationID: org.ID,
		UserID:         user.ID,
		RoleID:         role.ID,
		Name:           input.Name,
		Phone:          input.Phone,
		Email:          input.Email,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = s.transactor.RunInTx(func(repos repository.TransactionProvider) error {
		if err := repos.User().Create(user); err != nil {
			return err
		}
		if err := repos.Farm().Create(farm); err != nil {
			return err
		}
		producer.FarmID = farm.ID
		return repos.Producer().Create(producer)
	})
	if err != nil {
		return nil, nil, err
	}

	return user, farm, nil
}
