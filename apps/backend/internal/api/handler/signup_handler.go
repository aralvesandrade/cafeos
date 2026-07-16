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
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	Phone       string  `json:"phone"`
	FarmName    string  `json:"farm_name"`
	State       string  `json:"state"`
	City        string  `json:"city"`
	MainCrop    string  `json:"main_crop"`
	TotalAreaHA float64 `json:"total_area_ha"`
}

type registerResponse struct {
	UserID string `json:"user_id"`
	FarmID string `json:"farm_id"`
}

// Register realiza o cadastro público de um novo proprietário + sua fazenda
// @Summary Cadastro público de proprietário
// @Description Cria um usuário principal (proprietario) e sua fazenda na organização padrão da plataforma — não exige autenticação
// @Tags auth (Autenticação)
// @Accept json
// @Produce json
// @Param body body registerRequest true "Dados do proprietário e da fazenda"
// @Success 201 {object} registerResponse
// @Failure 400 {object} map[string]string
// @Router /auth/register [post]
func (h *SignupHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, farm, err := h.svc.RegisterPrincipal(service.RegisterPrincipalInput{
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		Phone:       req.Phone,
		FarmName:    req.FarmName,
		State:       req.State,
		City:        req.City,
		MainCrop:    req.MainCrop,
		TotalAreaHA: req.TotalAreaHA,
	})
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, registerResponse{UserID: user.ID, FarmID: farm.ID}, http.StatusCreated)
}
