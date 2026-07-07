package postgres

import (
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	infraRepo "github.com/aralvesandrade/cafeos/internal/infra/db/repository"
	"gorm.io/gorm"
)

type transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) repository.Transactor {
	return &transactor{db: db}
}

func (t *transactor) RunInTx(fn func(repos repository.TransactionProvider) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return fn(&txProvider{tx: tx})
	})
}

type txProvider struct {
	tx *gorm.DB
}

func (p *txProvider) Farm() repository.FarmRepository {
	return infraRepo.NewFarmRepository(p.tx)
}

func (p *txProvider) Plot() repository.PlotRepository {
	return infraRepo.NewPlotRepository(p.tx)
}

func (p *txProvider) Operation() repository.OperationRepository {
	return infraRepo.NewOperationRepository(p.tx)
}

func (p *txProvider) Harvest() repository.HarvestRepository {
	return infraRepo.NewHarvestRepository(p.tx)
}

func (p *txProvider) HarvestProduction() repository.HarvestProductionRepository {
	return infraRepo.NewHarvestProductionRepository(p.tx)
}

func (p *txProvider) Indicator() repository.IndicatorRepository {
	return infraRepo.NewIndicatorRepository(p.tx)
}

func (p *txProvider) User() repository.UserRepository {
	return infraRepo.NewUserRepository(p.tx)
}

func (p *txProvider) Tenant() repository.TenantRepository {
	return infraRepo.NewTenantRepository(p.tx)
}

func (p *txProvider) AgriculturalProduct() repository.AgriculturalProductRepository {
	return infraRepo.NewAgriculturalProductRepository(p.tx)
}

func (p *txProvider) Maintenance() repository.MaintenanceRepository {
	return infraRepo.NewMaintenanceRepository(p.tx)
}

func (p *txProvider) WorkShift() repository.WorkShiftRepository {
	return infraRepo.NewWorkShiftRepository(p.tx)
}

func (p *txProvider) Financial() repository.FinancialRepository {
	return infraRepo.NewFinancialRepository(p.tx)
}

func (p *txProvider) CostAllocation() repository.CostAllocationRepository {
	return infraRepo.NewCostAllocationRepository(p.tx)
}
