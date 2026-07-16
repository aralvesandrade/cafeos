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

type createTeamRequest struct {
	Name        string `json:"name"`
	Leader      string `json:"leader"`
	Description string `json:"description"`
}

// CreateTeam registra uma nova equipe de trabalho para a organização autenticada
// @Summary Criar equipe de trabalho
// @Description Registra uma nova equipe de trabalho na organização
// @Tags labor (Mão de Obra)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param team body createTeamRequest true "Dados da equipe"
// @Success 201 {object} entity.Team
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/teams [post]
func (h *LaborHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	var req createTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	team, err := h.svc.CreateTeam(organizationID, req.Name, req.Leader, req.Description)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, team, http.StatusCreated)
}

// ListTeams retorna todas as equipes de trabalho da organização autenticada
// @Summary Listar equipes de trabalho
// @Description Lista todas as equipes de trabalho pertencentes à organização
// @Tags labor (Mão de Obra)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.Team
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/teams [get]
func (h *LaborHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	teams, err := h.svc.ListTeams(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, teams, http.StatusOK)
}

// GetTeamByID retorna uma equipe de trabalho pelo seu ID
// @Summary Obter equipe de trabalho por ID
// @Description Retorna uma única equipe de trabalho
// @Tags labor (Mão de Obra)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Equipe"
// @Success 200 {object} entity.Team
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/teams/{id} [get]
func (h *LaborHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	team, err := h.svc.GetTeamByID(id)
	if err != nil {
		writeError(w, "team not found", http.StatusNotFound)
		return
	}
	writeJSON(w, team, http.StatusOK)
}

// UpdateTeam atualiza uma equipe de trabalho existente
// @Summary Atualizar equipe de trabalho
// @Description Atualiza dados da equipe de trabalho por ID (atualização parcial - somente os campos informados são alterados)
// @Tags labor (Mão de Obra)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Equipe"
// @Param team body entity.Team true "Dados atualizados da equipe"
// @Success 200 {object} entity.Team
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/teams/{id} [put]
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
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Leader != nil {
		existing.Leader = *req.Leader
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if err := h.svc.UpdateTeam(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// DeleteTeam remove uma equipe de trabalho pelo seu ID
// @Summary Excluir equipe de trabalho
// @Description Exclui uma equipe de trabalho por ID
// @Tags labor (Mão de Obra)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID da Equipe"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/teams/{id} [delete]
func (h *LaborHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteTeam(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Workers

type createWorkerRequest struct {
	TeamID     string  `json:"team_id"`
	Name       string  `json:"name"`
	Role       string  `json:"role"`
	Phone      string  `json:"phone"`
	HourlyRate float64 `json:"hourly_rate"`
}

// CreateWorker registra um novo trabalhador para a organização autenticada
// @Summary Criar trabalhador
// @Description Registra um novo trabalhador na organização
// @Tags labor (Mão de Obra)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param worker body createWorkerRequest true "Dados do trabalhador"
// @Success 201 {object} entity.Worker
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/workers [post]
func (h *LaborHandler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	var req createWorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	wkr, err := h.svc.CreateWorker(organizationID, req.TeamID, req.Name, req.Role, req.Phone, req.HourlyRate)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, wkr, http.StatusCreated)
}

// ListWorkers retorna todos os trabalhadores da organização autenticada
// @Summary Listar trabalhadores
// @Description Lista todos os trabalhadores pertencentes à organização
// @Tags labor (Mão de Obra)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.Worker
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/workers [get]
func (h *LaborHandler) ListWorkers(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	workers, err := h.svc.ListWorkers(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, workers, http.StatusOK)
}

// GetWorkerByID retorna um trabalhador pelo seu ID
// @Summary Obter trabalhador por ID
// @Description Retorna um único trabalhador
// @Tags labor (Mão de Obra)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Trabalhador"
// @Success 200 {object} entity.Worker
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/workers/{id} [get]
func (h *LaborHandler) GetWorkerByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	worker, err := h.svc.GetWorkerByID(id)
	if err != nil {
		writeError(w, "worker not found", http.StatusNotFound)
		return
	}
	writeJSON(w, worker, http.StatusOK)
}

// UpdateWorker atualiza um trabalhador existente
// @Summary Atualizar trabalhador
// @Description Atualiza dados do trabalhador por ID (atualização parcial - somente os campos informados são alterados)
// @Tags labor (Mão de Obra)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Trabalhador"
// @Param worker body entity.Worker true "Dados atualizados do trabalhador"
// @Success 200 {object} entity.Worker
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/workers/{id} [put]
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
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Role != nil {
		existing.Role = *req.Role
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.HourlyRate != nil {
		existing.HourlyRate = *req.HourlyRate
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	if err := h.svc.UpdateWorker(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, existing, http.StatusOK)
}

// DeleteWorker remove um trabalhador pelo seu ID
// @Summary Excluir trabalhador
// @Description Exclui um trabalhador por ID
// @Tags labor (Mão de Obra)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Trabalhador"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/workers/{id} [delete]
func (h *LaborHandler) DeleteWorker(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteWorker(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// WorkShifts

type createWorkShiftRequest struct {
	WorkerID    string  `json:"worker_id"`
	OperationID string  `json:"operation_id"`
	Hours       float64 `json:"hours"`
	Cost        float64 `json:"cost"`
	Date        string  `json:"date"`
	Notes       string  `json:"notes"`
}

// CreateWorkShift registra um novo turno de trabalho para um trabalhador
// @Summary Criar turno de trabalho
// @Description Registra um novo turno de trabalho de um trabalhador em uma operação
// @Tags labor (Mão de Obra)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param shift body createWorkShiftRequest true "Dados do turno de trabalho"
// @Success 201 {object} entity.WorkShift
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/shifts [post]
func (h *LaborHandler) CreateWorkShift(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	var req createWorkShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	date, _ := time.Parse("2006-01-02", req.Date)
	ws, err := h.svc.CreateWorkShift(organizationID, req.WorkerID, req.OperationID, req.Notes, req.Hours, req.Cost, date)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, ws, http.StatusCreated)
}

// ListWorkShifts retorna todos os turnos de trabalho da organização autenticada
// @Summary Listar turnos de trabalho
// @Description Lista todos os turnos de trabalho pertencentes à organização
// @Tags labor (Mão de Obra)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.WorkShift
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/shifts [get]
func (h *LaborHandler) ListWorkShifts(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	shifts, err := h.svc.ListWorkShifts(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, shifts, http.StatusOK)
}

// GetWorkShiftByID retorna um turno de trabalho pelo seu ID
// @Summary Obter turno de trabalho por ID
// @Description Retorna um único turno de trabalho
// @Tags labor (Mão de Obra)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Turno de Trabalho"
// @Success 200 {object} entity.WorkShift
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/shifts/{id} [get]
func (h *LaborHandler) GetWorkShiftByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	shift, err := h.svc.GetWorkShiftByID(id)
	if err != nil {
		writeError(w, "work shift not found", http.StatusNotFound)
		return
	}
	writeJSON(w, shift, http.StatusOK)
}

// DeleteWorkShift remove um turno de trabalho pelo seu ID
// @Summary Excluir turno de trabalho
// @Description Exclui um turno de trabalho por ID
// @Tags labor (Mão de Obra)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Turno de Trabalho"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/labor/shifts/{id} [delete]
func (h *LaborHandler) DeleteWorkShift(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteWorkShift(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
