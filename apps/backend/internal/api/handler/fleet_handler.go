package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type FleetHandler struct {
	svc *service.FleetService
}

func NewFleetHandler(svc *service.FleetService) *FleetHandler {
	return &FleetHandler{svc: svc}
}

type createVehicleRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Plate string `json:"plate"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

// CreateVehicle registra um novo veículo para o tenant autenticado
// @Summary Criar veículo
// @Description Registra um novo veículo na frota do tenant
// @Tags fleet (Frota)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param vehicle body createVehicleRequest true "Dados do veículo"
// @Success 201 {object} entity.Vehicle
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/vehicles [post]
func (h *FleetHandler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req createVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	v, err := h.svc.CreateVehicle(tenantID, req.Name, req.Type, req.Plate, req.Brand, req.Model, req.Year)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, v, http.StatusCreated)
}

// ListVehicles retorna todos os veículos do tenant autenticado
// @Summary Listar veículos
// @Description Lista todos os veículos pertencentes à frota do tenant
// @Tags fleet (Frota)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.Vehicle
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/vehicles [get]
func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	vehicles, err := h.svc.ListVehicles(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, vehicles, http.StatusOK)
}

// GetVehicleByID retorna um veículo pelo seu ID
// @Summary Obter veículo por ID
// @Description Retorna um único veículo
// @Tags fleet (Frota)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Veículo"
// @Success 200 {object} entity.Vehicle
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/vehicles/{id} [get]
func (h *FleetHandler) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, err := h.svc.GetVehicleByID(id)
	if err != nil {
		writeError(w, "vehicle not found", http.StatusNotFound)
		return
	}
	writeJSON(w, v, http.StatusOK)
}

// UpdateVehicle atualiza um veículo existente
// @Summary Atualizar veículo
// @Description Atualiza dados do veículo por ID (atualização parcial - somente os campos informados são alterados)
// @Tags fleet (Frota)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Veículo"
// @Param vehicle body entity.Vehicle true "Dados atualizados do veículo"
// @Success 200 {object} entity.Vehicle
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/vehicles/{id} [put]
func (h *FleetHandler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.svc.GetVehicleByID(id)
	if err != nil {
		writeError(w, "vehicle not found", http.StatusNotFound)
		return
	}
	var req struct {
		Name   *string `json:"name"`
		Type   *string `json:"type"`
		Plate  *string `json:"plate"`
		Brand  *string `json:"brand"`
		Model  *string `json:"model"`
		Year   *int    `json:"year"`
		Status *string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Type != nil {
		existing.Type = entity.VehicleType(*req.Type)
	}
	if req.Plate != nil {
		existing.Plate = *req.Plate
	}
	if req.Brand != nil {
		existing.Brand = *req.Brand
	}
	if req.Model != nil {
		existing.Model = *req.Model
	}
	if req.Year != nil {
		existing.Year = *req.Year
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if err := h.svc.UpdateVehicle(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// DeleteVehicle remove um veículo pelo seu ID
// @Summary Excluir veículo
// @Description Exclui um veículo por ID
// @Tags fleet (Frota)
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID do Veículo"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/vehicles/{id} [delete]
func (h *FleetHandler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteVehicle(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type createMaintenanceRequest struct {
	VehicleID   string  `json:"vehicle_id"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	Odometer    float64 `json:"odometer"`
	Date        string  `json:"date"`
	Notes       string  `json:"notes"`
}

// CreateMaintenance registra um novo registro de manutenção para um veículo
// @Summary Criar registro de manutenção
// @Description Registra um novo registro de manutenção para um veículo
// @Tags fleet (Frota)
// @Accept json
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param maintenance body createMaintenanceRequest true "Dados da manutenção"
// @Success 201 {object} entity.Maintenance
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/maintenance [post]
func (h *FleetHandler) CreateMaintenance(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req createMaintenanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", req.Date)
	m, err := h.svc.CreateMaintenance(tenantID, req.VehicleID, req.Type, req.Description, req.Notes, req.Cost, req.Odometer, date)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, m, http.StatusCreated)
}

// ListMaintenance retorna todos os registros de manutenção do tenant autenticado
// @Summary Listar registros de manutenção
// @Description Lista todos os registros de manutenção pertencentes à frota do tenant
// @Tags fleet (Frota)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Success 200 {array} entity.Maintenance
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/maintenance [get]
func (h *FleetHandler) ListMaintenance(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	items, err := h.svc.ListMaintenance(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, items, http.StatusOK)
}

// GetMaintenanceByID retorna um registro de manutenção pelo seu ID
// @Summary Obter registro de manutenção por ID
// @Description Retorna um único registro de manutenção
// @Tags fleet (Frota)
// @Produce json
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Manutenção"
// @Success 200 {object} entity.Maintenance
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/maintenance/{id} [get]
func (h *FleetHandler) GetMaintenanceByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := h.svc.GetMaintenanceByID(id)
	if err != nil {
		writeError(w, "maintenance not found", http.StatusNotFound)
		return
	}
	writeJSON(w, m, http.StatusOK)
}

// DeleteMaintenance remove um registro de manutenção pelo seu ID
// @Summary Excluir registro de manutenção
// @Description Exclui um registro de manutenção por ID
// @Tags fleet (Frota)
// @Param tenant_id path string true "ID do Tenant"
// @Param id path string true "ID da Manutenção"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{tenant_id}/fleet/maintenance/{id} [delete]
func (h *FleetHandler) DeleteMaintenance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteMaintenance(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
