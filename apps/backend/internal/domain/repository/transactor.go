package repository

type Transactor interface {
	RunInTx(fn func(repos TransactionProvider) error) error
}

type TransactionProvider interface {
	Farm() FarmRepository
	Plot() PlotRepository
	Operation() OperationRepository
	Harvest() HarvestRepository
	HarvestProduction() HarvestProductionRepository
	Indicator() IndicatorRepository
	User() UserRepository
	Tenant() TenantRepository
	AgriculturalProduct() AgriculturalProductRepository
}
