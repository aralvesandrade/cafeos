package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func TestPermissionService_SeedDefaults(t *testing.T) {
	repo := mock.NewInMemoryPermissionRepo()
	svc := NewPermissionService(repo)

	if err := svc.SeedDefaults("org-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	count, _ := repo.CountByOrganization("org-1")
	want := int64(len(entity.AllModules) * len(allRoles))
	if count != want {
		t.Errorf("expected %d seeded rows, got %d", want, count)
	}
}

func TestPermissionService_GetAccess_FallsBackToDefaultWhenUnseeded(t *testing.T) {
	repo := mock.NewInMemoryPermissionRepo()
	svc := NewPermissionService(repo)

	access, err := svc.GetAccess("org-unseeded", entity.RoleOperador, entity.ModuleFinancial)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access != entity.AccessNone {
		t.Errorf("expected operador_campo to have no access to financial by default, got %s", access)
	}

	access, err = svc.GetAccess("org-unseeded", entity.RoleFinanceiro, entity.ModuleFinancial)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access != entity.AccessWrite {
		t.Errorf("expected financeiro to have write access to financial by default, got %s", access)
	}
}

func TestPermissionService_Update_OverridesAndInvalidatesCache(t *testing.T) {
	repo := mock.NewInMemoryPermissionRepo()
	svc := NewPermissionService(repo)

	if err := svc.SeedDefaults("org-1"); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	// Warm the cache with the seeded value.
	if access, _ := svc.GetAccess("org-1", entity.RoleOperador, entity.ModuleFinancial); access != entity.AccessNone {
		t.Fatalf("expected seeded default 'none', got %s", access)
	}

	if err := svc.Update("org-1", entity.RoleOperador, entity.ModuleFinancial, entity.AccessRead); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	access, err := svc.GetAccess("org-1", entity.RoleOperador, entity.ModuleFinancial)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access != entity.AccessRead {
		t.Errorf("expected updated access 'read', got %s", access)
	}
}

func TestPermissionService_Update_Validation(t *testing.T) {
	repo := mock.NewInMemoryPermissionRepo()
	svc := NewPermissionService(repo)

	if err := svc.Update("org-1", "papel_invalido", entity.ModuleFinancial, entity.AccessRead); err == nil {
		t.Error("expected error for invalid role")
	}
	if err := svc.Update("org-1", entity.RoleOperador, "modulo_invalido", entity.AccessRead); err == nil {
		t.Error("expected error for invalid module")
	}
	if err := svc.Update("org-1", entity.RoleOperador, entity.ModuleFinancial, "acesso_invalido"); err == nil {
		t.Error("expected error for invalid access level")
	}
}

func TestPermissionService_SeedDefaultsIfMissing_DoesNotOverwrite(t *testing.T) {
	repo := mock.NewInMemoryPermissionRepo()
	svc := NewPermissionService(repo)

	if err := svc.SeedDefaults("org-1"); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	if err := svc.Update("org-1", entity.RoleOperador, entity.ModuleFinancial, entity.AccessRead); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if err := svc.SeedDefaultsIfMissing("org-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	access, _ := svc.GetAccess("org-1", entity.RoleOperador, entity.ModuleFinancial)
	if access != entity.AccessRead {
		t.Errorf("expected manual override 'read' to survive, got %s", access)
	}
}
