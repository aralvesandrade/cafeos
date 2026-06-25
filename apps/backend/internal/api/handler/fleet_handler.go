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

func (h *FleetHandler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	vehicles, err := h.svc.ListVehicles(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, vehicles, http.StatusOK)
}

func (h *FleetHandler) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, err := h.svc.GetVehicleByID(id)
	if err != nil {
		writeError(w, "vehicle not found", http.StatusNotFound)
		return
	}
	writeJSON(w, v, http.StatusOK)
}

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
	if req.Name != nil { existing.Name = *req.Name }
	if req.Type != nil { existing.Type = entity.VehicleType(*req.Type) }
	if req.Plate != nil { existing.Plate = *req.Plate }
	if req.Brand != nil { existing.Brand = *req.Brand }
	if req.Model != nil { existing.Model = *req.Model }
	if req.Year != nil { existing.Year = *req.Year }
	if req.Status != nil { existing.Status = *req.Status }
	if err := h.svc.UpdateVehicle(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

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

func (h *FleetHandler) ListMaintenance(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	items, err := h.svc.ListMaintenance(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, items, http.StatusOK)
}

func (h *FleetHandler) GetMaintenanceByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	m, err := h.svc.GetMaintenanceByID(id)
	if err != nil {
		writeError(w, "maintenance not found", http.StatusNotFound)
		return
	}
	writeJSON(w, m, http.StatusOK)
}

func (h *FleetHandler) DeleteMaintenance(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteMaintenance(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
