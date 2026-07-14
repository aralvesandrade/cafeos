package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type createUserRequest struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Role           string `json:"role"`
}

type updateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
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
	// Strip password hashes from response
	type userResponse struct {
		ID             string `json:"id"`
		OrganizationID string `json:"organization_id"`
		Name           string `json:"name"`
		Email          string `json:"email"`
		Role           string `json:"role"`
		IsActive       bool   `json:"is_active"`
		Status         string `json:"status"`
	}
	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		status := "active"
		if !u.IsActive {
			status = "inactive"
		}
		resp = append(resp, userResponse{
			ID:             u.ID,
			OrganizationID: u.OrganizationID,
			Name:           u.Name,
			Email:          u.Email,
			Role:           string(u.Role),
			IsActive:       u.IsActive,
			Status:         status,
		})
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

	if req.Name == "" || req.Email == "" || req.Password == "" {
		writeError(w, "name, email, and password are required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	role := entity.UserRole(req.Role)
	if role == "" {
		role = entity.RoleOperador
	}

	user := &entity.User{
		ID:             uuid.New().String(),
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Email:          req.Email,
		PasswordHash:   string(hash),
		Role:           role,
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
	if req.Role != nil {
		existing.Role = entity.UserRole(*req.Role)
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
