package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func TestRoleService_SeedDefaultsIfMissing_IsIdempotent(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if err := svc.SeedDefaultsIfMissing(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if err := svc.SeedDefaultsIfMissing(); err != nil {
		t.Fatalf("expected no error on second call, got %v", err)
	}

	roles, _ := roleRepo.List()
	want := len(entity.SystemRoleKeys) + len(entity.DefaultRoleKeys)
	if len(roles) != want {
		t.Errorf("expected %d roles, got %d", want, len(roles))
	}
}

func TestRoleService_Create_RejectsDuplicateKey(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if _, err := svc.Create("colhedor_chefe", "Colhedor Chefe"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, err := svc.Create("colhedor_chefe", "Outro Nome"); err == nil {
		t.Error("expected error for duplicate key")
	}
}

func TestRoleService_Update_RejectsSystemRole(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if err := svc.SeedDefaultsIfMissing(); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	roles, _ := roleRepo.List()
	var platformOwnerID string
	for _, r := range roles {
		if r.Key == entity.SystemRolePlatformOwner {
			platformOwnerID = r.ID
		}
	}

	if _, err := svc.Update(platformOwnerID, "New Name"); err != ErrRoleIsSystem {
		t.Errorf("expected ErrRoleIsSystem, got %v", err)
	}
}

func TestRoleService_Delete_RejectsSystemRole(t *testing.T) {
	roleRepo := mock.NewInMemoryRoleRepo()
	svc := NewRoleService(roleRepo, nil)

	if err := svc.SeedDefaultsIfMissing(); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	roles, _ := roleRepo.List()
	var platformOwnerID string
	for _, r := range roles {
		if r.Key == entity.SystemRolePlatformOwner {
			platformOwnerID = r.ID
		}
	}

	if err := svc.Delete(platformOwnerID); err != ErrRoleIsSystem {
		t.Errorf("expected ErrRoleIsSystem, got %v", err)
	}
}
