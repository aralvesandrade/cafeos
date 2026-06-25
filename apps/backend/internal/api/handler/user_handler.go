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
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type updateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *string `json:"role"`
}

// List returns all users
// @Summary List users
// @Description List all users (platform_owner only)
// @Tags users
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
		ID       string `json:"id"`
		TenantID string `json:"tenant_id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		IsActive bool   `json:"is_active"`
		Status   string `json:"status"`
	}
	resp := make([]userResponse, 0, len(users))
	for _, u := range users {
		status := "active"
		if !u.IsActive {
			status = "inactive"
		}
		resp = append(resp, userResponse{
			ID:       u.ID,
			TenantID: u.TenantID,
			Name:     u.Name,
			Email:    u.Email,
			Role:     string(u.Role),
			IsActive: u.IsActive,
			Status:   status,
		})
	}
	writeJSON(w, resp, http.StatusOK)
}

// Create registers a new user
// @Summary Create a user
// @Description Create a new user (platform_owner only)
// @Tags users
// @Accept json
// @Produce json
// @Param user body createUserRequest true "User data"
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
		ID:           uuid.New().String(),
		TenantID:     req.TenantID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.repo.Create(user); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, user, http.StatusCreated)
}

// Update updates an existing user
// @Summary Update a user
// @Description Update user data (platform_owner only)
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body updateUserRequest true "Updated user data"
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

// Delete removes a user
// @Summary Delete a user
// @Description Delete a user by ID (platform_owner only)
// @Tags users
// @Param id path string true "User ID"
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
