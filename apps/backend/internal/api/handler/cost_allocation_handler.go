package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type CostAllocationHandler struct {
	svc *service.CostAllocationService
}

func NewCostAllocationHandler(svc *service.CostAllocationService) *CostAllocationHandler {
	return &CostAllocationHandler{svc: svc}
}

type createAllocationRequest struct {
	HarvestID    string             `json:"harvest_id"`
	CostCenterID string             `json:"cost_center_id"`
	Description  string             `json:"description"`
	TotalAmount  float64            `json:"total_amount"`
	Method       string             `json:"method"`
	Date         string             `json:"date"`
	Percentages  map[string]float64 `json:"percentages"`
}

type allocationItemResponse struct {
	PlotID     string  `json:"plot_id"`
	PlotName   string  `json:"plot_name"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// Create registra um novo rateio de custo para um centro de custo da colheita
// @Summary Criar rateio de custo
// @Description Registra um novo rateio de custo, distribuindo um valor entre talhões
// @Tags cost-allocations (Rateios de Custo)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param allocation body createAllocationRequest true "Dados do rateio de custo"
// @Success 201 {object} entity.CostAllocation
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/cost-allocations [post]
func (h *CostAllocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createAllocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	input := service.CreateAllocationInput{
		OrganizationID: organizationID,
		HarvestID:      req.HarvestID,
		CostCenterID:   req.CostCenterID,
		Description:    req.Description,
		TotalAmount:    req.TotalAmount,
		Method:         entity.AllocationMethod(req.Method),
		Date:           date,
		Percentages:    req.Percentages,
	}

	allocation, err := h.svc.Create(input)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, allocation, http.StatusCreated)
}

// ListByHarvest retorna todos os rateios de custo de uma colheita
// @Summary Listar rateios de custo por colheita
// @Description Lista todos os rateios de custo pertencentes a uma colheita
// @Tags cost-allocations (Rateios de Custo)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param harvest_id path string true "ID da Colheita"
// @Success 200 {array} entity.CostAllocation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/harvests/{harvest_id}/cost-allocations [get]
func (h *CostAllocationHandler) ListByHarvest(w http.ResponseWriter, r *http.Request) {
	harvestID := r.PathValue("harvest_id")
	allocs, err := h.svc.ListByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, allocs, http.StatusOK)
}

// GetByID retorna um rateio de custo pelo seu ID
// @Summary Obter rateio de custo por ID
// @Description Retorna um único rateio de custo
// @Tags cost-allocations (Rateios de Custo)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Rateio de Custo"
// @Success 200 {object} entity.CostAllocation
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/cost-allocations/{id} [get]
func (h *CostAllocationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "cost allocation not found", http.StatusNotFound)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

// Delete remove um rateio de custo pelo seu ID
// @Summary Excluir rateio de custo
// @Description Exclui um rateio de custo por ID
// @Tags cost-allocations (Rateios de Custo)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Rateio de Custo"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/cost-allocations/{id} [delete]
func (h *CostAllocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
