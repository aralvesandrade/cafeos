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

func (h *StockHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	items, err := h.svc.ListItems(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, items, http.StatusOK)
}

func (h *StockHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := h.svc.GetItemByID(id)
	if err != nil {
		writeError(w, "item not found", http.StatusNotFound)
		return
	}
	writeJSON(w, item, http.StatusOK)
}

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

func (h *StockHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	movs, err := h.svc.ListMovements(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, movs, http.StatusOK)
}
