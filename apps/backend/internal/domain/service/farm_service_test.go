package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func TestFarmService_Create(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	farm, err := svc.Create("tenant-1", "Fazenda Boa Vista", "João", "MG", 100.0, 80.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if farm.Name != "Fazenda Boa Vista" {
		t.Errorf("expected name 'Fazenda Boa Vista', got %s", farm.Name)
	}
	if farm.TenantID != "tenant-1" {
		t.Errorf("expected tenant-1, got %s", farm.TenantID)
	}
	if farm.PlantedAreaHA != 80.0 {
		t.Errorf("expected 80.0, got %f", farm.PlantedAreaHA)
	}
}

func TestFarmService_Create_Validation(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	tests := []struct {
		name    string
		farmName string
		total   float64
		planted float64
	}{
		{"empty name", "", 100, 80},
		{"zero total", "Farm", 0, 0},
		{"planted > total", "Farm", 50, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create("t1", tt.farmName, "", "", tt.total, tt.planted)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestFarmService_ListByTenant(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	svc.Create("t1", "Farm A", "", "", 100, 80)
	svc.Create("t1", "Farm B", "", "", 200, 150)
	svc.Create("t2", "Farm C", "", "", 50, 30)

	farms, err := svc.ListByTenant("t1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(farms) != 2 {
		t.Errorf("expected 2 farms, got %d", len(farms))
	}
}

func TestFarmService_GetByID_NotFound(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	_, err := svc.GetByID("non-existent")
	if err == nil {
		t.Error("expected error for non-existent farm")
	}
}

func TestFarmService_Update(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	farm, _ := svc.Create("t1", "Old Name", "", "", 100, 80)
	farm.Name = "New Name"
	err := svc.Update(farm)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := svc.GetByID(farm.ID)
	if updated.Name != "New Name" {
		t.Errorf("expected 'New Name', got %s", updated.Name)
	}
}

func TestFarmService_Update_Validation(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	farm, _ := svc.Create("t1", "Farm", "", "", 100, 80)
	farm.Name = ""
	err := svc.Update(farm)
	if err == nil {
		t.Error("expected error for empty name")
	}

	farm.Name = "Farm"
	farm.PlantedAreaHA = 200
	farm.TotalAreaHA = 100
	err = svc.Update(farm)
	if err == nil {
		t.Error("expected error for planted > total")
	}
}

func TestFarmService_Delete(t *testing.T) {
	repo := mock.NewInMemoryFarmRepo()
	svc := NewFarmService(repo)

	farm, _ := svc.Create("t1", "Farm", "", "", 100, 80)
	err := svc.Delete(farm.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = svc.GetByID(farm.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestFarmService_EntityMapping(t *testing.T) {
	farm := &entity.Farm{
		ID:           "1",
		TenantID:     "t1",
		Name:         "Test Farm",
		Owner:        "Owner",
		Location:     "Location",
		TotalAreaHA:  100.5,
		PlantedAreaHA: 80.3,
	}

	if farm.Name != "Test Farm" { t.Error("name mismatch") }
	if farm.TotalAreaHA != 100.5 { t.Error("total area mismatch") }
	if farm.PlantedAreaHA != 80.3 { t.Error("planted area mismatch") }
}
