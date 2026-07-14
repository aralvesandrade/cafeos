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

// Create registra um novo centro de custo para o tenant autenticado
// @Summary Criar centro de custo
// @Description Registra um novo centro de custo no tenant
// @Tags cost-centers (Centros de Custo)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param cost_center body createCostCenterRequest true "Dados do centro de custo"
// @Success 201 {object} entity.CostCenter
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/cost-centers [post]
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

// SenarCategories retorna o catálogo fixo de categorias de despesa SENAR/CEPEA,
// usado pelos clientes para pré-preencher novos centros de custo com um cost_group conhecido.
// @Summary Listar categorias SENAR/CEPEA
// @Description Retorna o catálogo fixo de categorias de despesa SENAR/CEPEA
// @Tags cost-centers (Centros de Custo)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.SenarCostCategory
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/cost-centers/senar-categories [get]
func (h *CostCenterHandler) SenarCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, h.svc.SenarCategories(), http.StatusOK)
}

// List retorna todos os centros de custo do tenant autenticado
// @Summary Listar centros de custo
// @Description Lista todos os centros de custo pertencentes ao tenant
// @Tags cost-centers (Centros de Custo)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.CostCenter
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/cost-centers [get]
func (h *CostCenterHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	ccs, err := h.svc.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ccs, http.StatusOK)
}

// GetByID retorna um centro de custo pelo seu ID
// @Summary Obter centro de custo por ID
// @Description Retorna um único centro de custo
// @Tags cost-centers (Centros de Custo)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Centro de Custo"
// @Success 200 {object} entity.CostCenter
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/cost-centers/{id} [get]
func (h *CostCenterHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cc, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "cost center not found", http.StatusNotFound)
		return
	}
	writeJSON(w, cc, http.StatusOK)
}

// Update atualiza um centro de custo existente
// @Summary Atualizar centro de custo
// @Description Atualiza dados do centro de custo por ID (atualização parcial - somente os campos informados são alterados)
// @Tags cost-centers (Centros de Custo)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Centro de Custo"
// @Param cost_center body entity.CostCenter true "Dados atualizados do centro de custo"
// @Success 200 {object} entity.CostCenter
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/cost-centers/{id} [put]
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

// Delete remove um centro de custo pelo seu ID
// @Summary Excluir centro de custo
// @Description Exclui um centro de custo por ID
// @Tags cost-centers (Centros de Custo)
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Centro de Custo"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/cost-centers/{id} [delete]
func (h *CostCenterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
