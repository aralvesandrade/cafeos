package entity

import "time"

type PlotStage string

const (
	PlotStageFormacao PlotStage = "formacao"
	PlotStageProducao PlotStage = "producao"
)

type Plot struct {
	ID             string  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string  `json:"organization_id" gorm:"type:uuid;not null;index"`
	FarmID         string  `json:"farm_id" gorm:"type:uuid;not null;index"`
	Name           string  `json:"name" gorm:"not null"`
	AreaHA         float64 `json:"area_ha" gorm:"type:numeric(12,2);default:0"`
	Cultivar       string  `json:"cultivar" gorm:"default:''"`
	PlantingYear   int     `json:"planting_year" gorm:"default:0"`
	Altitude       int     `json:"altitude" gorm:"default:0"`
	SoilType       string  `json:"soil_type" gorm:"default:''"`

	Leased     bool      `json:"leased" gorm:"default:false"`
	Stage      PlotStage `json:"stage" gorm:"default:'formacao'"`
	Irrigation string    `json:"irrigation" gorm:"default:''"`

	ActivationDate   *time.Time `json:"activation_date"`
	PlantingDate     *time.Time `json:"planting_date"`
	DeactivationDate *time.Time `json:"deactivation_date"`

	Intercropped  bool   `json:"intercropped" gorm:"default:false"`
	SecondaryCrop string `json:"secondary_crop" gorm:"default:''"`
	Notes         string `json:"notes" gorm:"default:''"`

	CropType string `json:"crop_type" gorm:"default:''"`

	FormationCostPerHA float64 `json:"formation_cost_per_ha" gorm:"type:numeric(14,2);default:0"`
	UsefulLifeYears    int     `json:"useful_life_years" gorm:"default:0"`
	RowSpacingM        float64 `json:"row_spacing_m" gorm:"type:numeric(6,2);default:0"`
	PlantSpacingM      float64 `json:"plant_spacing_m" gorm:"type:numeric(6,2);default:0"`
	PlantCount         int     `json:"plant_count" gorm:"default:0"`

	DamAreaHA          float64 `json:"dam_area_ha" gorm:"type:numeric(12,2);default:0"`
	ImprovementsAreaHA float64 `json:"improvements_area_ha" gorm:"type:numeric(12,2);default:0"`
	RoadsAreaHA        float64 `json:"roads_area_ha" gorm:"type:numeric(12,2);default:0"`
	APPAreaHA          float64 `json:"app_area_ha" gorm:"type:numeric(12,2);default:0"`
	LegalReserveAreaHA float64 `json:"legal_reserve_area_ha" gorm:"type:numeric(12,2);default:0"`

	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Farm         Farm         `json:"-" gorm:"foreignKey:FarmID"`
}

func (Plot) TableName() string {
	return "plots"
}
