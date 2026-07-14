package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type OperationHandler struct {
	svc     *service.OperationService
	plotSvc *service.PlotService
	farmSvc *service.FarmService
}

func NewOperationHandler(svc *service.OperationService, plotSvc *service.PlotService, farmSvc *service.FarmService) *OperationHandler {
	return &OperationHandler{svc: svc, plotSvc: plotSvc, farmSvc: farmSvc}
}

// canAccessOperation reports whether the request may access an operation:
// proprietario is restricted to operations on plots of farms it owns.
func (h *OperationHandler) canAccessOperation(r *http.Request, organizationID string, op *entity.Operation) bool {
	userID, restricted := restrictedOwnerID(r)
	if !restricted {
		return true
	}
	ownedFarms, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
	if err != nil {
		return false
	}
	owned, err := h.plotSvc.PlotIDsForFarms(organizationID, ownedFarms)
	if err != nil {
		return false
	}
	return owned[op.PlotID]
}

type createOperationRequest struct {
	PlotID       string  `json:"plot_id"`
	HarvestID    *string `json:"harvest_id"`
	CostCenterID *string `json:"cost_center_id"`
	Type         string  `json:"type"`
	Date         string  `json:"date"`
	Responsible  string  `json:"responsible"`
	ProductUsed  string  `json:"product_used"`
	Quantity     float64 `json:"quantity"`
	Cost         float64 `json:"cost"`
	Notes        string  `json:"notes"`
}

// Create registra uma nova operação agrícola
// @Summary Registrar operação
// @Description Registra uma operação agrícola (adubação, pulverização, irrigação, poda, colheita)
// @Tags operations (Operações)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param operation body createOperationRequest true "Dados da operação"
// @Success 201 {object} SwaggerOperation
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations [post]
func (h *OperationHandler) Create(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)

	var req createOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		date = time.Now()
	}

	op, err := h.svc.Create(organizationID, req.PlotID, entity.OperationType(req.Type), date, req.Responsible, req.ProductUsed, req.Quantity, req.Cost, req.Notes, req.HarvestID, req.CostCenterID)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, op, http.StatusCreated)
}

// GetByID retorna uma operação pelo seu ID
// @Summary Obter operação por ID
// @Description Retorna uma única operação
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Operação"
// @Success 200 {object} SwaggerOperation
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations/{id} [get]
func (h *OperationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	op, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "operation not found", http.StatusNotFound)
		return
	}
	if !h.canAccessOperation(r, organizationID, op) {
		writeError(w, "operation not found", http.StatusNotFound)
		return
	}
	writeJSON(w, op, http.StatusOK)
}

// ListByPlot retorna todas as operações de um talhão
// @Summary Listar operações por talhão
// @Description Lista todas as operações de um talhão específico
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param plot_id path string true "ID do Talhão"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/plots/{plot_id}/operations [get]
func (h *OperationHandler) ListByPlot(w http.ResponseWriter, r *http.Request) {
	plotID := r.PathValue("plot_id")
	ops, err := h.svc.ListByPlot(plotID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, ops, http.StatusOK)
}

// List retorna todas as operações da organização autenticada
// @Summary Listar todas as operações
// @Description Lista todas as operações de todos os talhões da organização
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Router /api/v1/{organization_id}/operations [get]
func (h *OperationHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	ops, err := h.svc.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ops, err = h.filterOperations(r, organizationID, ops)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, ops, http.StatusOK)
}

// filterOperations applies proprietario ownership scoping and the optional
// ?farm_id= query filter, in that order, to a list of operations.
func (h *OperationHandler) filterOperations(r *http.Request, organizationID string, ops []*entity.Operation) ([]*entity.Operation, error) {
	if userID, restricted := restrictedOwnerID(r); restricted {
		ownedFarms, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
		if err != nil {
			return nil, err
		}
		owned, err := h.plotSvc.PlotIDsForFarms(organizationID, ownedFarms)
		if err != nil {
			return nil, err
		}
		filtered := ops[:0]
		for _, op := range ops {
			if owned[op.PlotID] {
				filtered = append(filtered, op)
			}
		}
		ops = filtered
	}

	if farmID := r.URL.Query().Get("farm_id"); farmID != "" {
		plotIDs, err := h.plotSvc.PlotIDsForFarm(organizationID, farmID)
		if err != nil {
			return nil, err
		}
		filtered := ops[:0]
		for _, op := range ops {
			if plotIDs[op.PlotID] {
				filtered = append(filtered, op)
			}
		}
		ops = filtered
	}

	return ops, nil
}

// ListRecent retorna as operações mais recentes
// @Summary Listar operações recentes
// @Description Lista as operações mais recentes, limitadas por parâmetro de consulta
// @Tags operations (Operações)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param limit query int false "Máximo de resultados (padrão 10)"
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Success 200 {array} SwaggerOperation
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations/recent [get]
func (h *OperationHandler) ListRecent(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	ops, err := h.svc.ListRecent(organizationID, limit)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ops, err = h.filterOperations(r, organizationID, ops)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, ops, http.StatusOK)
}

// Delete remove uma operação
// @Summary Excluir operação
// @Description Exclui uma operação por ID
// @Tags operations (Operações)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Operação"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/operations/{id} [delete]
func (h *OperationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	op, err := h.svc.GetByID(id)
	if err != nil {
		writeError(w, "operation not found", http.StatusNotFound)
		return
	}
	if !h.canAccessOperation(r, organizationID, op) {
		writeError(w, "operation not found", http.StatusNotFound)
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
