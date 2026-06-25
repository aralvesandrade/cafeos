package repository

import "github.com/aralvesandrade/cafeos/internal/domain/entity"

type TeamRepository interface {
	Create(t *entity.Team) error
	GetByID(id string) (*entity.Team, error)
	ListByTenant(tenantID string) ([]*entity.Team, error)
	Update(t *entity.Team) error
	Delete(id string) error
}

type WorkerRepository interface {
	Create(w *entity.Worker) error
	GetByID(id string) (*entity.Worker, error)
	ListByTenant(tenantID string) ([]*entity.Worker, error)
	ListByTeam(teamID string) ([]*entity.Worker, error)
	Update(w *entity.Worker) error
	Delete(id string) error
}

type WorkShiftRepository interface {
	Create(ws *entity.WorkShift) error
	GetByID(id string) (*entity.WorkShift, error)
	ListByTenant(tenantID string) ([]*entity.WorkShift, error)
	ListByWorker(workerID string) ([]*entity.WorkShift, error)
	Delete(id string) error
}
