package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type OperationTypeHandler struct {
	svc *service.OperationTypeService
}

func NewOperationTypeHandler(svc *service.OperationTypeService) *OperationTypeHandler {
	return &OperationTypeHandler{svc: svc}
}

type createOperationTypeRequest struct {
	Name  string `json:"name"`
	Code  string `json:"code"`
	Color string `json:"color"`
}

// Create registra um novo tipo de operação para a organização autenticada
// @Summary Criar tipo de operação
// @Description Registra um novo tipo de operação agrícola na organização
// @Tags operation-types (Tipos de Operação)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param operation_type body createOperationTypeRequest true "Dados do tipo de operação"
// @Success 201 {object} entity.OperationType
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operation-types [post]
func (h *OperationTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createOperationTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ot, err := h.svc.Create(organizationID, req.Name, req.Code, req.Color)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, ot, http.StatusCreated)
}

// List retorna todos os tipos de operação da organização autenticada
// @Summary Listar tipos de operação
// @Description Lista todos os tipos de operação pertencentes à organização
// @Tags operation-types (Tipos de Operação)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.OperationType
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operation-types [get]
func (h *OperationTypeHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	ots, err := h.svc.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ots, http.StatusOK)
}

// GetByID retorna um tipo de operação pelo seu ID
// @Summary Obter tipo de operação por ID
// @Description Retorna um único tipo de operação
// @Tags operation-types (Tipos de Operação)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Tipo de Operação"
// @Success 200 {object} entity.OperationType
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operation-types/{id} [get]
func (h *OperationTypeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ot, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "operation type not found", http.StatusNotFound)
		return
	}
	writeJSON(w, ot, http.StatusOK)
}

// Update atualiza um tipo de operação existente
// @Summary Atualizar tipo de operação
// @Description Atualiza dados do tipo de operação por ID (atualização parcial - somente os campos informados são alterados)
// @Tags operation-types (Tipos de Operação)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Tipo de Operação"
// @Param operation_type body createOperationTypeRequest true "Dados atualizados do tipo de operação"
// @Success 200 {object} entity.OperationType
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operation-types/{id} [put]
func (h *OperationTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "operation type not found", http.StatusNotFound)
		return
	}

	var input struct {
		Name  *string `json:"name"`
		Code  *string `json:"code"`
		Color *string `json:"color"`
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
	if input.Color != nil {
		existing.Color = *input.Color
	}

	if err := h.svc.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, existing, http.StatusOK)
}

// Delete remove um tipo de operação pelo seu ID
// @Summary Excluir tipo de operação
// @Description Exclui um tipo de operação por ID
// @Tags operation-types (Tipos de Operação)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Tipo de Operação"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operation-types/{id} [delete]
func (h *OperationTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
