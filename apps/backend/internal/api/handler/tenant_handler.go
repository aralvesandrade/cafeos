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

// List returns all tenants
// @Summary List tenants
// @Description List all tenants (platform_owner only)
// @Tags tenants
// @Produce json
// @Success 200 {array} entity.Tenant
// @Security BearerAuth
// @Router /api/v1/tenants [get]
func (h *TenantHandler) List(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.repo.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, tenants, http.StatusOK)
}

// Create registers a new tenant
// @Summary Create a tenant
// @Description Create a new tenant (platform_owner only)
// @Tags tenants
// @Accept json
// @Produce json
// @Param tenant body createTenantRequest true "Tenant data"
// @Success 201 {object} entity.Tenant
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tenants [post]
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

// GetByID returns a tenant by its ID
// @Summary Get tenant by ID
// @Description Returns a single tenant (platform_owner only)
// @Tags tenants
// @Produce json
// @Param id path string true "Tenant ID"
// @Success 200 {object} entity.Tenant
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tenants/{id} [get]
func (h *TenantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	tenant, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "tenant not found", http.StatusNotFound)
		return
	}
	writeJSON(w, tenant, http.StatusOK)
}

// Update updates an existing tenant
// @Summary Update a tenant
// @Description Update tenant data (platform_owner only)
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path string true "Tenant ID"
// @Param tenant body updateTenantRequest true "Updated tenant data"
// @Success 200 {object} entity.Tenant
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tenants/{id} [put]
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

// Delete removes a tenant
// @Summary Delete a tenant
// @Description Delete a tenant by ID (platform_owner only)
// @Tags tenants
// @Param id path string true "Tenant ID"
// @Success 204 "No Content"
// @Security BearerAuth
// @Router /api/v1/tenants/{id} [delete]
func (h *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
