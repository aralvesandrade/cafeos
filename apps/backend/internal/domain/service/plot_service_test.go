package service

import (
	"testing"

	mock "github.com/aralvesandrade/cafeos/internal/domain/service/testing"
)

func TestPlotService_Create(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	plot, err := svc.Create("t1", "farm-1", "Talhão A", "Catuaí", "Argiloso", 10.5, 2020, 950)
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
}

func TestPlotService_Create_Validation(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	_, err := svc.Create("t1", "f1", "", "Catuaí", "Argiloso", 10, 2020, 950)
	if err == nil {
		t.Error("expected error for empty name")
	}

	_, err = svc.Create("t1", "f1", "Plot", "Catuaí", "Argiloso", 0, 2020, 950)
	if err == nil {
		t.Error("expected error for zero area")
	}
}

func TestPlotService_ListByFarm(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	svc.Create("t1", "farm-1", "Plot A", "Catuaí", "", 10, 2020, 0)
	svc.Create("t1", "farm-1", "Plot B", "Mundo Novo", "", 15, 2019, 0)
	svc.Create("t1", "farm-2", "Plot C", "Catuaí", "", 5, 2021, 0)

	plots, _ := svc.ListByFarm("farm-1")
	if len(plots) != 2 {
		t.Errorf("expected 2 plots, got %d", len(plots))
	}
}

func TestPlotService_Update(t *testing.T) {
	repo := mock.NewInMemoryPlotRepo()
	svc := NewPlotService(repo)

	plot, _ := svc.Create("t1", "f1", "Old", "Catuaí", "", 10, 2020, 900)
	plot.Name = "Updated"
	err := svc.Update(plot)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
