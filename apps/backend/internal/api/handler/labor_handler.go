package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type LaborHandler struct {
	svc *service.LaborService
}

func NewLaborHandler(svc *service.LaborService) *LaborHandler {
	return &LaborHandler{svc: svc}
}

// Teams
func (h *LaborHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req struct {
		Name        string `json:"name"`
		Leader      string `json:"leader"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	team, err := h.svc.CreateTeam(tenantID, req.Name, req.Leader, req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, team, http.StatusCreated)
}

func (h *LaborHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	teams, err := h.svc.ListTeams(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, teams, http.StatusOK)
}

func (h *LaborHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.svc.GetTeamByID(id)
	if err != nil {
		writeError(w, "team not found", http.StatusNotFound)
		return
	}
	var req struct {
		Name        *string `json:"name"`
		Leader      *string `json:"leader"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name != nil { existing.Name = *req.Name }
	if req.Leader != nil { existing.Leader = *req.Leader }
	if req.Description != nil { existing.Description = *req.Description }
	if err := h.svc.UpdateTeam(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

func (h *LaborHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteTeam(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Workers
func (h *LaborHandler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req struct {
		TeamID     string  `json:"team_id"`
		Name       string  `json:"name"`
		Role       string  `json:"role"`
		Phone      string  `json:"phone"`
		HourlyRate float64 `json:"hourly_rate"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	wkr, err := h.svc.CreateWorker(tenantID, req.TeamID, req.Name, req.Role, req.Phone, req.HourlyRate)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, wkr, http.StatusCreated)
}

func (h *LaborHandler) ListWorkers(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	workers, err := h.svc.ListWorkers(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, workers, http.StatusOK)
}

func (h *LaborHandler) UpdateWorker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing, err := h.svc.GetWorkerByID(id)
	if err != nil {
		writeError(w, "worker not found", http.StatusNotFound)
		return
	}
	var req struct {
		Name       *string  `json:"name"`
		Role       *string  `json:"role"`
		Phone      *string  `json:"phone"`
		HourlyRate *float64 `json:"hourly_rate"`
		IsActive   *bool    `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name != nil { existing.Name = *req.Name }
	if req.Role != nil { existing.Role = *req.Role }
	if req.Phone != nil { existing.Phone = *req.Phone }
	if req.HourlyRate != nil { existing.HourlyRate = *req.HourlyRate }
	if req.IsActive != nil { existing.IsActive = *req.IsActive }
	if err := h.svc.UpdateWorker(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

func (h *LaborHandler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteWorker(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// WorkShifts
func (h *LaborHandler) CreateWorkShift(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	var req struct {
		WorkerID    string  `json:"worker_id"`
		OperationID string  `json:"operation_id"`
		Hours       float64 `json:"hours"`
		Cost        float64 `json:"cost"`
		Date        string  `json:"date"`
		Notes       string  `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", req.Date)
	ws, err := h.svc.CreateWorkShift(tenantID, req.WorkerID, req.OperationID, req.Notes, req.Hours, req.Cost, date)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, ws, http.StatusCreated)
}

func (h *LaborHandler) ListWorkShifts(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := r.Context().Value(middleware.TenantIDKey).(string)
	shifts, err := h.svc.ListWorkShifts(tenantID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, shifts, http.StatusOK)
}

func (h *LaborHandler) DeleteWorkShift(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteWorkShift(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
