package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type StockHandler struct {
	svc *service.StockService
}

func NewStockHandler(svc *service.StockService) *StockHandler {
	return &StockHandler{svc: svc}
}

type createStockItemRequest struct {
	ProductID  string  `json:"product_id"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	Batch      string  `json:"batch"`
	ExpiryDate string  `json:"expiry_date"`
	MinStock   float64 `json:"min_stock"`
	Location   string  `json:"location"`
	Notes      string  `json:"notes"`
}

// CreateItem registra um novo item de estoque para o tenant autenticado
// @Summary Criar item de estoque
// @Description Registra um novo item de estoque no tenant
// @Tags stock (Estoque)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param item body createStockItemRequest true "Dados do item de estoque"
// @Success 201 {object} entity.StockItem
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/items [post]
func (h *StockHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req createStockItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	var expiry *time.Time
	if req.ExpiryDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpiryDate)
		if err == nil {
			expiry = &t
		}
	}
	item, err := h.svc.CreateItem(tenantID, req.ProductID, req.Unit, req.Batch, req.Location, req.Notes, req.Quantity, req.MinStock, expiry)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, item, http.StatusCreated)
}

// ListItems retorna todos os itens de estoque do tenant autenticado
// @Summary Listar itens de estoque
// @Description Lista todos os itens de estoque pertencentes ao tenant
// @Tags stock (Estoque)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.StockItem
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/items [get]
func (h *StockHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	items, err := h.svc.ListItems(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, items, http.StatusOK)
}

// GetItemByID retorna um item de estoque pelo seu ID
// @Summary Obter item de estoque por ID
// @Description Retorna um único item de estoque
// @Tags stock (Estoque)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Item de Estoque"
// @Success 200 {object} entity.StockItem
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/items/{id} [get]
func (h *StockHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.svc.GetItemByID(id)
	if err != nil {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	writeJSON(w, item, http.StatusOK)
}

// UpdateItem atualiza um item de estoque existente
// @Summary Atualizar item de estoque
// @Description Atualiza dados do item de estoque por ID (atualização parcial - somente os campos informados são alterados)
// @Tags stock (Estoque)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Item de Estoque"
// @Param item body entity.StockItem true "Dados atualizados do item de estoque"
// @Success 200 {object} entity.StockItem
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/items/{id} [put]
func (h *StockHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.svc.GetItemByID(id)
	if err != nil {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	var req struct {
		Quantity *float64 `json:"quantity"`
		Location *string  `json:"location"`
		MinStock *float64 `json:"min_stock"`
		Notes    *string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Quantity != nil {
		existing.Quantity = *req.Quantity
	}
	if req.Location != nil {
		existing.Location = *req.Location
	}
	if req.MinStock != nil {
		existing.MinStock = *req.MinStock
	}
	if req.Notes != nil {
		existing.Notes = *req.Notes
	}
	if err := h.svc.UpdateItem(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// DeleteItem remove um item de estoque pelo seu ID
// @Summary Excluir item de estoque
// @Description Exclui um item de estoque por ID
// @Tags stock (Estoque)
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Item de Estoque"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/items/{id} [delete]
func (h *StockHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteItem(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type recordMovementRequest struct {
	ItemID    string  `json:"item_id"`
	Type      string  `json:"type"`
	Quantity  float64 `json:"quantity"`
	Date      string  `json:"date"`
	Reference string  `json:"reference"`
	Notes     string  `json:"notes"`
}

// RecordMovement registra uma movimentação de estoque (entrada/saída) para um item
// @Summary Registrar movimentação de estoque
// @Description Registra uma movimentação de estoque (entrada/saída) para um item de estoque
// @Tags stock (Estoque)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param movement body recordMovementRequest true "Dados da movimentação de estoque"
// @Success 201 {object} entity.StockMovement
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/movements [post]
func (h *StockHandler) RecordMovement(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req recordMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", req.Date)
	mov, err := h.svc.RecordMovement(tenantID, req.ItemID, req.Type, req.Reference, req.Notes, req.Quantity, date)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, mov, http.StatusCreated)
}

// ListMovements retorna todas as movimentações de estoque do tenant autenticado
// @Summary Listar movimentações de estoque
// @Description Lista todas as movimentações de estoque pertencentes ao tenant
// @Tags stock (Estoque)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.StockMovement
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/stock/movements [get]
func (h *StockHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	movs, err := h.svc.ListMovements(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, movs, http.StatusOK)
}
