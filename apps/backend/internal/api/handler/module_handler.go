package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type ModuleHandler struct {
	svc *service.ModuleService
}

func NewModuleHandler(svc *service.ModuleService) *ModuleHandler {
	return &ModuleHandler{svc: svc}
}

// List retorna o catálogo de módulos da aplicação
// @Summary Listar módulos
// @Description Lista os módulos fixos da aplicação com nome e ordem de exibição (requer write em "permissions")
// @Tags modules (Módulos)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.Module
// @Security BearerAuth
// @Router /api/v1/{organization_id}/modules [get]
func (h *ModuleHandler) List(w http.ResponseWriter, r *http.Request) {
	modules, err := h.svc.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, modules, http.StatusOK)
}

type updateModuleRequest struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// Update renomeia/reordena um módulo existente. Não é possível criar um
// módulo novo por essa rota — a chave é fixa e amarrada às rotas em código.
// @Summary Atualizar módulo
// @Description Atualiza nome/ordem de exibição de um módulo existente (requer write em "permissions")
// @Tags modules (Módulos)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param key path string true "Chave do Módulo"
// @Param module body updateModuleRequest true "Dados do módulo"
// @Success 200 {object} entity.Module
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/modules/{key} [put]
func (h *ModuleHandler) Update(w http.ResponseWriter, r *http.Request) {
	key := entity.ModuleKey(r.PathValue("key"))

	var req updateModuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	module, err := h.svc.UpdateMeta(key, req.Name, req.Order)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, module, http.StatusOK)
}
