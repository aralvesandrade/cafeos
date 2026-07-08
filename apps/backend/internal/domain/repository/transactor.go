package repository

type Transactor interface {
	RunInTx(fn func(repos TransactionProvider) error) error
}

type TransactionProvider interface {
	Farm() FarmRepository
	Producer() ProducerRepository
	Plot() PlotRepository
	Operation() OperationRepository
	Harvest() HarvestRepository
	HarvestProduction() HarvestProductionRepository
	Indicator() IndicatorRepository
	User() UserRepository
	Tenant() TenantRepository
	AgriculturalProduct() AgriculturalProductRepository
	Maintenance() MaintenanceRepository
	WorkShift() WorkShiftRepository
	Financial() FinancialRepository
	CostAllocation() CostAllocationRepository
	CostCenter() CostCenterRepository
}
