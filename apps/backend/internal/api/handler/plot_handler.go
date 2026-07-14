package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type PlotHandler struct {
	svc     *service.PlotService
	farmSvc *service.FarmService
}

func NewPlotHandler(svc *service.PlotService, farmSvc *service.FarmService) *PlotHandler {
	return &PlotHandler{svc: svc, farmSvc: farmSvc}
}

// canAccessPlot mirrors FarmHandler.canAccessFarm: proprietario is restricted
// to plots on farms it owns, every other role sees all plots in its organization.
func (h *PlotHandler) canAccessPlot(r *http.Request, organizationID string, plot *entity.Plot) bool {
	userID, restricted := restrictedOwnerID(r)
	if !restricted {
		return true
	}
	owned, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
	if err != nil {
		return false
	}
	return owned[plot.FarmID]
}

func parseDateField(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type createPlotRequest struct {
	FarmID       string  `json:"farm_id"`
	Name         string  `json:"name"`
	AreaHA       float64 `json:"area_ha"`
	Cultivar     string  `json:"cultivar"`
	PlantingYear int     `json:"planting_year"`
	Altitude     int     `json:"altitude"`
	SoilType     string  `json:"soil_type"`

	Leased     bool   `json:"leased"`
	Stage      string `json:"stage"`
	Irrigation string `json:"irrigation"`

	ActivationDate   *string `json:"activation_date"`
	PlantingDate     *string `json:"planting_date"`
	DeactivationDate *string `json:"deactivation_date"`

	Intercropped  bool   `json:"intercropped"`
	SecondaryCrop string `json:"secondary_crop"`
	Notes         string `json:"notes"`

	CropType string `json:"crop_type"`

	FormationCostPerHA float64 `json:"formation_cost_per_ha"`
	UsefulLifeYears    int     `json:"useful_life_years"`
	RowSpacingM        float64 `json:"row_spacing_m"`
	PlantSpacingM      float64 `json:"plant_spacing_m"`
	PlantCount         int     `json:"plant_count"`

	DamAreaHA          float64 `json:"dam_area_ha"`
	ImprovementsAreaHA float64 `json:"improvements_area_ha"`
	RoadsAreaHA        float64 `json:"roads_area_ha"`
	APPAreaHA          float64 `json:"app_area_ha"`
	LegalReserveAreaHA float64 `json:"legal_reserve_area_ha"`
}

func (req *createPlotRequest) toEntity(organizationID string) (*entity.Plot, error) {
	activation, err := parseDateField(req.ActivationDate)
	if err != nil {
		return nil, err
	}
	planting, err := parseDateField(req.PlantingDate)
	if err != nil {
		return nil, err
	}
	deactivation, err := parseDateField(req.DeactivationDate)
	if err != nil {
		return nil, err
	}

	return &entity.Plot{
		OrganizationID:     organizationID,
		FarmID:             req.FarmID,
		Name:               req.Name,
		AreaHA:             req.AreaHA,
		Cultivar:           req.Cultivar,
		PlantingYear:       req.PlantingYear,
		Altitude:           req.Altitude,
		SoilType:           req.SoilType,
		Leased:             req.Leased,
		Stage:              entity.PlotStage(req.Stage),
		Irrigation:         req.Irrigation,
		ActivationDate:     activation,
		PlantingDate:       planting,
		DeactivationDate:   deactivation,
		Intercropped:       req.Intercropped,
		SecondaryCrop:      req.SecondaryCrop,
		Notes:              req.Notes,
		CropType:           req.CropType,
		FormationCostPerHA: req.FormationCostPerHA,
		UsefulLifeYears:    req.UsefulLifeYears,
		RowSpacingM:        req.RowSpacingM,
		PlantSpacingM:      req.PlantSpacingM,
		PlantCount:         req.PlantCount,
		DamAreaHA:          req.DamAreaHA,
		ImprovementsAreaHA: req.ImprovementsAreaHA,
		RoadsAreaHA:        req.RoadsAreaHA,
		APPAreaHA:          req.APPAreaHA,
		LegalReserveAreaHA: req.LegalReserveAreaHA,
	}, nil
}

// Create registra um novo talhão para a organização autenticada
// @Summary Criar talhão
// @Description Registra um novo talhão em uma fazenda
// @Tags plots (Talhões)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param plot body createPlotRequest true "Dados do talhão"
// @Success 201 {object} SwaggerPlot
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/plots [post]
func (h *PlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createPlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	plotEntity, err := req.toEntity(organizationID)
	if err != nil {
		writeError(w, "invalid date, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	plot, err := h.svc.Create(plotEntity)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, plot, http.StatusCreated)
}

// GetByID retorna um talhão pelo seu ID
// @Summary Obter talhão por ID
// @Description Retorna um único talhão
// @Tags plots (Talhões)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Talhão"
// @Success 200 {object} SwaggerPlot
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/plots/{id} [get]
func (h *PlotHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	plot, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}
	if !h.canAccessPlot(r, organizationID, plot) {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}
	writeJSON(w, plot, http.StatusOK)
}

// ListByFarm retorna todos os talhões de uma fazenda
// @Summary Listar talhões por fazenda
// @Description Lista todos os talhões pertencentes a uma fazenda
// @Tags plots (Talhões)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param farm_id path string true "ID da Fazenda"
// @Success 200 {array} SwaggerPlot
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/farms/{farm_id}/plots [get]
func (h *PlotHandler) ListByFarm(w http.ResponseWriter, r *http.Request) {
	farmID := r.PathValue("farm_id")
	plots, err := h.svc.ListByFarm(farmID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, plots, http.StatusOK)
}

// List retorna todos os talhões da organização autenticada
// @Summary Listar todos os talhões
// @Description Lista todos os talhões de todas as fazendas da organização
// @Tags plots (Talhões)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} SwaggerPlot
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Router /api/v1/{organization_id}/plots [get]
func (h *PlotHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	plots, err := h.svc.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userID, restricted := restrictedOwnerID(r); restricted {
		owned, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filtered := plots[:0]
		for _, p := range plots {
			if owned[p.FarmID] {
				filtered = append(filtered, p)
			}
		}
		plots = filtered
	}

	if farmID := r.URL.Query().Get("farm_id"); farmID != "" {
		filtered := plots[:0]
		for _, p := range plots {
			if p.FarmID == farmID {
				filtered = append(filtered, p)
			}
		}
		plots = filtered
	}

	writeJSON(w, plots, http.StatusOK)
}

// Update atualiza um talhão existente
// @Summary Atualizar talhão
// @Description Atualiza dados do talhão por ID (atualização parcial - somente os campos informados são alterados)
// @Tags plots (Talhões)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Talhão"
// @Param plot body SwaggerPlot true "Dados atualizados do talhão"
// @Success 200 {object} SwaggerPlot
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/plots/{id} [put]
func (h *PlotHandler) Update(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")

	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}
	if !h.canAccessPlot(r, organizationID, existing) {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name         *string  `json:"name"`
		FarmID       *string  `json:"farm_id"`
		AreaHA       *float64 `json:"area_ha"`
		Cultivar     *string  `json:"cultivar"`
		SoilType     *string  `json:"soil_type"`
		Altitude     *int     `json:"altitude"`
		PlantingYear *int     `json:"planting_year"`

		Leased     *bool   `json:"leased"`
		Stage      *string `json:"stage"`
		Irrigation *string `json:"irrigation"`

		ActivationDate   *string `json:"activation_date"`
		PlantingDate     *string `json:"planting_date"`
		DeactivationDate *string `json:"deactivation_date"`

		Intercropped  *bool   `json:"intercropped"`
		SecondaryCrop *string `json:"secondary_crop"`
		Notes         *string `json:"notes"`

		CropType *string `json:"crop_type"`

		FormationCostPerHA *float64 `json:"formation_cost_per_ha"`
		UsefulLifeYears    *int     `json:"useful_life_years"`
		RowSpacingM        *float64 `json:"row_spacing_m"`
		PlantSpacingM      *float64 `json:"plant_spacing_m"`
		PlantCount         *int     `json:"plant_count"`

		DamAreaHA          *float64 `json:"dam_area_ha"`
		ImprovementsAreaHA *float64 `json:"improvements_area_ha"`
		RoadsAreaHA        *float64 `json:"roads_area_ha"`
		APPAreaHA          *float64 `json:"app_area_ha"`
		LegalReserveAreaHA *float64 `json:"legal_reserve_area_ha"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.FarmID != nil {
		existing.FarmID = *input.FarmID
	}
	if input.AreaHA != nil {
		existing.AreaHA = *input.AreaHA
	}
	if input.Cultivar != nil {
		existing.Cultivar = *input.Cultivar
	}
	if input.SoilType != nil {
		existing.SoilType = *input.SoilType
	}
	if input.Altitude != nil {
		existing.Altitude = *input.Altitude
	}
	if input.PlantingYear != nil {
		existing.PlantingYear = *input.PlantingYear
	}
	if input.Leased != nil {
		existing.Leased = *input.Leased
	}
	if input.Stage != nil {
		existing.Stage = entity.PlotStage(*input.Stage)
	}
	if input.Irrigation != nil {
		existing.Irrigation = *input.Irrigation
	}
	if input.ActivationDate != nil {
		t, err := parseDateField(input.ActivationDate)
		if err != nil {
			writeError(w, "invalid activation_date, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		existing.ActivationDate = t
	}
	if input.PlantingDate != nil {
		t, err := parseDateField(input.PlantingDate)
		if err != nil {
			writeError(w, "invalid planting_date, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		existing.PlantingDate = t
	}
	if input.DeactivationDate != nil {
		t, err := parseDateField(input.DeactivationDate)
		if err != nil {
			writeError(w, "invalid deactivation_date, expected YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		existing.DeactivationDate = t
	}
	if input.Intercropped != nil {
		existing.Intercropped = *input.Intercropped
	}
	if input.SecondaryCrop != nil {
		existing.SecondaryCrop = *input.SecondaryCrop
	}
	if input.Notes != nil {
		existing.Notes = *input.Notes
	}
	if input.CropType != nil {
		existing.CropType = *input.CropType
	}
	if input.FormationCostPerHA != nil {
		existing.FormationCostPerHA = *input.FormationCostPerHA
	}
	if input.UsefulLifeYears != nil {
		existing.UsefulLifeYears = *input.UsefulLifeYears
	}
	if input.RowSpacingM != nil {
		existing.RowSpacingM = *input.RowSpacingM
	}
	if input.PlantSpacingM != nil {
		existing.PlantSpacingM = *input.PlantSpacingM
	}
	if input.PlantCount != nil {
		existing.PlantCount = *input.PlantCount
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

	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// Delete remove um talhão
// @Summary Excluir talhão
// @Description Exclui um talhão por ID
// @Tags plots (Talhões)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Talhão"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/plots/{id} [delete]
func (h *PlotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}
	if !h.canAccessPlot(r, organizationID, existing) {
		writeError(w, "plot not found", http.StatusNotFound)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
