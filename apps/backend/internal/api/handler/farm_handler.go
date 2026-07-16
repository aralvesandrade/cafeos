package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

// Types aliases for swagger documentation
type (
	SwaggerFarm              = entity.Farm
	SwaggerProducer          = entity.Producer
	SwaggerPlot              = entity.Plot
	SwaggerOperation         = entity.Operation
	SwaggerHarvest           = entity.Harvest
	SwaggerHarvestProduction = entity.HarvestProduction
	SwaggerIndicator         = entity.Indicator
)

type FarmHandler struct {
	svc     *service.FarmService
	userSvc *service.UserService
	roleSvc *service.RoleService
}

func NewFarmHandler(svc *service.FarmService, userSvc *service.UserService, roleSvc *service.RoleService) *FarmHandler {
	return &FarmHandler{svc: svc, userSvc: userSvc, roleSvc: roleSvc}
}

type producerRequest struct {
	UserID        string  `json:"user_id"`
	RoleID        string  `json:"role_id"`
	CPF           string  `json:"cpf"`
	Name          string  `json:"name"`
	RG            string  `json:"rg"`
	IssuingBody   string  `json:"issuing_body"`
	Gender        string  `json:"gender"`
	BirthDate     *string `json:"birth_date"`
	MaritalStatus string  `json:"marital_status"`
	Phone         string  `json:"phone"`
	Email         string  `json:"email"`
	Education     string  `json:"education"`
}

func (p *producerRequest) toEntity() (*entity.Producer, error) {
	producer := &entity.Producer{
		UserID:        p.UserID,
		RoleID:        p.RoleID,
		CPF:           p.CPF,
		Name:          p.Name,
		RG:            p.RG,
		IssuingBody:   p.IssuingBody,
		Gender:        p.Gender,
		MaritalStatus: p.MaritalStatus,
		Phone:         p.Phone,
		Email:         p.Email,
		Education:     p.Education,
	}
	if p.BirthDate != nil && *p.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *p.BirthDate)
		if err != nil {
			return nil, err
		}
		producer.BirthDate = &t
	}
	return producer, nil
}

func toProducerEntities(reqs []producerRequest) ([]*entity.Producer, error) {
	producers := make([]*entity.Producer, 0, len(reqs))
	for i := range reqs {
		p, err := reqs[i].toEntity()
		if err != nil {
			return nil, err
		}
		producers = append(producers, p)
	}
	return producers, nil
}

type createFarmRequest struct {
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Location string `json:"location"`

	Phone                    string `json:"phone"`
	Activities               string `json:"activities"`
	MainCrop                 string `json:"main_crop"`
	SecondaryCrop            string `json:"secondary_crop"`
	State                    string `json:"state"`
	City                     string `json:"city"`
	Address                  string `json:"address"`
	ProductionSystem         string `json:"production_system"`
	CommercializationProduct string `json:"commercialization_product"`

	HasNoCNPJ              bool   `json:"has_no_cnpj"`
	CNPJ                   string `json:"cnpj"`
	HasNoNIRF              bool   `json:"has_no_nirf"`
	NIRF                   string `json:"nirf"`
	HasNoINCRA             bool   `json:"has_no_incra"`
	INCRARegistration      string `json:"incra_registration"`
	HasNoStateRegistration bool   `json:"has_no_state_registration"`
	StateRegistration      string `json:"state_registration"`
	HasNoDAP               bool   `json:"has_no_dap"`
	DAP                    string `json:"dap"`
	HasNoCAR               bool   `json:"has_no_car"`
	CAR                    string `json:"car"`

	FullyLeased    bool    `json:"fully_leased"`
	LandValuePerHA float64 `json:"land_value_per_ha"`

	TotalAreaHA   float64 `json:"total_area_ha"`
	PlantedAreaHA float64 `json:"planted_area_ha"`

	DamAreaHA                   float64 `json:"dam_area_ha"`
	ImprovementsAreaHA          float64 `json:"improvements_area_ha"`
	RoadsAreaHA                 float64 `json:"roads_area_ha"`
	APPAreaHA                   float64 `json:"app_area_ha"`
	LegalReserveAreaHA          float64 `json:"legal_reserve_area_ha"`
	NativeVegetationAreaHA      float64 `json:"native_vegetation_area_ha"`
	LivestockAreaNotCoveredHA   float64 `json:"livestock_area_not_covered_ha"`
	AgricultureAreaNotCoveredHA float64 `json:"agriculture_area_not_covered_ha"`
	NonAgriculturalAreaHA       float64 `json:"non_agricultural_area_ha"`

	Producers []producerRequest `json:"producers"`
}

