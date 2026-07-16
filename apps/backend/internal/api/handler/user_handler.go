package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/api/middleware"
	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/aralvesandrade/cafeos/internal/domain/service"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo    repository.UserRepository
	roleSvc *service.RoleService
	userSvc *service.UserService
}

func NewUserHandler(repo repository.UserRepository, roleSvc *service.RoleService, userSvc *service.UserService) *UserHandler {
	return &UserHandler{repo: repo, roleSvc: roleSvc, userSvc: userSvc}
}

type createUserRequest struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	// RoleID is preferred; Role (a role key, e.g. "operador_campo") is
	// accepted for backward compatibility and resolved via RoleService.
	RoleID string `json:"role_id"`
	Role   string `json:"role"`
}

type updateUserRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	RoleID   *string `json:"role_id"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

// resolveRoleID validates the request's role reference against the global
// roles catalog and returns a concrete RoleID, defaulting to
// "operador_campo" when nothing was supplied.
func (h *UserHandler) resolveRoleID(roleID, roleKey string) (string, error) {
	if roleID != "" {
		role, err := h.roleSvc.GetByID(roleID)
		if err != nil {
			return "", err
		}
		return role.ID, nil
	}
	if roleKey == "" {
		roleKey = "operador_campo"
	}
	role, err := h.roleSvc.FindByKey(roleKey)
	if err != nil {
		return "", err
	}
	return role.ID, nil
}

// resolveRoleIDForCreate is like resolveRoleID, but the organization's very
// first user always becomes organization_admin regardless of what was
// requested — every organization needs one administrator able to register
// farms and other users, and nothing else guarantees that today.
func (h *UserHandler) resolveRoleIDForCreate(organizationID, roleID, roleKey string) (string, error) {
	count, err := h.repo.CountByOrganization(organizationID)
	if err != nil {
		return "", err
	}
	if count == 0 {
		role, err := h.roleSvc.FindByKey(entity.SystemRoleOrganizationAdmin)
		if err != nil {
			return "", err
		}
		return role.ID, nil
	}
	return h.resolveRoleID(roleID, roleKey)
}

// List retorna todos os usuários
// @Summary Listar usuários
// @Description Lista todos os usuários (somente platform_owner)
// @Tags users (Usuários)
// @Produce json
// @Success 200 {array} entity.User
// @Security BearerAuth
// @Router /api/v1/admin/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List()
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}
	writeJSON(w, resp, http.StatusOK)
}

// Create registra um novo usuário
// @Summary Criar usuário
// @Description Cria um novo usuário (somente platform_owner)
// @Tags users (Usuários)
// @Accept json
// @Produce json
// @Param user body createUserRequest true "Dados do usuário"
// @Success 201 {object} entity.User
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.OrganizationID == "" {
		writeError(w, "name, email, password, and organization_id are required", http.StatusBadRequest)
		return
	}

	roleID, err := h.resolveRoleIDForCreate(req.OrganizationID, req.RoleID, req.Role)
	if err != nil {
		writeError(w, "invalid role", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		ID:             uuid.New().String(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   string(hash),
		RoleID:         roleID,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.repo.Create(user); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, user, http.StatusCreated)
}

// Update atualiza um usuário existente
// @Summary Atualizar usuário
// @Description Atualiza dados do usuário (somente platform_owner)
// @Tags users (Usuários)
// @Accept json
// @Produce json
// @Param id path string true "ID do Usuário"
// @Param user body updateUserRequest true "Dados atualizados do usuário"
// @Success 200 {object} entity.User
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.RoleID != nil || req.Role != nil {
		var roleID, roleKey string
		if req.RoleID != nil {
			roleID = *req.RoleID
		}
		if req.Role != nil {
			roleKey = *req.Role
		}
		resolved, err := h.resolveRoleID(roleID, roleKey)
		if err != nil {
			writeError(w, "invalid role", http.StatusBadRequest)
			return
		}
		existing.RoleID = resolved
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, existing, http.StatusOK)
}

// Delete remove um usuário
// @Summary Excluir usuário
// @Description Exclui um usuário por ID (somente platform_owner)
// @Tags users (Usuários)
// @Param id path string true "ID do Usuário"
// @Success 204 "No Content"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.repo.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type userResponse struct {
	ID              string  `json:"id"`
	OrganizationID  string  `json:"organization_id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	RoleID          string  `json:"role_id"`
	Role            string  `json:"role"`
	IsActive        bool    `json:"is_active"`
	Status          string  `json:"status"`
	ManagedByUserID *string `json:"managed_by_user_id"`
	IsPrincipal     bool    `json:"is_principal"`
	PlanID          *string `json:"plan_id,omitempty"`
}

func toUserResponse(u *entity.User) userResponse {
	status := "active"
	if !u.IsActive {
		status = "inactive"
	}
	return userResponse{
		ID:              u.ID,
		OrganizationID:  u.OrganizationID,
		Name:            u.Name,
		Email:           u.Email,
		RoleID:          u.RoleID,
		Role:            u.Role.Key,
		IsActive:        u.IsActive,
		Status:          status,
		ManagedByUserID: u.ManagedByUserID,
		IsPrincipal:     u.ManagedByUserID == nil,
		PlanID:          u.PlanID,
	}
}

// ListForOrg retorna os usuários da própria organização
// @Summary Listar usuários da organização
// @Description Lista os usuários da organização autenticada (requer acesso ao módulo "users")
// @Tags users (Usuários)
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Success 200 {array} entity.User
// @Security BearerAuth
// @Router /api/v1/{organization_id}/users [get]
func (h *UserHandler) ListForOrg(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	users, err := h.repo.ListByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, toUserResponse(u))
	}
	writeJSON(w, resp, http.StatusOK)
}

