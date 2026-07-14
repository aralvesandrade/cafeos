package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type HarvestHandler struct {
	svc     *service.HarvestService
	plotSvc *service.PlotService
	farmSvc *service.FarmService
}

func NewHarvestHandler(svc *service.HarvestService, plotSvc *service.PlotService, farmSvc *service.FarmService) *HarvestHandler {
	return &HarvestHandler{svc: svc, plotSvc: plotSvc, farmSvc: farmSvc}
}

type createHarvestRequest struct {
	Year                int     `json:"year"`
	Description         string  `json:"description"`
	EstimatedProduction float64 `json:"estimated_production"`
}

// Create registra uma nova safra
// @Summary Criar colheita
// @Description Registra uma nova colheita (safra) para um ano
// @Tags harvests (Colheitas)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param harvest body createHarvestRequest true "Dados da colheita"
// @Success 201 {object} SwaggerHarvest
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests [post]
func (h *HarvestHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createHarvestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	harvest, err := h.svc.Create(organizationID, req.Year, req.Description, req.EstimatedProduction)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, harvest, http.StatusCreated)
}

// GetByID retorna uma colheita pelo seu ID
// @Summary Obter colheita por ID
// @Description Retorna uma única colheita (safra)
// @Tags harvests (Colheitas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Colheita"
// @Success 200 {object} SwaggerHarvest
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests/{id} [get]
func (h *HarvestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	harvest, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "harvest not found", http.StatusNotFound)
		return
	}
	writeJSON(w, harvest, http.StatusOK)
}

// List retorna todas as colheitas da organização autenticada
// @Summary Listar colheitas
// @Description Lista todas as colheitas (safras) da organização
// @Tags harvests (Colheitas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} SwaggerHarvest
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests [get]
func (h *HarvestHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	harvests, err := h.svc.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, harvests, http.StatusOK)
}

// Finalize marca uma colheita como concluída e calcula os indicadores
// @Summary Finalizar colheita
// @Description Finaliza uma colheita (safra), calculando todos os indicadores (sacas/ha, custo/saca, etc.)
// @Tags harvests (Colheitas)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Colheita"
// @Success 200 "OK"
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests/{id}/finalize [put]
func (h *HarvestHandler) Finalize(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Finalize(id); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type recordProductionRequest struct {
	PlotID   string  `json:"plot_id"`
	Quantity float64 `json:"quantity"`
	Notes    string  `json:"notes"`
}

// RecordProduction registra a produção de um talhão em uma colheita
// @Summary Registrar produção da colheita
// @Description Registra a quantidade produzida em um talhão específico dentro de uma colheita
// @Tags harvests (Colheitas)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Colheita"
// @Param production body recordProductionRequest true "Dados de produção"
// @Success 201 {object} SwaggerHarvestProduction
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests/{id}/production [post]
func (h *HarvestHandler) RecordProduction(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	harvestID := r.PathValue("id")

	var req recordProductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	prod, err := h.svc.RecordProduction(organizationID, harvestID, req.PlotID, req.Quantity, req.Notes)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, prod, http.StatusCreated)
}

// GetProduction retorna todos os registros de produção de uma colheita
// @Summary Obter produção da colheita
// @Description Obtém todos os registros de produção de uma colheita específica
// @Tags harvests (Colheitas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Colheita"
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Success 200 {array} SwaggerHarvestProduction
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests/{id}/production [get]
func (h *HarvestHandler) GetProduction(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	harvestID := r.PathValue("id")
	productions, err := h.svc.GetProductionByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userID, restricted := restrictedOwnerID(r); restricted {
		ownedFarms, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ownedPlots, err := h.plotSvc.PlotIDsForFarms(organizationID, ownedFarms)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filtered := productions[:0]
		for _, p := range productions {
			if ownedPlots[p.PlotID] {
				filtered = append(filtered, p)
			}
		}
		productions = filtered
	}

	if farmID := r.URL.Query().Get("farm_id"); farmID != "" {
		plotIDs, err := h.plotSvc.PlotIDsForFarm(organizationID, farmID)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filtered := productions[:0]
		for _, p := range productions {
			if plotIDs[p.PlotID] {
				filtered = append(filtered, p)
			}
		}
		productions = filtered
	}

	writeJSON(w, productions, http.StatusOK)
}