func (req *createFarmRequest) toEntity(organizationID string) *entity.Farm {
	return &entity.Farm{
		OrganizationID:              organizationID,
		Name:                        req.Name,
		Owner:                       req.Owner,
		Location:                    req.Location,
		Phone:                       req.Phone,
		Activities:                  req.Activities,
		MainCrop:                    req.MainCrop,
		SecondaryCrop:               req.SecondaryCrop,
		State:                       req.State,
		City:                        req.City,
		Address:                     req.Address,
		ProductionSystem:            req.ProductionSystem,
		CommercializationProduct:    req.CommercializationProduct,
		HasNoCNPJ:                   req.HasNoCNPJ,
		CNPJ:                        req.CNPJ,
		HasNoNIRF:                   req.HasNoNIRF,
		NIRF:                        req.NIRF,
		HasNoINCRA:                  req.HasNoINCRA,
		INCRARegistration:           req.INCRARegistration,
		HasNoStateRegistration:      req.HasNoStateRegistration,
		StateRegistration:           req.StateRegistration,
		HasNoDAP:                    req.HasNoDAP,
		DAP:                         req.DAP,
		HasNoCAR:                    req.HasNoCAR,
		CAR:                         req.CAR,
		FullyLeased:                 req.FullyLeased,
		LandValuePerHA:              req.LandValuePerHA,
		TotalAreaHA:                 req.TotalAreaHA,
		PlantedAreaHA:               req.PlantedAreaHA,
		DamAreaHA:                   req.DamAreaHA,
		ImprovementsAreaHA:          req.ImprovementsAreaHA,
		RoadsAreaHA:                 req.RoadsAreaHA,
		APPAreaHA:                   req.APPAreaHA,
		LegalReserveAreaHA:          req.LegalReserveAreaHA,
		NativeVegetationAreaHA:      req.NativeVegetationAreaHA,
		LivestockAreaNotCoveredHA:   req.LivestockAreaNotCoveredHA,
		AgricultureAreaNotCoveredHA: req.AgricultureAreaNotCoveredHA,
		NonAgriculturalAreaHA:       req.NonAgriculturalAreaHA,
	}
}

// Create registra uma nova fazenda para a organização autenticada
// @Summary Criar fazenda
// @Description Registra uma nova fazenda na organização
// @Tags farms (Fazendas)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param farm body createFarmRequest true "Dados da fazenda"
// @Success 201 {object} SwaggerFarm
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/farms [post]
func (h *FarmHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	callerUserID, _ := r.Context().Value(middleware.UserIDKey).(string)
	callerRole, _ := r.Context().Value(middleware.RoleKey).(string)

	if callerRole != entity.SystemRolePlatformOwner {
		isPrincipal, err := h.userSvc.IsPrincipal(callerUserID)
		if err != nil {
			writeError(w, "invalid caller", http.StatusBadRequest)
			return
		}
		if !isPrincipal {
			writeError(w, service.ErrNotPrincipal.Error(), http.StatusForbidden)
			return
		}
	}

	var req createFarmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	producers, err := toProducerEntities(req.Producers)
	if err != nil {
		writeError(w, "invalid producer birth_date, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Cada user_id vinculado no payload precisa ser o próprio principal ou
	// um sub-usuário dele — nunca um usuário de fora da sua hierarquia.
	if callerRole != entity.SystemRolePlatformOwner {
		hasCallerAsProducer := false
		for _, p := range producers {
			if p.UserID == callerUserID {
				hasCallerAsProducer = true
				continue
			}
			if err := h.userSvc.EnsureManages(callerUserID, p.UserID); err != nil {
				writeError(w, "user_id must be the principal themselves or one of their sub-users", http.StatusForbidden)
				return
			}
		}
		if !hasCallerAsProducer {
			role, err := h.roleSvc.FindByKey(entity.RoleKeyProprietario)
			if err != nil {
				writeError(w, "failed to resolve owner role", http.StatusInternalServerError)
				return
			}
			caller, err := h.userSvc.Get(callerUserID)
			if err != nil {
				writeError(w, "invalid caller", http.StatusBadRequest)
				return
			}
			producers = append(producers, &entity.Producer{UserID: callerUserID, RoleID: role.ID, Name: caller.Name, Email: caller.Email})
		}
	}

	farm, err := h.svc.Create(req.toEntity(organizationID), producers)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, farm, http.StatusCreated)
}

// GetByID retorna uma fazenda pelo seu ID
// @Summary Obter fazenda por ID
// @Description Retorna uma única fazenda
// @Tags farms (Fazendas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Fazenda"
// @Success 200 {object} SwaggerFarm
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/farms/{id} [get]
func (h *FarmHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	farm, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "farm not found", http.StatusNotFound)
		return
	}
	if !h.canAccessFarm(r, id) {
		writeError(w, "farm not found", http.StatusNotFound)
		return
	}
	writeJSON(w, farm, http.StatusOK)
}

