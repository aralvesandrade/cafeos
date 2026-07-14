package service

import (
	"errors"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type LaborService struct {
	teamRepo   repository.TeamRepository
	workerRepo repository.WorkerRepository
	shiftRepo  repository.WorkShiftRepository
}

func NewLaborService(teamRepo repository.TeamRepository, workerRepo repository.WorkerRepository, shiftRepo repository.WorkShiftRepository) *LaborService {
	return &LaborService{teamRepo: teamRepo, workerRepo: workerRepo, shiftRepo: shiftRepo}
}

// Team operations
func (s *LaborService) CreateTeam(organizationID, name, leader, description string) (*entity.Team, error) {
	if name == "" {
		return nil, errors.New("team name is required")
	}
	t := &entity.Team{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Name:           name,
		Leader:         leader,
		Description:    description,
		CreatedAt:      time.Now(),
	}
	if err := s.teamRepo.Create(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *LaborService) ListTeams(organizationID string) ([]*entity.Team, error) {
	return s.teamRepo.ListByOrganization(organizationID)
}

func (s *LaborService) GetTeamByID(id string) (*entity.Team, error) {
	return s.teamRepo.GetByID(id)
}

func (s *LaborService) UpdateTeam(t *entity.Team) error {
	return s.teamRepo.Update(t)
}

func (s *LaborService) DeleteTeam(id string) error {
	return s.teamRepo.Delete(id)
}

// Worker operations
func (s *LaborService) CreateWorker(organizationID, teamID, name, role, phone string, hourlyRate float64) (*entity.Worker, error) {
	if name == "" {
		return nil, errors.New("worker name is required")
	}
	w := &entity.Worker{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		TeamID:         teamID,
		Name:           name,
		Role:           role,
		Phone:          phone,
		HourlyRate:     hourlyRate,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.workerRepo.Create(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *LaborService) ListWorkers(organizationID string) ([]*entity.Worker, error) {
	return s.workerRepo.ListByOrganization(organizationID)
}

func (s *LaborService) GetWorkerByID(id string) (*entity.Worker, error) {
	return s.workerRepo.GetByID(id)
}

func (s *LaborService) UpdateWorker(w *entity.Worker) error {
	w.UpdatedAt = time.Now()
	return s.workerRepo.Update(w)
}

func (s *LaborService) DeleteWorker(id string) error {
	return s.workerRepo.Delete(id)
}

// WorkShift operations
func (s *LaborService) CreateWorkShift(organizationID, workerID, operationID, notes string, hours, cost float64, date time.Time) (*entity.WorkShift, error) {
	if workerID == "" {
		return nil, errors.New("worker is required")
	}
	if hours <= 0 {
		return nil, errors.New("hours must be greater than zero")
	}
	var opRef *string
	if operationID != "" {
		opRef = &operationID
	}
	ws := &entity.WorkShift{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		WorkerID:       workerID,
		OperationID:    opRef,
		Date:           date,
		Hours:          hours,
		Cost:           cost,
		Notes:          notes,
		CreatedAt:      time.Now(),
	}
	if err := s.shiftRepo.Create(ws); err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *LaborService) ListWorkShifts(organizationID string) ([]*entity.WorkShift, error) {
	return s.shiftRepo.ListByOrganization(organizationID)
}

func (s *LaborService) DeleteWorkShift(id string) error {
	return s.shiftRepo.Delete(id)
}
