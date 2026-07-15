package service

import (
	"testing"
	"time"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
	"github.com/google/uuid"
)

// seedTestRoles populates the given role repo with the ten default roles
// (two system + eight org-scoped) for "org-1", returning role IDs by key.
func seedTestRoles(t *testing.T, repo *mock.InMemoryRoleRepo) map[string]string {
	t.Helper()
	ids := make(map[string]string)
	for _, key := range entity.SystemRoleKeys {
		role := &entity.Role{ID: uuid.New().String(), Key: key, Name: key, IsSystem: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if err := repo.Create(role); err != nil {
			t.Fatalf("seed role %s: %v", key, err)
		}
		ids[key] = role.ID
	}
	org := "org-1"
	for _, key := range entity.DefaultOrgRoleKeys {
		role := &entity.Role{ID: uuid.New().String(), OrganizationID: &org, Key: key, Name: key, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		if err := repo.Create(role); err != nil {
			t.Fatalf("seed role %s: %v", key, err)
		}
		ids[key] = role.ID
	}
	return ids
}

func TestPermissionService_SeedDefaults(t *testing.T) {
	permRepo := mock.NewInMemoryPermissionRepo()
	roleRepo := mock.NewInMemoryRoleRepo()
	seedTestRoles(t, roleRepo)
	svc := NewPermissionService(permRepo, roleRepo)

	if err := svc.SeedDefaults("org-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	count, _ := permRepo.CountByOrganization("org-1")
	want := int64(len(entity.AllModules) * 10)
	if count != want {
		t.Errorf("expected %d seeded rows, got %d", want, count)
	}
}

func TestPermissionService_GetAccess_FallsBackToDefaultWhenUnseeded(t *testing.T) {
	permRepo := mock.NewInMemoryPermissionRepo()
	roleRepo := mock.NewInMemoryRoleRepo()
	seedTestRoles(t, roleRepo)
	svc := NewPermissionService(permRepo, roleRepo)

	access, err := svc.GetAccess("org-unseeded", "operador_campo", entity.ModuleFinancial)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access != entity.AccessNone {
		t.Errorf("expected operador_campo to have no access to financial by default, got %s", access)
	}

	access, err = svc.GetAccess("org-unseeded", "financeiro", entity.ModuleFinancial)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access != entity.AccessWrite {
		t.Errorf("expected financeiro to have write access to financial by default, got %s", access)
	}
}

func TestPermissionService_Update_OverridesAndInvalidatesCache(t *testing.T) {
	permRepo := mock.NewInMemoryPermissionRepo()
	roleRepo := mock.NewInMemoryRoleRepo()
	ids := seedTestRoles(t, roleRepo)
	svc := NewPermissionService(permRepo, roleRepo)

	if err := svc.SeedDefaults("org-1"); err != nil {
		t.Fatalf("seed failed: %v", err)
	}

	// Warm the cache with the seeded value.
	if access, _ := svc.GetAccess("org-1", "operador_campo", entity.ModuleFinancial); access != entity.AccessNone {
		t.Fatalf("expected seeded default 'none', got %s", access)
	}

	if err := svc.Update("org-1", ids["operador_campo"], entity.ModuleFinancial, entity.AccessRead); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	access, err := svc.GetAccess("org-1", "operador_campo", entity.ModuleFinancial)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if access != entity.AccessRead {
		t.Errorf("expected updated access 'read', got %s", access)
	}
}

func TestPermissionService_Update_Validation(t *testing.T) {
	permRepo := mock.NewInMemoryPermissionRepo()
	roleRepo := mock.NewInMemoryRoleRepo()
	ids := seedTestRoles(t, roleRepo)
	svc := NewPermissionService(permRepo, roleRepo)

	if err := svc.Update("org-1", "role-invalido", entity.ModuleFinancial, entity.AccessRead); err == nil {
		t.Error("expected error for invalid role")
	}
	if err := svc.Update("org-1", ids["operador_campo"], "modulo_invalido", entity.AccessRead); err == nil {
		t.Error("expected error for invalid module")
	}
	if err := svc.Update("org-1", ids["operador_campo"], entity.ModuleFinancial, "acesso_invalido"); err == nil {
		t.Error("expected error for invalid access level")
	}
}

func TestPermissionService_SeedDefaultsIfMissing_DoesNotOverwrite(t *testing.T) {
	permRepo := mock.NewInMemoryPermissionRepo()
	roleRepo := mock.NewInMemoryRoleRepo()
	ids := seedTestRoles(t, roleRepo)
	svc := NewPermissionService(permRepo, roleRepo)

	if err := svc.SeedDefaults("org-1"); err != nil {
		t.Fatalf("seed failed: %v", err)
	}
	if err := svc.Update("org-1", ids["operador_campo"], entity.ModuleFinancial, entity.AccessRead); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	if err := svc.SeedDefaultsIfMissing("org-1"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	access, _ := svc.GetAccess("org-1", "operador_campo", entity.ModuleFinancial)
	if access != entity.AccessRead {
		t.Errorf("expected manual override 'read' to survive, got %s", access)
	}
}
