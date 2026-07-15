package service

import (
	"errors"
	"regexp"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

type PlanService struct {
	repo repository.PlanRepository
}

func NewPlanService(repo repository.PlanRepository) *PlanService {
	return &PlanService{repo: repo}
}

func validatePlan(name, slug string, priceCents int64) error {
	if name == "" {
		return errors.New("plan name is required")
	}
	if slug == "" || !slugPattern.MatchString(slug) {
		return errors.New("invalid slug: use lowercase letters, numbers and hyphens")
	}
	if priceCents < 0 {
		return errors.New("price_cents must not be negative")
	}
	return nil
}

func (s *PlanService) Create(p *entity.Plan) error {
	if err := validatePlan(p.Name, p.Slug, p.PriceCents); err != nil {
		return err
	}
	if existing, err := s.repo.GetBySlug(p.Slug); err == nil && existing.ID != "" {
		return errors.New("slug already in use")
	}

	p.ID = uuid.New().String()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return s.repo.Create(p)
}

func (s *PlanService) GetByID(id string) (*entity.Plan, error) {
	return s.repo.GetByID(id)
}

func (s *PlanService) List() ([]*entity.Plan, error) {
	return s.repo.List()
}

// ListActive returns only active plans, for public/unauthenticated display.
func (s *PlanService) ListActive() ([]*entity.Plan, error) {
	all, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	active := make([]*entity.Plan, 0, len(all))
	for _, p := range all {
		if p.Active {
			active = append(active, p)
		}
	}
	return active, nil
}

func (s *PlanService) Update(p *entity.Plan) error {
	if err := validatePlan(p.Name, p.Slug, p.PriceCents); err != nil {
		return err
	}
	if existing, err := s.repo.GetBySlug(p.Slug); err == nil && existing.ID != "" && existing.ID != p.ID {
		return errors.New("slug already in use")
	}
	p.UpdatedAt = time.Now()
	return s.repo.Update(p)
}

func (s *PlanService) Delete(id string) error {
	return s.repo.Delete(id)
}
