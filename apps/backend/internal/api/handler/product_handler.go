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

// List retorna todos os produtos agrícolas da organização autenticada
// @Summary Listar produtos agrícolas
// @Description Lista todos os produtos agrícolas pertencentes à organização
// @Tags agricultural-products (Produtos Agrícolas)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.AgriculturalProduct
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/agricultural-products [get]
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	products, err := h.repo.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, products, http.StatusOK)
}