// canAccessFarm reports whether the authenticated request may access the given
// farm: proprietario is restricted to farms it owns (via the linked producer),
// every other role sees all farms in its organization.
func (h *FarmHandler) canAccessFarm(r *http.Request, farmID string) bool {
	role, _ := r.Context().Value(middleware.RoleKey).(string)
	if role != entity.RoleKeyProprietario {
		return true
	}
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	return h.svc.IsOwner(farmID, userID)
}

// List retorna todas as fazendas da organização autenticada
// @Summary Listar fazendas
// @Description Lista todas as fazendas pertencentes à organização
// @Tags farms (Fazendas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} SwaggerFarm
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/farms [get]
func (h *FarmHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	var farms []*entity.Farm
	var err error
	if role == entity.RoleKeyProprietario {
		userID, _ := r.Context().Value(middleware.UserIDKey).(string)
		farms, err = h.svc.ListByOwner(organizationID, userID)
	} else {
		farms, err = h.svc.ListByOrganization(organizationID)
	}
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, farms, http.StatusOK)
}

// Update atualiza uma fazenda existente
// @Summary Atualizar fazenda
// @Description Atualiza dados da fazenda por ID (atualização parcial - somente os campos informados são alterados)
// @Tags farms (Fazendas)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Fazenda"
// @Param farm body SwaggerFarm true "Dados atualizados da fazenda"
// @Success 200 {object} SwaggerFarm
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/farms/{id} [put]
func (h *FarmHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "farm not found", http.StatusNotFound)
		return
	}
	if !h.canAccessFarm(r, id) {
		writeError(w, "farm not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name     *string `json:"name"`
		Owner    *string `json:"owner"`
		Location *string `json:"location"`

		Phone                    *string `json:"phone"`
		Activities               *string `json:"activities"`
		MainCrop                 *string `json:"main_crop"`
		SecondaryCrop            *string `json:"secondary_crop"`
		State                    *string `json:"state"`
		City                     *string `json:"city"`
		Address                  *string `json:"address"`
		ProductionSystem         *string `json:"production_system"`
		CommercializationProduct *string `json:"commercialization_product"`

		HasNoCNPJ              *bool   `json:"has_no_cnpj"`
		CNPJ                   *string `json:"cnpj"`
		HasNoNIRF              *bool   `json:"has_no_nirf"`
		NIRF                   *string `json:"nirf"`
		HasNoINCRA             *bool   `json:"has_no_incra"`
		INCRARegistration      *string `json:"incra_registration"`
		HasNoStateRegistration *bool   `json:"has_no_state_registration"`
		StateRegistration      *string `json:"state_registration"`
		HasNoDAP               *bool   `json:"has_no_dap"`
		DAP                    *string `json:"dap"`
		HasNoCAR               *bool   `json:"has_no_car"`
		CAR                    *string `json:"car"`

		FullyLeased    *bool    `json:"fully_leased"`
		LandValuePerHA *float64 `json:"land_value_per_ha"`

		TotalAreaHA   *float64 `json:"total_area_ha"`
		PlantedAreaHA *float64 `json:"planted_area_ha"`

		DamAreaHA                   *float64 `json:"dam_area_ha"`
		ImprovementsAreaHA          *float64 `json:"improvements_area_ha"`
		RoadsAreaHA                 *float64 `json:"roads_area_ha"`
		APPAreaHA                   *float64 `json:"app_area_ha"`
		LegalReserveAreaHA          *float64 `json:"legal_reserve_area_ha"`
		NativeVegetationAreaHA      *float64 `json:"native_vegetation_area_ha"`
		LivestockAreaNotCoveredHA   *float64 `json:"livestock_area_not_covered_ha"`
		AgricultureAreaNotCoveredHA *float64 `json:"agriculture_area_not_covered_ha"`
		NonAgriculturalAreaHA       *float64 `json:"non_agricultural_area_ha"`

		Producers *[]producerRequest `json:"producers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Owner != nil {
		existing.Owner = *input.Owner
	}
	if input.Location != nil {
		existing.Location = *input.Location
	}
	if input.Phone != nil {
		existing.Phone = *input.Phone
	}
	if input.Activities != nil {
		existing.Activities = *input.Activities
	}
	if input.MainCrop != nil {
		existing.MainCrop = *input.MainCrop
	}
	if input.SecondaryCrop != nil {
		existing.SecondaryCrop = *input.SecondaryCrop
	}
	if input.State != nil {
		existing.State = *input.State
	}
	if input.City != nil {
		existing.City = *input.City
	}
	if input.Address != nil {
		existing.Address = *input.Address
	}
	if input.ProductionSystem != nil {
		existing.ProductionSystem = *input.ProductionSystem
	}
	if input.CommercializationProduct != nil {
		existing.CommercializationProduct = *input.CommercializationProduct
	}
	if input.HasNoCNPJ != nil {
		existing.HasNoCNPJ = *input.HasNoCNPJ
	}
	if input.CNPJ != nil {
		existing.CNPJ = *input.CNPJ
	}
	if input.HasNoNIRF != nil {
		existing.HasNoNIRF = *input.HasNoNIRF
	}
	if input.NIRF != nil {
		existing.NIRF = *input.NIRF
	}
	if input.HasNoINCRA != nil {
		existing.HasNoINCRA = *input.HasNoINCRA
	}
	if input.INCRARegistration != nil {
		existing.INCRARegistration = *input.INCRARegistration
	}
	if input.HasNoStateRegistration != nil {
		existing.HasNoStateRegistration = *input.HasNoStateRegistration
	}
	if input.StateRegistration != nil {
		existing.StateRegistration = *input.StateRegistration
	}
	if input.HasNoDAP != nil {
		existing.HasNoDAP = *input.HasNoDAP
	}
	if input.DAP != nil {
		existing.DAP = *input.DAP
	}
	if input.HasNoCAR != nil {
		existing.HasNoCAR = *input.HasNoCAR
	}
	if input.CAR != nil {
		existing.CAR = *input.CAR
	}
	if input.FullyLeased != nil {
		existing.FullyLeased = *input.FullyLeased
	}
	if input.LandValuePerHA != nil {
		existing.LandValuePerHA = *input.LandValuePerHA
	}
	if input.TotalAreaHA != nil {
		existing.TotalAreaHA = *input.TotalAreaHA
	}
	if input.PlantedAreaHA != nil {
		existing.PlantedAreaHA = *input.PlantedAreaHA
	}
	if input.DamAreaHA != nil {
		existing.DamAreaHA = *input.DamAreaHA
	}
	if input.ImprovementsAreaHA != nil {
		existing.ImprovementsAreaHA = *input.ImprovementsAreaHA
	}
	if input.RoadsAreaHA != nil {
		existing.RoadsAreaHA = *input.RoadsAreaHA
	}
	if input.APPAreaHA != nil {
		existing.APPAreaHA = *input.APPAreaHA
	}
	if input.LegalReserveAreaHA != nil {
		existing.LegalReserveAreaHA = *input.LegalReserveAreaHA
	}
	if input.NativeVegetationAreaHA != nil {
		existing.NativeVegetationAreaHA = *input.NativeVegetationAreaHA
	}
	if input.LivestockAreaNotCoveredHA != nil {
		existing.LivestockAreaNotCoveredHA = *input.LivestockAreaNotCoveredHA
	}
	if input.AgricultureAreaNotCoveredHA != nil {
		existing.AgricultureAreaNotCoveredHA = *input.AgricultureAreaNotCoveredHA
	}
	if input.NonAgriculturalAreaHA != nil {
		existing.NonAgriculturalAreaHA = *input.NonAgriculturalAreaHA
	}

	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Producers != nil {
		callerUserID, _ := r.Context().Value(middleware.UserIDKey).(string)
		callerRole, _ := r.Context().Value(middleware.RoleKey).(string)
		if callerRole != entity.SystemRolePlatformOwner {
			isPrincipal, err := h.userSvc.IsPrincipal(callerUserID)
			if err != nil {
				writeError(w, "invalid caller", http.StatusBadRequest)
				return
			}
			if !isPrincipal || !h.svc.IsOwner(existing.ID, callerUserID) {
				writeError(w, "only a principal that owns this farm can link users to it", http.StatusForbidden)
				return
			}
		}

		producers, err := toProducerEntities(*input.Producers)
		if err != nil {
			writeError(w, "invalid producer birth_date, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		if callerRole != entity.SystemRolePlatformOwner {
			for _, p := range producers {
				if p.UserID == callerUserID {
					continue
				}
				if err := h.userSvc.EnsureManages(callerUserID, p.UserID); err != nil {
					writeError(w, "user_id must be the principal themselves or one of their sub-users", http.StatusForbidden)
					return
				}
			}
		}

		saved, err := h.svc.SetProducers(existing.ID, producers)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
		existing.Producers = make([]entity.Producer, len(saved))
		for i, p := range saved {
			existing.Producers[i] = *p
		}
	}

	writeJSON(w, existing, http.StatusOK)
}

// Delete remove uma fazenda
// @Summary Excluir fazenda
// @Description Exclui uma fazenda por ID
// @Tags farms (Fazendas)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Fazenda"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/farms/{id} [delete]
func (h *FarmHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !h.canAccessFarm(r, id) {
		writeError(w, "farm not found", http.StatusNotFound)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
