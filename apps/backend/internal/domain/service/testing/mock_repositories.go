package testing

import (
	"errors"
	"sync"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type InMemoryFarmRepo struct {
	mu     sync.RWMutex
	farms  map[string]*entity.Farm
}

func NewInMemoryFarmRepo() *InMemoryFarmRepo {
	return &InMemoryFarmRepo{farms: make(map[string]*entity.Farm)}
}

func (r *InMemoryFarmRepo) Create(f *entity.Farm) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.farms[f.ID] = f
	return nil
}

func (r *InMemoryFarmRepo) GetByID(id string) (*entity.Farm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.farms[id]
	if !ok {
		return nil, errors.New("farm not found")
	}
	return f, nil
}

func (r *InMemoryFarmRepo) ListByTenant(tenantID string) ([]*entity.Farm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Farm
	for _, f := range r.farms {
		if f.TenantID == tenantID {
			result = append(result, f)
		}
	}
	return result, nil
}

func (r *InMemoryFarmRepo) Update(f *entity.Farm) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.farms[f.ID] = f
	return nil
}

func (r *InMemoryFarmRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.farms, id)
	return nil
}

type InMemoryPlotRepo struct {
	mu    sync.RWMutex
	plots map[string]*entity.Plot
}

func NewInMemoryPlotRepo() *InMemoryPlotRepo {
	return &InMemoryPlotRepo{plots: make(map[string]*entity.Plot)}
}

func (r *InMemoryPlotRepo) Create(p *entity.Plot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plots[p.ID] = p
	return nil
}

func (r *InMemoryPlotRepo) GetByID(id string) (*entity.Plot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.plots[id]
	if !ok {
		return nil, errors.New("plot not found")
	}
	return p, nil
}

func (r *InMemoryPlotRepo) ListByFarm(farmID string) ([]*entity.Plot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Plot
	for _, p := range r.plots {
		if p.FarmID == farmID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemoryPlotRepo) ListByTenant(tenantID string) ([]*entity.Plot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Plot
	for _, p := range r.plots {
		if p.TenantID == tenantID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemoryPlotRepo) Update(p *entity.Plot) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plots[p.ID] = p
	return nil
}

func (r *InMemoryPlotRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.plots, id)
	return nil
}

type InMemoryOperationRepo struct {
	mu         sync.RWMutex
	operations map[string]*entity.Operation
}

func NewInMemoryOperationRepo() *InMemoryOperationRepo {
	return &InMemoryOperationRepo{operations: make(map[string]*entity.Operation)}
}

func (r *InMemoryOperationRepo) Create(op *entity.Operation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.operations[op.ID] = op
	return nil
}

func (r *InMemoryOperationRepo) GetByID(id string) (*entity.Operation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	op, ok := r.operations[id]
	if !ok {
		return nil, errors.New("operation not found")
	}
	return op, nil
}

func (r *InMemoryOperationRepo) ListByPlot(plotID string) ([]*entity.Operation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Operation
	for _, op := range r.operations {
		if op.PlotID == plotID {
			result = append(result, op)
		}
	}
	return result, nil
}

func (r *InMemoryOperationRepo) ListByTenant(tenantID string) ([]*entity.Operation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Operation
	for _, op := range r.operations {
		if op.TenantID == tenantID {
			result = append(result, op)
		}
	}
	return result, nil
}

func (r *InMemoryOperationRepo) ListByTenantAndPeriod(tenantID string, start, end string) ([]*entity.Operation, error) {
	return r.ListByTenant(tenantID)
}

func (r *InMemoryOperationRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.operations, id)
	return nil
}

type InMemoryHarvestRepo struct {
	mu       sync.RWMutex
	harvests map[string]*entity.Harvest
}

func NewInMemoryHarvestRepo() *InMemoryHarvestRepo {
	return &InMemoryHarvestRepo{harvests: make(map[string]*entity.Harvest)}
}

func (r *InMemoryHarvestRepo) Create(h *entity.Harvest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.harvests[h.ID] = h
	return nil
}

func (r *InMemoryHarvestRepo) GetByID(id string) (*entity.Harvest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.harvests[id]
	if !ok {
		return nil, errors.New("harvest not found")
	}
	return h, nil
}

func (r *InMemoryHarvestRepo) ListByTenant(tenantID string) ([]*entity.Harvest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Harvest
	for _, h := range r.harvests {
		if h.TenantID == tenantID {
			result = append(result, h)
		}
	}
	return result, nil
}

func (r *InMemoryHarvestRepo) Update(h *entity.Harvest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.harvests[h.ID] = h
	return nil
}

type InMemoryHarvestProductionRepo struct {
	mu          sync.RWMutex
	productions map[string]*entity.HarvestProduction
}

func NewInMemoryHarvestProductionRepo() *InMemoryHarvestProductionRepo {
	return &InMemoryHarvestProductionRepo{productions: make(map[string]*entity.HarvestProduction)}
}

func (r *InMemoryHarvestProductionRepo) Create(hp *entity.HarvestProduction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.productions[hp.ID] = hp
	return nil
}

func (r *InMemoryHarvestProductionRepo) GetByID(id string) (*entity.HarvestProduction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	hp, ok := r.productions[id]
	if !ok {
		return nil, errors.New("harvest production not found")
	}
	return hp, nil
}

func (r *InMemoryHarvestProductionRepo) ListByHarvest(harvestID string) ([]*entity.HarvestProduction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.HarvestProduction
	for _, hp := range r.productions {
		if hp.HarvestID == harvestID {
			result = append(result, hp)
		}
	}
	return result, nil
}

func (r *InMemoryHarvestProductionRepo) ListByPlot(plotID string) ([]*entity.HarvestProduction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.HarvestProduction
	for _, hp := range r.productions {
		if hp.PlotID == plotID {
			result = append(result, hp)
		}
	}
	return result, nil
}

func (r *InMemoryHarvestProductionRepo) Update(hp *entity.HarvestProduction) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.productions[hp.ID] = hp
	return nil
}

type InMemoryIndicatorRepo struct {
	mu         sync.RWMutex
	indicators map[string]*entity.Indicator
}

func NewInMemoryIndicatorRepo() *InMemoryIndicatorRepo {
	return &InMemoryIndicatorRepo{indicators: make(map[string]*entity.Indicator)}
}

func (r *InMemoryIndicatorRepo) Create(ind *entity.Indicator) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.indicators[ind.ID] = ind
	return nil
}

func (r *InMemoryIndicatorRepo) ListByHarvest(harvestID string) ([]*entity.Indicator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Indicator
	for _, ind := range r.indicators {
		if ind.HarvestID == harvestID {
			result = append(result, ind)
		}
	}
	return result, nil
}

func (r *InMemoryIndicatorRepo) ListByTenant(tenantID string) ([]*entity.Indicator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Indicator
	for _, ind := range r.indicators {
		if ind.TenantID == tenantID {
			result = append(result, ind)
		}
	}
	return result, nil
}

func (r *InMemoryIndicatorRepo) ListByTenantAndType(tenantID string, t entity.IndicatorType) ([]*entity.Indicator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Indicator
	for _, ind := range r.indicators {
		if ind.TenantID == tenantID && ind.Type == t {
			result = append(result, ind)
		}
	}
	return result, nil
}
