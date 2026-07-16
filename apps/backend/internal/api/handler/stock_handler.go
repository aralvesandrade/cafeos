package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type StockHandler struct {
	svc     *service.StockService
	farmSvc *service.FarmService
}

func NewStockHandler(svc *service.StockService, farmSvc *service.FarmService) *StockHandler {
	return &StockHandler{svc: svc, farmSvc: farmSvc}
}

type createStockItemRequest struct {
	ProductID  string  `json:"product_id"`
	FarmID     *string `json:"farm_id"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	Batch      string  `json:"batch"`
	ExpiryDate string  `json:"expiry_date"`
	MinStock   float64 `json:"min_stock"`
	Location   string  `json:"location"`
	Notes      string  `json:"notes"`
}

// canAccessStockItem mirrors FinancialHandler.canAccessFinancial: proprietario
// may only see/edit items linked to farms it owns, or items with no farm link.
func (h *StockHandler) canAccessStockItem(r *http.Request, organizationID string, item *entity.StockItem) bool {
	if item.FarmID == nil {
		return true
	}
	userID, restricted := restrictedOwnerID(r)
	if !restricted {
		return true
	}
	owned, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
	if err != nil {
		return false
	}
	return owned[*item.FarmID]
}

// CreateItem registra um novo item de estoque para a organização autenticada
// @Summary Criar item de estoque
// @Description Registra um novo item de estoque na organização
// @Tags stock (Estoque)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param item body createStockItemRequest true "Dados do item de estoque"
// @Success 201 {object} entity.StockItem
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/items [post]
func (h *StockHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
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
	item, err := h.svc.CreateItem(organizationID, req.ProductID, req.Unit, req.Batch, req.Location, req.Notes, req.FarmID, req.Quantity, req.MinStock, expiry)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, item, http.StatusCreated)
}

// ListItems retorna todos os itens de estoque da organização autenticada
// @Summary Listar itens de estoque
// @Description Lista todos os itens de estoque pertencentes à organização
// @Tags stock (Estoque)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Success 200 {array} entity.StockItem
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/items [get]
func (h *StockHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	items, err := h.svc.ListItems(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userID, restricted := restrictedOwnerID(r); restricted {
		owned, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filtered := items[:0]
		for _, it := range items {
			if it.FarmID == nil || owned[*it.FarmID] {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}

	if farmID := r.URL.Query().Get("farm_id"); farmID != "" {
		filtered := items[:0]
		for _, it := range items {
			if it.FarmID != nil && *it.FarmID == farmID {
				filtered = append(filtered, it)
			}
		}
		items = filtered
	}

	writeJSON(w, items, http.StatusOK)
}

// GetItemByID retorna um item de estoque pelo seu ID
// @Summary Obter item de estoque por ID
// @Description Retorna um único item de estoque
// @Tags stock (Estoque)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Item de Estoque"
// @Success 200 {object} entity.StockItem
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/items/{id} [get]
func (h *StockHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	item, err := h.svc.GetItemByID(id)
	if err != nil {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	if !h.canAccessStockItem(r, organizationID, item) {
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
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Item de Estoque"
// @Param item body entity.StockItem true "Dados atualizados do item de estoque"
// @Success 200 {object} entity.StockItem
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/items/{id} [put]
func (h *StockHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	existing, err := h.svc.GetItemByID(id)
	if err != nil {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	if !h.canAccessStockItem(r, organizationID, existing) {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	var req struct {
		Quantity *float64 `json:"quantity"`
		Location *string  `json:"location"`
		MinStock *float64 `json:"min_stock"`
		FarmID   *string  `json:"farm_id"`
		Notes    *string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Quantity != nil {
		existing.Quantity = *req.Quantity
	}
	if req.FarmID != nil {
		existing.FarmID = req.FarmID
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
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Item de Estoque"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/items/{id} [delete]
func (h *StockHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")
	existing, err := h.svc.GetItemByID(id)
	if err != nil {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	if !h.canAccessStockItem(r, organizationID, existing) {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
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
// @Param organization_id path string true "ID da Organização"
// @Param movement body recordMovementRequest true "Dados da movimentação de estoque"
// @Success 201 {object} entity.StockMovement
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/movements [post]
func (h *StockHandler) RecordMovement(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	var req recordMovementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", req.Date)
	mov, err := h.svc.RecordMovement(organizationID, req.ItemID, req.Type, req.Reference, req.Notes, req.Quantity, date)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, mov, http.StatusCreated)
}

// GetMovementByID retorna uma movimentação de estoque pelo seu ID
// @Summary Obter movimentação de estoque por ID
// @Description Retorna uma única movimentação de estoque
// @Tags stock (Estoque)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Movimentação"
// @Success 200 {object} entity.StockMovement
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/movements/{id} [get]
func (h *StockHandler) GetMovementByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	mov, err := h.svc.GetMovementByID(id)
	if err != nil {
		writeError(w, "movement not found", http.StatusNotFound)
		return
	}
	writeJSON(w, mov, http.StatusOK)
}

// ListMovements retorna todas as movimentações de estoque da organização autenticada
// @Summary Listar movimentações de estoque
// @Description Lista todas as movimentações de estoque pertencentes à organização
// @Tags stock (Estoque)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param farm_id query string false "Filtrar por ID da Fazenda"
// @Success 200 {array} entity.StockMovement
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/stock/movements [get]
func (h *StockHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	movs, err := h.svc.ListMovements(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, restricted := restrictedOwnerID(r)
	farmID := r.URL.Query().Get("farm_id")
	if restricted || farmID != "" {
		items, err := h.svc.ListItems(organizationID)
		if err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		itemFarm := make(map[string]*string, len(items))
		for _, it := range items {
			itemFarm[it.ID] = it.FarmID
		}

		if restricted {
			owned, err := h.farmSvc.OwnedFarmIDs(organizationID, userID)
			if err != nil {
				writeError(w, err.Error(), http.StatusInternalServerError)
				return
			}
			filtered := movs[:0]
			for _, m := range movs {
				fid := itemFarm[m.ItemID]
				if fid == nil || owned[*fid] {
					filtered = append(filtered, m)
				}
			}
			movs = filtered
		}

		if farmID != "" {
			filtered := movs[:0]
			for _, m := range movs {
				fid := itemFarm[m.ItemID]
				if fid != nil && *fid == farmID {
					filtered = append(filtered, m)
				}
			}
			movs = filtered
		}
	}

	writeJSON(w, movs, http.StatusOK)
}