// CreateForOrg cria um usuário na própria organização
// @Summary Criar usuário na organização
// @Description Cria um usuário na organização autenticada (requer write no módulo "users")
// @Tags users (Usuários)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param user body createUserRequest true "Dados do usuário"
// @Success 201 {object} entity.User
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/users [post]
func (h *UserHandler) CreateForOrg(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	callerUserID, _ := r.Context().Value(middleware.UserIDKey).(string)
	callerRole, _ := r.Context().Value(middleware.RoleKey).(string)

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, "name, email, and password are required", http.StatusBadRequest)
		return
	}

	// A organização pode não ter nenhum usuário ainda (primeiro usuário
	// vira principal automaticamente) — só exige IsPrincipal quando já
	// existe alguém na organização. platform_owner sempre pode criar.
	count, err := h.repo.CountByOrganization(organizationID)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 && callerRole != entity.SystemRolePlatformOwner {
		isPrincipal, err := h.userSvc.IsPrincipal(callerUserID)
		if err != nil {
			writeError(w, "invalid caller", http.StatusBadRequest)
			return
		}
		if !isPrincipal {
			writeError(w, service.ErrNotPrincipal.Error(), http.StatusForbidden)
			return
		}
	}

	roleID, err := h.resolveRoleIDForCreate(organizationID, req.RoleID, req.Role)
	if err != nil {
		writeError(w, "invalid role", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &entity.User{
		ID:             uuid.New().String(),
		OrganizationID: organizationID,
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   string(hash),
		RoleID:         roleID,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	// Sub-usuário: o primeiro usuário da organização (count==0) é sempre
	// principal (ManagedByUserID nulo); os demais ficam vinculados ao
	// principal que os criou — nunca são principais por si mesmos.
	// (platform_owner criando via bypass não vira "gerente" de ninguém —
	// nesse caso o usuário nasce como um novo principal independente.)
	if count > 0 && callerRole != entity.SystemRolePlatformOwner {
		user.ManagedByUserID = &callerUserID
	}

	if err := h.repo.Create(user); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, toUserResponse(user), http.StatusCreated)
}

// UpdateForOrg atualiza um usuário da própria organização
// @Summary Atualizar usuário da organização
// @Description Atualiza um usuário, restrito à organização autenticada (requer write no módulo "users")
// @Tags users (Usuários)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Usuário"
// @Param user body updateUserRequest true "Dados atualizados do usuário"
// @Success 200 {object} entity.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/users/{id} [put]
func (h *UserHandler) UpdateForOrg(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(id)
	if err != nil || existing.OrganizationID != organizationID {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.RoleID != nil || req.Role != nil {
		var roleID, roleKey string
		if req.RoleID != nil {
			roleID = *req.RoleID
		}
		if req.Role != nil {
			roleKey = *req.Role
		}
		resolved, err := h.resolveRoleID(roleID, roleKey)
		if err != nil {
			writeError(w, "invalid role", http.StatusBadRequest)
			return
		}
		existing.RoleID = resolved
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(existing); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, toUserResponse(existing), http.StatusOK)
}

// DeleteForOrg remove um usuário da própria organização
// @Summary Excluir usuário da organização
// @Description Exclui um usuário, restrito à organização autenticada (requer write no módulo "users")
// @Tags users (Usuários)
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Usuário"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/users/{id} [delete]
func (h *UserHandler) DeleteForOrg(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(id)
	if err != nil || existing.OrganizationID != organizationID {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type assignPlanRequest struct {
	PlanID string `json:"plan_id"`
}

// AssignPlanForOrg vincula um Plan ao usuário principal (dono do grupo de
// fazendas), responsável pelo pagamento — sub-usuário nunca tem plano
// próprio.
// @Summary Atribuir plano ao usuário principal
// @Description Vincula um Plan ao usuário — só permitido para o próprio usuário principal ou platform_owner
// @Tags users (Usuários)
// @Accept json
// @Produce json
// @Param organization_id path string true "ID da Organização"
// @Param id path string true "ID do Usuário"
// @Param plan body assignPlanRequest true "Plano a atribuir"
// @Success 200 {object} entity.User
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/{organization_id}/users/{id}/plan [patch]
func (h *UserHandler) AssignPlanForOrg(w http.ResponseWriter, r *http.Request) {
	organizationID, _ := r.Context().Value(middleware.OrganizationIDKey).(string)
	callerUserID, _ := r.Context().Value(middleware.UserIDKey).(string)
	callerRole, _ := r.Context().Value(middleware.RoleKey).(string)
	id := r.PathValue("id")

	existing, err := h.repo.GetByID(id)
	if err != nil || existing.OrganizationID != organizationID {
		writeError(w, "user not found", http.StatusNotFound)
		return
	}

	if callerRole != entity.SystemRolePlatformOwner && callerUserID != id {
		writeError(w, "only the principal user themselves (or platform_owner) can assign their plan", http.StatusForbidden)
		return
	}

	var req assignPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PlanID == "" {
		writeError(w, "plan_id is required", http.StatusBadRequest)
		return
	}

	if err := h.userSvc.AssignPlan(id, req.PlanID); err != nil {
		if err == service.ErrNotPrincipal {
			writeError(w, err.Error(), http.StatusForbidden)
			return
		}
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.repo.GetByID(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, toUserResponse(updated), http.StatusOK)
}
