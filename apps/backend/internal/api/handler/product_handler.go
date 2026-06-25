package handler

import (
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
)

type ProductHandler struct {
	repo repository.AgriculturalProductRepository
}

func NewProductHandler(repo repository.AgriculturalProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	products, err := h.repo.ListByTenant(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, products, http.StatusOK)
}
