package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
	"github.com/google/uuid"
)

type OrganizationHandler struct {
	repo    repository.OrganizationRepository
	permSvc *service.PermissionService
}

func NewOrganizationHandler(repo repository.OrganizationRepository, permSvc *service.PermissionService) *OrganizationHandler {
	return &OrganizationHandler{repo: repo, permSvc: permSvc}
}

type createOrganizationRequest struct {
	Name      string `json:"name"`
	BrandName string `json:"brand_name"`
}

type updateOrganizationRequest struct {
	Name      *string `json:"name"`
	BrandName *string `json:"brand_name"`
	Status    *string `json:"status"`
}

// List retorna todas as organizações
// @Summary Listar organizações
// @Description Lista todas as organizações (somente platform_owner)
// @Tags organizations (Organizações)
// @Produce json
// @Success 200 {array} entity.Organization
// @Security BearerAuth
// @Router /api/v1/admin/organizations [get]
func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	organizations, err := h.repo.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, organizations, http.StatusOK)
}

// Create registra uma nova organização
// @Summary Criar organização
// @Description Cria uma nova organização (somente platform_owner)
// @Tags organizations (Organizações)
// @Accept json
// @Produce json
// @Param organization body createOrganizationRequest true "Dados da organização"
// @Success 201 {object} entity.Organization
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/organizations [post]
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		writeError(w, "organization name is required", http.StatusBadRequest)
		return
	}

	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	organization := &entity.Organization{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Slug:         slug,
		BrandName:    req.BrandName,
		Status:       "active",
		PrimaryColor: "#2E7D32",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.Create(organization); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.permSvc.SeedDefaults(organization.ID); err != nil {
		slog.Error("failed to seed default permissions for new organization", "organization_id", organization.ID, "error", err)
	}

	writeJSON(w, organization, http.StatusCreated)
}

// GetByID retorna uma organização pelo seu ID
// @Summary Obter organização por ID
// @Description Retorna uma única organização (somente platform_owner)
// @Tags organizations (Organizações)
// @Produce json
// @Param id path string true "ID da Organização"
// @Success 200 {object} entity.Organization
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/organizations/{id} [get]
func (h *OrganizationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	organization, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "organization not found", http.StatusNotFound)
		return
	}
	writeJSON(w, organization, http.StatusOK)
}

// Update atualiza uma organização existente
// @Summary Atualizar organização
// @Description Atualiza dados da organização (somente platform_owner)
// @Tags organizations (Organizações)
// @Accept json
// @Produce json
// @Param id path string true "ID da Organização"
// @Param organization body updateOrganizationRequest true "Dados atualizados da organização"
// @Success 200 {object} entity.Organization
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/organizations/{id} [put]
func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "organization not found", http.StatusNotFound)
		return
	}

	var req updateOrganizationRequest
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

// Delete remove uma organização
// @Summary Excluir organização
// @Description Exclui uma organização por ID (somente platform_owner)
// @Tags organizations (Organizações)
// @Param id path string true "ID da Organização"
// @Success 204 "No Content"
// @Security BearerAuth
// @Router /api/v1/admin/organizations/{id} [delete]
func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
