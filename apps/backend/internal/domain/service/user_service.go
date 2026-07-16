package service

import (
	"errors"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/domain/repository"
)

var ErrNotPrincipal = errors.New("only a principal user (farm group owner) can perform this action")

// UserService concentra as regras da hierarquia de 2 níveis entre usuários
// principais (donos de um grupo de fazendas, responsáveis pelo próprio
// Plan) e os sub-usuários que eles criam.
type UserService struct {
	repo     repository.UserRepository
	planRepo repository.PlanRepository
}

func NewUserService(repo repository.UserRepository, planRepo repository.PlanRepository) *UserService {
	return &UserService{repo: repo, planRepo: planRepo}
}

// Get retorna o usuário pelo ID.
func (s *UserService) Get(userID string) (*entity.User, error) {
	return s.repo.GetByID(userID)
}

// IsPrincipal retorna true quando o usuário é raiz da própria hierarquia
// (ManagedByUserID nulo) — dono de fazendas, responsável pelo próprio
// plano, e o único autorizado a criar sub-usuários e fazendas.
func (s *UserService) IsPrincipal(userID string) (bool, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return false, err
	}
	return user.ManagedByUserID == nil, nil
}

// EnsurePrincipal valida que managerID pode criar sub-usuários/fazendas
// (é principal) e que targetUserID, quando informado, já é um sub-usuário
// dele — usado ao vincular uma fazenda a um usuário específico.
func (s *UserService) EnsureManages(managerID, targetUserID string) error {
	target, err := s.repo.GetByID(targetUserID)
	if err != nil {
		return err
	}
	if target.ManagedByUserID == nil || *target.ManagedByUserID != managerID {
		return errors.New("user is not a sub-user of the requesting principal")
	}
	return nil
}

// AssignPlan vincula um Plan ao usuário — só permitido quando o usuário é
// principal (sub-usuário nunca tem plano/pagamento próprio).
func (s *UserService) AssignPlan(userID, planID string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.ManagedByUserID != nil {
		return ErrNotPrincipal
	}
	if _, err := s.planRepo.GetByID(planID); err != nil {
		return errors.New("plan not found")
	}
	user.PlanID = &planID
	return s.repo.Update(user)
}
