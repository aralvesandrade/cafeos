package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type PlanHandler struct {
	svc *service.PlanService
}

func NewPlanHandler(svc *service.PlanService) *PlanHandler {
	return &PlanHandler{svc: svc}
}

type createPlanRequest struct {
	Name            string                 `json:"name"`
	Slug            string                 `json:"slug"`
	Description     string                 `json:"description"`
	PriceCents      int64                  `json:"price_cents"`
	BillingInterval string                 `json:"billing_interval"`
	MaxFarms        int                    `json:"max_farms"`
	MaxUsers        int                    `json:"max_users"`
	Features        entity.PlanFeatureList `json:"features"`
	Active          *bool                  `json:"active"`
	Featured        bool                   `json:"featured"`
	DisplayOrder    int                    `json:"display_order"`
}

type updatePlanRequest struct {
	Name            *string                 `json:"name"`
	Slug            *string                 `json:"slug"`
	Description     *string                 `json:"description"`
	PriceCents      *int64                  `json:"price_cents"`
	BillingInterval *string                 `json:"billing_interval"`
	MaxFarms        *int                    `json:"max_farms"`
	MaxUsers        *int                    `json:"max_users"`
	Features        *entity.PlanFeatureList `json:"features"`
	Active          *bool                   `json:"active"`
	Featured        *bool                   `json:"featured"`
	DisplayOrder    *int                    `json:"display_order"`
}

// List retorna todos os planos
// @Summary Listar planos
// @Description Lista todos os planos de assinatura (somente platform_owner)
// @Tags plans (Planos)
// @Produce json
// @Success 200 {array} entity.Plan
// @Security BearerAuth
// @Router /api/v1/admin/plans [get]
func (h *PlanHandler) List(w http.ResponseWriter, r *http.Request) {
	plans, err := h.svc.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, plans, http.StatusOK)
}

// ListPublic retorna os planos ativos, sem autenticação
// @Summary Listar planos públicos
// @Description Lista os planos de assinatura ativos (público, usado na landing page)
// @Tags plans (Planos)
// @Produce json
// @Success 200 {array} entity.Plan
// @Router /api/v1/public/plans [get]
func (h *PlanHandler) ListPublic(w http.ResponseWriter, r *http.Request) {
	plans, err := h.svc.ListActive()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, plans, http.StatusOK)
}

// Create cria um novo plano
// @Summary Criar plano
// @Description Cria um novo plano de assinatura (somente platform_owner)
// @Tags plans (Planos)
// @Accept json
// @Produce json
// @Param plan body createPlanRequest true "Dados do plano"
// @Success 201 {object} entity.Plan
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/plans [post]
func (h *PlanHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}
	billingInterval := req.BillingInterval
	if billingInterval == "" {
		billingInterval = "monthly"
	}

	plan := &entity.Plan{
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		PriceCents:      req.PriceCents,
		BillingInterval: billingInterval,
		MaxFarms:        req.MaxFarms,
		MaxUsers:        req.MaxUsers,
		Features:        req.Features,
		Active:          active,
		Featured:        req.Featured,
		DisplayOrder:    req.DisplayOrder,
	}

	if err := h.svc.Create(plan); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, plan, http.StatusCreated)
}

// GetByID retorna um plano pelo seu ID
// @Summary Obter plano por ID
// @Description Retorna um único plano (somente platform_owner)
// @Tags plans (Planos)
// @Produce json
// @Param id path string true "ID do Plano"
// @Success 200 {object} entity.Plan
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/plans/{id} [get]
func (h *PlanHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	plan, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plan not found", http.StatusNotFound)
		return
	}
	writeJSON(w, plan, http.StatusOK)
}

// Update atualiza um plano existente
// @Summary Atualizar plano
// @Description Atualiza dados do plano (somente platform_owner)
// @Tags plans (Planos)
// @Accept json
// @Produce json
// @Param id path string true "ID do Plano"
// @Param plan body updatePlanRequest true "Dados atualizados do plano"
// @Success 200 {object} entity.Plan
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/plans/{id} [put]
func (h *PlanHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "plan not found", http.StatusNotFound)
		return
	}

	var req updatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Slug != nil {
		existing.Slug = *req.Slug
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.PriceCents != nil {
		existing.PriceCents = *req.PriceCents
	}
	if req.BillingInterval != nil {
		existing.BillingInterval = *req.BillingInterval
	}
	if req.MaxFarms != nil {
		existing.MaxFarms = *req.MaxFarms
	}
	if req.MaxUsers != nil {
		existing.MaxUsers = *req.MaxUsers
	}
	if req.Features != nil {
		existing.Features = *req.Features
	}
	if req.Active != nil {
		existing.Active = *req.Active
	}
	if req.Featured != nil {
		existing.Featured = *req.Featured
	}
	if req.DisplayOrder != nil {
		existing.DisplayOrder = *req.DisplayOrder
	}

	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, existing, http.StatusOK)
}

// Delete remove um plano
// @Summary Excluir plano
// @Description Exclui um plano por ID (somente platform_owner)
// @Tags plans (Planos)
// @Param id path string true "ID do Plano"
// @Success 204 "No Content"
// @Security BearerAuth
// @Router /api/v1/admin/plans/{id} [delete]
func (h *PlanHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
