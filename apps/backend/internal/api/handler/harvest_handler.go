package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type HarvestHandler struct {
	svc *service.HarvestService
}

func NewHarvestHandler(svc *service.HarvestService) *HarvestHandler {
	return &HarvestHandler{svc: svc}
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
// @Param tenant_id path string true "ID do Tenant"
// @Param harvest body createHarvestRequest true "Dados da colheita"
// @Success 201 {object} SwaggerHarvest
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests [post]
func (h *HarvestHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createHarvestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	harvest, err := h.svc.Create(tenantID, req.Year, req.Description, req.EstimatedProduction)
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
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Colheita"
// @Success 200 {object} SwaggerHarvest
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id} [get]
func (h *HarvestHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	harvest, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "harvest not found", http.StatusNotFound)
		return
	}
	writeJSON(w, harvest, http.StatusOK)
}

// List retorna todas as colheitas do tenant autenticado
// @Summary Listar colheitas
// @Description Lista todas as colheitas (safras) do tenant
// @Tags harvests (Colheitas)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} SwaggerHarvest
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests [get]
func (h *HarvestHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	harvests, err := h.svc.ListByTenant(tenantID)
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
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Colheita"
// @Success 200 "OK"
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id}/finalize [put]
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
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Colheita"
// @Param production body recordProductionRequest true "Dados de produção"
// @Success 201 {object} SwaggerHarvestProduction
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id}/production [post]
func (h *HarvestHandler) RecordProduction(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	harvestID := r.PathValue("id")

	var req recordProductionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	prod, err := h.svc.RecordProduction(tenantID, harvestID, req.PlotID, req.Quantity, req.Notes)
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
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Colheita"
// @Success 200 {array} SwaggerHarvestProduction
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/harvests/{id}/production [get]
func (h *HarvestHandler) GetProduction(w http.ResponseWriter, r *http.Request) {
	harvestID := r.PathValue("id")
	productions, err := h.svc.GetProductionByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, productions, http.StatusOK)
}
