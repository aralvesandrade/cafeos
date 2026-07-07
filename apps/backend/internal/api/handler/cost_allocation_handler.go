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

func (h *CostAllocationHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

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
		TenantID:     tenantID,
		HarvestID:    req.HarvestID,
		CostCenterID: req.CostCenterID,
		Description:  req.Description,
		TotalAmount:  req.TotalAmount,
		Method:       entity.AllocationMethod(req.Method),
		Date:         date,
		Percentages:  req.Percentages,
	}

	allocation, err := h.svc.Create(input)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, allocation, http.StatusCreated)
}

func (h *CostAllocationHandler) ListByHarvest(w http.ResponseWriter, r *http.Request) {
	harvestID := r.PathValue("harvest_id")
	allocs, err := h.svc.ListByHarvest(harvestID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, allocs, http.StatusOK)
}

func (h *CostAllocationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "cost allocation not found", http.StatusNotFound)
		return
	}
	writeJSON(w, a, http.StatusOK)
}

func (h *CostAllocationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
