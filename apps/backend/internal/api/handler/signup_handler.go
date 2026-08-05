package handler

import (
	"encoding/json"
	"net/http"

	"github.com/aralvesandrade/cafeos/internal/domain/service"
)

type SignupHandler struct {
	svc *service.SignupService
}

func NewSignupHandler(svc *service.SignupService) *SignupHandler {
	return &SignupHandler{svc: svc}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PlanSlug string `json:"plan_slug"`
}

type registerResponse struct {
	UserID string `json:"user_id"`
}

// Register realiza o cadastro público de um novo proprietário
// @Summary Cadastro público de proprietário
// @Description Cria um usuário principal (proprietario) na organização padrão da plataforma — não exige autenticação. A fazenda é cadastrada depois, já autenticado.
// @Tags auth (Autenticação)
// @Accept json
// @Produce json
// @Param body body registerRequest true "Dados do proprietário"
// @Success 201 {object} registerResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *SignupHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.svc.RegisterPrincipal(service.RegisterPrincipalInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		PlanSlug: req.PlanSlug,
	})
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, registerResponse{UserID: user.ID}, http.StatusCreated)
}
