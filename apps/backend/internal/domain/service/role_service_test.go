package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func TestRoleService_SeedSystemRolesIfMissing_IsIdempotent(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if err := svc.SeedSystemRolesIfMissing(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := svc.SeedSystemRolesIfMissing(); err != nil {
		t.Fatalf("expected no error on second call, got %v", err)
	}

	roles, _ := roleRepo.ListForOrganization("")
	if len(roles) != len(entity.SystemRoleKeys) {
		t.Errorf("expected %d system roles, got %d", len(entity.SystemRoleKeys), len(roles))
	}
}

func TestRoleService_SeedDefaultsIfMissing_DoesNotDuplicate(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if err := svc.SeedDefaultsIfMissing("org-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := svc.SeedDefaultsIfMissing("org-1"); err != nil {
		t.Fatalf("expected no error on second call, got %v", err)
	}

	count, _ := roleRepo.CountByOrganization("org-1")
	if count != int64(len(entity.DefaultOrgRoleKeys)) {
		t.Errorf("expected %d default roles, got %d", len(entity.DefaultOrgRoleKeys), count)
	}
}

func TestRoleService_Create_RejectsDuplicateKey(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if _, err := svc.Create("org-1", "colhedor_chefe", "Colhedor Chefe"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := svc.Create("org-1", "colhedor_chefe", "Outro Nome"); err == nil {
		t.Error("expected error for duplicate key within the same organization")
	}
}

func TestRoleService_Update_RejectsSystemRole(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if err := svc.SeedSystemRolesIfMissing(); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	roles, _ := roleRepo.ListForOrganization("org-1")
	var platformOwnerID string
	for _, r := range roles {
		if r.Key == entity.SystemRolePlatformOwner {
			platformOwnerID = r.ID
		}
	}

	if _, err := svc.Update("org-1", platformOwnerID, "New Name"); err != ErrRoleIsSystem {
		t.Errorf("expected ErrRoleIsSystem, got %v", err)
	}
}

func TestRoleService_Delete_RejectsRoleFromAnotherOrg(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	role, err := svc.Create("org-1", "colhedor_chefe", "Colhedor Chefe")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := svc.Delete("org-2", role.ID); err != ErrForbiddenOrg {
		t.Errorf("expected ErrForbiddenOrg, got %v", err)
	}
}
