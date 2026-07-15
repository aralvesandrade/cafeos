package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type PermissionHandler struct {
	svc *service.PermissionService
}

func NewPermissionHandler(svc *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

// List retorna a matriz completa de permissões da organização
// @Summary Listar matriz de permissões
// @Description Lista o acesso configurado de cada papel em cada módulo (requer write em "permissions")
// @Tags permissions (Permissões)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.RolePermission
// @Security BearerAuth
// @Router /api/v1/{organization_id}/permissions [get]
func (h *PermissionHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	permissions, err := h.svc.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, permissions, http.StatusOK)
}

type updatePermissionEntry struct {
	RoleID string             `json:"role_id"`
	Module entity.ModuleKey   `json:"module"`
	Access entity.AccessLevel `json:"access"`
}

// Update sobrescreve o acesso de uma ou mais células da matriz
// @Summary Atualizar matriz de permissões
// @Description Atualiza o acesso (none|read|write) de papéis por módulo (requer write em "permissions")
// @Tags permissions (Permissões)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param permissions body []updatePermissionEntry true "Lista de alterações"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/permissions [put]
func (h *PermissionHandler) Update(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var entries []updatePermissionEntry
	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	for _, entry := range entries {
		if err := h.svc.Update(organizationID, entry.RoleID, entry.Module, entry.Access); err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	writeJSON(w, map[string]string{"status": "updated"}, http.StatusOK)
}

// Mine retorna o acesso do usuário autenticado em cada módulo
// @Summary Meus acessos
// @Description Retorna o mapa módulo -> acesso do papel do usuário autenticado
// @Tags permissions (Permissões)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/permissions/me [get]
func (h *PermissionHandler) Mine(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	result := make(map[entity.ModuleKey]entity.AccessLevel, len(entity.AllModules))
	for _, module := range entity.AllModules {
		access, err := h.svc.GetAccess(organizationID, role, module)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result[module] = access
	}

	writeJSON(w, result, http.StatusOK)
}
