package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type CostCenterHandler struct {
	svc *service.CostCenterService
}

func NewCostCenterHandler(svc *service.CostCenterService) *CostCenterHandler {
	return &CostCenterHandler{svc: svc}
}

type createCostCenterRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Type        string `json:"type"`
	CostGroup   string `json:"cost_group"`
	Description string `json:"description"`
}

func (h *CostCenterHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)

	var req createCostCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cc, err := h.svc.Create(tenantID, req.Name, req.Code, entity.CostCenterType(req.Type), entity.CostGroup(req.CostGroup), req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, cc, http.StatusCreated)
}

// SenarCategories returns the fixed SENAR/CEPEA despesa category catalog,
// used by clients to pre-fill new cost centers with a known cost_group.
func (h *CostCenterHandler) SenarCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.svc.SenarCategories(), http.StatusOK)
}

func (h *CostCenterHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	ccs, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ccs, http.StatusOK)
}

func (h *CostCenterHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cc, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "cost center not found", http.StatusNotFound)
		return
	}
	writeJSON(w, cc, http.StatusOK)
}

func (h *CostCenterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "cost center not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Code        *string `json:"code"`
		Type        *string `json:"type"`
		CostGroup   *string `json:"cost_group"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Code != nil {
		existing.Code = *input.Code
	}
	if input.Type != nil {
		existing.Type = entity.CostCenterType(*input.Type)
	}
	if input.CostGroup != nil {
		existing.CostGroup = entity.CostGroup(*input.CostGroup)
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}

	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, existing, http.StatusOK)
}

func (h *CostCenterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
