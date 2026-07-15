package testing

import (
	"errors"
	"sync"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
)

type InMemoryFarmRepo struct {
	mu    sync.RWMutex
	farms map[string]*entity.Farm
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

func (r *InMemoryFarmRepo) ListByOrganization(organizationID string) ([]*entity.Farm, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Farm
	for _, f := range r.farms {
		if f.OrganizationID == organizationID {
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

type InMemoryProducerRepo struct {
	mu        sync.RWMutex
	producers map[string]*entity.Producer
}

func NewInMemoryProducerRepo() *InMemoryProducerRepo {
	return &InMemoryProducerRepo{producers: make(map[string]*entity.Producer)}
}

func (r *InMemoryProducerRepo) Create(p *entity.Producer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.producers[p.ID] = p
	return nil
}

func (r *InMemoryProducerRepo) GetByFarmID(farmID string) (*entity.Producer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.producers {
		if p.FarmID == farmID {
			return p, nil
		}
	}
	return nil, errors.New("producer not found")
}

func (r *InMemoryProducerRepo) Update(p *entity.Producer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.producers[p.ID] = p
	return nil
}

func (r *InMemoryProducerRepo) DeleteByFarmID(farmID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for id, p := range r.producers {
		if p.FarmID == farmID {
			delete(r.producers, id)
		}
	}
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

func (r *InMemoryPlotRepo) ListByOrganization(organizationID string) ([]*entity.Plot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Plot
	for _, p := range r.plots {
		if p.OrganizationID == organizationID {
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

func (r *InMemoryOperationRepo) Update(op *entity.Operation) error {
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

func (r *InMemoryOperationRepo) ListByOrganization(organizationID string) ([]*entity.Operation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Operation
	for _, op := range r.operations {
		if op.OrganizationID == organizationID {
			result = append(result, op)
		}
	}
	return result, nil
}

func (r *InMemoryOperationRepo) ListByOrganizationAndPeriod(organizationID string, start, end string) ([]*entity.Operation, error) {
	return r.ListByOrganization(organizationID)
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

func (r *InMemoryHarvestRepo) ListByOrganization(organizationID string) ([]*entity.Harvest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Harvest
	for _, h := range r.harvests {
		if h.OrganizationID == organizationID {
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

func (r *InMemoryIndicatorRepo) ListByOrganization(organizationID string) ([]*entity.Indicator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Indicator
	for _, ind := range r.indicators {
		if ind.OrganizationID == organizationID {
			result = append(result, ind)
		}
	}
	return result, nil
}

func (r *InMemoryIndicatorRepo) ListByOrganizationAndType(organizationID string, t entity.IndicatorType) ([]*entity.Indicator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Indicator
	for _, ind := range r.indicators {
		if ind.OrganizationID == organizationID && ind.Type == t {
			result = append(result, ind)
		}
	}
	return result, nil
}

type InMemoryPermissionRepo struct {
	mu          sync.RWMutex
	permissions map[string]*entity.RolePermission
}

func NewInMemoryPermissionRepo() *InMemoryPermissionRepo {
	return &InMemoryPermissionRepo{permissions: make(map[string]*entity.RolePermission)}
}

func (r *InMemoryPermissionRepo) key(organizationID, roleID string, module entity.ModuleKey) string {
	return organizationID + "|" + roleID + "|" + string(module)
}

func (r *InMemoryPermissionRepo) ListByOrganization(organizationID string) ([]*entity.RolePermission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.RolePermission
	for _, p := range r.permissions {
		if p.OrganizationID == organizationID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *InMemoryPermissionRepo) Upsert(p *entity.RolePermission) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.permissions[r.key(p.OrganizationID, p.RoleID, p.Module)] = p
	return nil
}

func (r *InMemoryPermissionRepo) CountByOrganization(organizationID string) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var count int64
	for _, p := range r.permissions {
		if p.OrganizationID == organizationID {
			count++
		}
	}
	return count, nil
}

type InMemoryRoleRepo struct {
	mu    sync.RWMutex
	roles map[string]*entity.Role
}

func NewInMemoryRoleRepo() *InMemoryRoleRepo {
	return &InMemoryRoleRepo{roles: make(map[string]*entity.Role)}
}

func (r *InMemoryRoleRepo) List() ([]*entity.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*entity.Role
	for _, role := range r.roles {
		result = append(result, role)
	}
	return result, nil
}

func (r *InMemoryRoleRepo) GetByID(id string) (*entity.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	role, ok := r.roles[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return role, nil
}

func (r *InMemoryRoleRepo) FindByKey(key string) (*entity.Role, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, role := range r.roles {
		if role.Key == key {
			return role, nil
		}
	}
	return nil, errors.New("role not found")
}

func (r *InMemoryRoleRepo) Create(role *entity.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.roles[role.ID] = role
	return nil
}

func (r *InMemoryRoleRepo) Update(role *entity.Role) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.roles[role.ID] = role
	return nil
}

func (r *InMemoryRoleRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.roles, id)
	return nil
}

func (r *InMemoryRoleRepo) Count() (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return int64(len(r.roles)), nil
}
