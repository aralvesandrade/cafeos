package service

import (
	"testing"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func TestPlotService_Create(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	plot, err := svc.Create(&entity.Plot{
		OrganizationID: "t1",
		FarmID:         "farm-1",
		Name:           "Talhão A",
		Cultivar:       "Catuaí",
		SoilType:       "Argiloso",
		AreaHA:         10.5,
		PlantingYear:   2020,
		Altitude:       950,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if plot.Name != "Talhão A" {
		t.Errorf("expected 'Talhão A', got %s", plot.Name)
	}
	if plot.Cultivar != "Catuaí" {
		t.Errorf("expected 'Catuaí', got %s", plot.Cultivar)
	}
	if plot.AreaHA != 10.5 {
		t.Errorf("expected 10.5, got %f", plot.AreaHA)
	}
	if plot.Stage != entity.PlotStageFormacao {
		t.Errorf("expected default stage 'formacao', got %s", plot.Stage)
	}
}

func TestPlotService_Create_Validation(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	_, err := svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "f1", Name: "", Cultivar: "Catuaí", SoilType: "Argiloso", AreaHA: 10, PlantingYear: 2020, Altitude: 950})
	if err == nil {
		t.Error("expected error for empty name")
	}

	_, err = svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "f1", Name: "Plot", Cultivar: "Catuaí", SoilType: "Argiloso", AreaHA: 0, PlantingYear: 2020, Altitude: 950})
	if err == nil {
		t.Error("expected error for zero area")
	}

	_, err = svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "f1", Name: "Plot", AreaHA: 10, Stage: "invalid"})
	if err == nil {
		t.Error("expected error for invalid stage")
	}
}

func TestPlotService_ListByFarm(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "farm-1", Name: "Plot A", Cultivar: "Catuaí", AreaHA: 10, PlantingYear: 2020})
	svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "farm-1", Name: "Plot B", Cultivar: "Mundo Novo", AreaHA: 15, PlantingYear: 2019})
	svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "farm-2", Name: "Plot C", Cultivar: "Catuaí", AreaHA: 5, PlantingYear: 2021})

	plots, _ := svc.ListByFarm("farm-1")
	if len(plots) != 2 {
		t.Errorf("expected 2 plots, got %d", len(plots))
	}
}

func TestPlotService_Update(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	plot, _ := svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "f1", Name: "Old", Cultivar: "Catuaí", AreaHA: 10, PlantingYear: 2020, Altitude: 900})
	plot.Name = "Updated"
	err := svc.Update(plot)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestPlotService_Update_InvalidDates(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	plot, _ := svc.Create(&entity.Plot{OrganizationID: "t1", FarmID: "f1", Name: "Plot", AreaHA: 10})

	activation := plot.CreatedAt
	deactivation := activation.AddDate(0, 0, -1)
	plot.ActivationDate = &activation
	plot.DeactivationDate = &deactivation

	if err := svc.Update(plot); err == nil {
		t.Error("expected error for deactivation date before activation date")
	}
}
