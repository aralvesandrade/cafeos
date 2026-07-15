package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type RoleHandler struct {
	svc *service.RoleService
}

func NewRoleHandler(svc *service.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

// List retorna o catálogo global de papéis (papéis de sistema + demais)
// @Summary Listar papéis
// @Description Lista o catálogo global de papéis, compartilhado por todas as organizações (requer write em "permissions")
// @Tags roles (Papéis)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.Role
// @Security BearerAuth
// @Router /api/v1/{organization_id}/roles [get]
func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, roles, http.StatusOK)
}

type createRoleRequest struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// Create cria um novo papel no catálogo global
// @Summary Criar papel
// @Description Cria um papel no catálogo global, compartilhado por todas as organizações (requer write em "permissions")
// @Tags roles (Papéis)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param role body createRoleRequest true "Dados do papel"
// @Success 201 {object} entity.Role
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/roles [post]
func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	role, err := h.svc.Create(req.Key, req.Name)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, role, http.StatusCreated)
}

type updateRoleRequest struct {
	Name string `json:"name"`
}

// Update renomeia um papel do catálogo global
// @Summary Atualizar papel
// @Description Renomeia um papel do catálogo global; papéis de sistema não podem ser alterados (requer write em "permissions")
// @Tags roles (Papéis)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Papel"
// @Param role body updateRoleRequest true "Dados do papel"
// @Success 200 {object} entity.Role
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/roles/{id} [put]
func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	role, err := h.svc.Update(id, req.Name)
	if err != nil {
		writeError(w, err.Error(), roleErrStatus(err))
		return
	}
	writeJSON(w, role, http.StatusOK)
}

// Delete remove um papel do catálogo global
// @Summary Excluir papel
// @Description Exclui um papel do catálogo global, se não estiver em uso por nenhum usuário (requer write em "permissions")
// @Tags roles (Papéis)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Papel"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/roles/{id} [delete]
func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), roleErrStatus(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func roleErrStatus(err error) int {
	switch {
	case errors.Is(err, service.ErrRoleNotFound):
		return http.StatusNotFound
	case errors.Is(err, service.ErrRoleIsSystem), errors.Is(err, service.ErrRoleInUse):
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
