package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo  repository.UserRepository
	tenantRepo repository.TenantRepository
	jwtSecret string
}

func NewAuthHandler(userRepo repository.UserRepository, tenantRepo repository.TenantRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, tenantRepo: tenantRepo, jwtSecret: jwtSecret}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token    string `json:"token"`
	TenantID string `json:"tenant_id"`
	User     struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	} `json:"user"`
}

// Login autentica um usuário e retorna um token JWT
// @Summary Login
// @Description Autentica com e-mail e senha
// @Tags auth (Autenticação)
// @Accept json
// @Produce json
// @Param body body loginRequest true "Credenciais"
// @Success 200 {object} loginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		writeError(w, "email and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		writeError(w, "account is inactive", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	tenant, err := h.tenantRepo.GetByID(user.TenantID)
	if err != nil {
		writeError(w, "tenant not found", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      string(user.Role),
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		writeError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := loginResponse{
		Token:    tokenStr,
		TenantID: tenant.Slug,
	}
	resp.User.ID = user.ID
	resp.User.Email = user.Email
	resp.User.Name = user.Name
	resp.User.Role = string(user.Role)

	writeJSON(w, resp, http.StatusOK)
}
