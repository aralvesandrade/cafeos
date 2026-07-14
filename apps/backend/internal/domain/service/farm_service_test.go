package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func newFarmSvc() (*FarmService, *mock.InMemoryFarmRepo) {
	repo := mock.NewInMemoryFarmRepo()
	producerRepo := mock.NewInMemoryProducerRepo()
	return NewFarmService(repo, producerRepo), repo
}

func TestFarmService_Create(t *testing.T) {
	svc, _ := newFarmSvc()

	farm, err := svc.Create(&entity.Farm{
		OrganizationID: "organization-1",
		Name:           "Fazenda Boa Vista",
		Owner:          "João",
		Location:       "MG",
		TotalAreaHA:    100.0,
		PlantedAreaHA:  80.0,
	}, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if farm.Name != "Fazenda Boa Vista" {
		t.Errorf("expected name 'Fazenda Boa Vista', got %s", farm.Name)
	}
	if farm.OrganizationID != "organization-1" {
		t.Errorf("expected organization-1, got %s", farm.OrganizationID)
	}
	if farm.PlantedAreaHA != 80.0 {
		t.Errorf("expected 80.0, got %f", farm.PlantedAreaHA)
	}
}

func TestFarmService_Create_WithProducer(t *testing.T) {
	svc, _ := newFarmSvc()

	farm, err := svc.Create(&entity.Farm{
		OrganizationID: "organization-1",
		Name:           "Fazenda Boa Vista",
		TotalAreaHA:    100.0,
	}, &entity.Producer{
		Name: "Carlos Eduardo Rosa",
		CPF:  "930.744.338-68",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if farm.Producer == nil {
		t.Fatal("expected producer to be set")
	}
	if farm.Producer.FarmID != farm.ID {
		t.Errorf("expected producer.farm_id %s, got %s", farm.ID, farm.Producer.FarmID)
	}
}

func TestFarmService_Create_Validation(t *testing.T) {
	svc, _ := newFarmSvc()

	tests := []struct {
		name     string
		farmName string
		total    float64
		planted  float64
	}{
		{"empty name", "", 100, 80},
		{"zero total", "Farm", 0, 0},
		{"planted > total", "Farm", 50, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(&entity.Farm{
				OrganizationID: "t1",
				Name:           tt.farmName,
				TotalAreaHA:    tt.total,
				PlantedAreaHA:  tt.planted,
			}, nil)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestFarmService_ListByOrganization(t *testing.T) {
	svc, _ := newFarmSvc()

	svc.Create(&entity.Farm{OrganizationID: "t1", Name: "Farm A", TotalAreaHA: 100, PlantedAreaHA: 80}, nil)
	svc.Create(&entity.Farm{OrganizationID: "t1", Name: "Farm B", TotalAreaHA: 200, PlantedAreaHA: 150}, nil)
	svc.Create(&entity.Farm{OrganizationID: "t2", Name: "Farm C", TotalAreaHA: 50, PlantedAreaHA: 30}, nil)

	farms, err := svc.ListByOrganization("t1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(farms) != 2 {
		t.Errorf("expected 2 farms, got %d", len(farms))
	}
}

func TestFarmService_GetByID_NotFound(t *testing.T) {
	svc, _ := newFarmSvc()

	_, err := svc.GetByID("non-existent")
	if err == nil {
		t.Error("expected error for non-existent farm")
	}
}

func TestFarmService_Update(t *testing.T) {
	svc, _ := newFarmSvc()

	farm, _ := svc.Create(&entity.Farm{OrganizationID: "t1", Name: "Old Name", TotalAreaHA: 100, PlantedAreaHA: 80}, nil)
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
	svc, _ := newFarmSvc()

	farm, _ := svc.Create(&entity.Farm{OrganizationID: "t1", Name: "Farm", TotalAreaHA: 100, PlantedAreaHA: 80}, nil)
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
	svc, _ := newFarmSvc()

	farm, _ := svc.Create(&entity.Farm{OrganizationID: "t1", Name: "Farm", TotalAreaHA: 100, PlantedAreaHA: 80}, nil)
	err := svc.Delete(farm.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = svc.GetByID(farm.ID)
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestFarmService_UpsertProducer(t *testing.T) {
	svc, _ := newFarmSvc()

	farm, _ := svc.Create(&entity.Farm{OrganizationID: "t1", Name: "Farm", TotalAreaHA: 100, PlantedAreaHA: 80}, nil)

	producer, err := svc.UpsertProducer(farm.ID, &entity.Producer{Name: "Carlos"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if producer.FarmID != farm.ID {
		t.Errorf("expected farm_id %s, got %s", farm.ID, producer.FarmID)
	}

	updated, err := svc.UpsertProducer(farm.ID, &entity.Producer{Name: "Carlos Eduardo"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.ID != producer.ID {
		t.Error("expected upsert to reuse the same producer ID")
	}
	if updated.Name != "Carlos Eduardo" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}
}

func TestFarmService_EntityMapping(t *testing.T) {
	farm := &entity.Farm{
		ID:             "1",
		OrganizationID: "t1",
		Name:           "Test Farm",
		Owner:          "Owner",
		Location:       "Location",
		TotalAreaHA:    100.5,
		PlantedAreaHA:  80.3,
	}

	if farm.Name != "Test Farm" {
		t.Error("name mismatch")
	}
	if farm.TotalAreaHA != 100.5 {
		t.Error("total area mismatch")
	}
	if farm.PlantedAreaHA != 80.3 {
		t.Error("planted area mismatch")
	}
}
