package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
)

type TenantHandler struct {
	repo repository.TenantRepository
}

func NewTenantHandler(repo repository.TenantRepository) *TenantHandler {
	return &TenantHandler{repo: repo}
}

type createTenantRequest struct {
	Name      string `json:"name"`
	BrandName string `json:"brand_name"`
	Plan      string `json:"plan"`
}

type updateTenantRequest struct {
	Name      *string `json:"name"`
	BrandName *string `json:"brand_name"`
	Plan      *string `json:"plan"`
	Status    *string `json:"status"`
}

// List retorna todos os tenants
// @Summary Listar tenants
// @Description Lista todos os tenants (somente platform_owner)
// @Tags tenants (Inquilinos)
// @Produce json
// @Success 200 {array} entity.Tenant
// @Security BearerAuth
// @Router /api/v1/admin/tenants [get]
func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.repo.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, tenants, http.StatusOK)
}

// Create registra um novo tenant
// @Summary Criar tenant
// @Description Cria um novo tenant (somente platform_owner)
// @Tags tenants (Inquilinos)
// @Accept json
// @Produce json
// @Param tenant body createTenantRequest true "Dados do tenant"
// @Success 201 {object} entity.Tenant
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/tenants [post]
func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		writeError(w, "tenant name is required", http.StatusBadRequest)
		return
	}

	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	tenant := &entity.Tenant{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Slug:         slug,
		BrandName:    req.BrandName,
		Plan:         req.Plan,
		Status:       "active",
		PrimaryColor: "#2E7D32",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.Create(tenant); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, tenant, http.StatusCreated)
}

// GetByID retorna um tenant pelo seu ID
// @Summary Obter tenant por ID
// @Description Retorna um único tenant (somente platform_owner)
// @Tags tenants (Inquilinos)
// @Produce json
// @Param id path string true "ID do Tenant"
// @Success 200 {object} entity.Tenant
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/tenants/{id} [get]
func (h *TenantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tenant, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "tenant not found", http.StatusNotFound)
		return
	}
	writeJSON(w, tenant, http.StatusOK)
}

// Update atualiza um tenant existente
// @Summary Atualizar tenant
// @Description Atualiza dados do tenant (somente platform_owner)
// @Tags tenants (Inquilinos)
// @Accept json
// @Produce json
// @Param id path string true "ID do Tenant"
// @Param tenant body updateTenantRequest true "Dados atualizados do tenant"
// @Success 200 {object} entity.Tenant
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/tenants/{id} [put]
func (h *TenantHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "tenant not found", http.StatusNotFound)
		return
	}

	var req updateTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.BrandName != nil {
		existing.BrandName = *req.BrandName
	}
	if req.Plan != nil {
		existing.Plan = *req.Plan
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, existing, http.StatusOK)
}

// Delete remove um tenant
// @Summary Excluir tenant
// @Description Exclui um tenant por ID (somente platform_owner)
// @Tags tenants (Inquilinos)
// @Param id path string true "ID do Tenant"
// @Success 204 "No Content"
// @Security BearerAuth
// @Router /api/v1/admin/tenants/{id} [delete]
func (h *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
