package repository

import (
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) WithTx(tx *gorm.DB) *TeamRepository {
	return &TeamRepository{db: tx}
}

func (r *TeamRepository) Create(t *entity.Team) error {
	return r.db.Create(t).Error
}

func (r *TeamRepository) GetByID(id string) (*entity.Team, error) {
	var t entity.Team
	err := r.db.First(&t, "id = ?", id).Error
	return &t, err
}

func (r *TeamRepository) ListByOrganization(organizationID string) ([]*entity.Team, error) {
	var items []*entity.Team
	err := r.db.Where("organization_id = ?", organizationID).Order("name").Find(&items).Error
	return items, err
}

func (r *TeamRepository) Update(t *entity.Team) error {
	return r.db.Save(t).Error
}

func (r *TeamRepository) Delete(id string) error {
	return r.db.Delete(&entity.Team{}, "id = ?", id).Error
}

type WorkerRepository struct {
	db *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

func (r *WorkerRepository) WithTx(tx *gorm.DB) *WorkerRepository {
	return &WorkerRepository{db: tx}
}

func (r *WorkerRepository) Create(w *entity.Worker) error {
	return r.db.Create(w).Error
}

func (r *WorkerRepository) GetByID(id string) (*entity.Worker, error) {
	var w entity.Worker
	err := r.db.Preload("Team").First(&w, "id = ?", id).Error
	return &w, err
}

func (r *WorkerRepository) ListByOrganization(organizationID string) ([]*entity.Worker, error) {
	var items []*entity.Worker
	err := r.db.Preload("Team").Where("organization_id = ?", organizationID).Order("name").Find(&items).Error
	return items, err
}

func (r *WorkerRepository) ListByTeam(teamID string) ([]*entity.Worker, error) {
	var items []*entity.Worker
	err := r.db.Where("team_id = ?", teamID).Order("name").Find(&items).Error
	return items, err
}

func (r *WorkerRepository) Update(w *entity.Worker) error {
	return r.db.Save(w).Error
}

func (r *WorkerRepository) Delete(id string) error {
	return r.db.Delete(&entity.Worker{}, "id = ?", id).Error
}

type WorkShiftRepository struct {
	db *gorm.DB
}

func NewWorkShiftRepository(db *gorm.DB) *WorkShiftRepository {
	return &WorkShiftRepository{db: db}
}

func (r *WorkShiftRepository) WithTx(tx *gorm.DB) *WorkShiftRepository {
	return &WorkShiftRepository{db: tx}
}

func (r *WorkShiftRepository) Create(ws *entity.WorkShift) error {
	return r.db.Create(ws).Error
}

func (r *WorkShiftRepository) GetByID(id string) (*entity.WorkShift, error) {
	var ws entity.WorkShift
	err := r.db.Preload("Worker").First(&ws, "id = ?", id).Error
	return &ws, err
}

func (r *WorkShiftRepository) ListByOrganization(organizationID string) ([]*entity.WorkShift, error) {
	var items []*entity.WorkShift
	err := r.db.Preload("Worker").Where("organization_id = ?", organizationID).Order("date DESC").Find(&items).Error
	return items, err
}

func (r *WorkShiftRepository) ListByWorker(workerID string) ([]*entity.WorkShift, error) {
	var items []*entity.WorkShift
	err := r.db.Where("worker_id = ?", workerID).Order("date DESC").Find(&items).Error
	return items, err
}

func (r *WorkShiftRepository) Delete(id string) error {
	return r.db.Delete(&entity.WorkShift{}, "id = ?", id).Error
}
