package entity

import "time"

type Farm struct {
	ID             string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string `json:"organization_id" gorm:"type:uuid;not null;index"`
	Name           string `json:"name" gorm:"not null"`
	Owner          string `json:"owner" gorm:"default:''"`
	Location       string `json:"location" gorm:"default:''"`

	Phone                    string `json:"phone" gorm:"default:''"`
	Activities               string `json:"activities" gorm:"default:''"`
	MainCrop                 string `json:"main_crop" gorm:"default:''"`
	SecondaryCrop            string `json:"secondary_crop" gorm:"default:''"`
	State                    string `json:"state" gorm:"default:''"`
	City                     string `json:"city" gorm:"default:''"`
	Address                  string `json:"address" gorm:"default:''"`
	ProductionSystem         string `json:"production_system" gorm:"default:''"`
	CommercializationProduct string `json:"commercialization_product" gorm:"default:''"`

	HasNoCNPJ              bool   `json:"has_no_cnpj" gorm:"default:false"`
	CNPJ                   string `json:"cnpj" gorm:"default:''"`
	HasNoNIRF              bool   `json:"has_no_nirf" gorm:"default:false"`
	NIRF                   string `json:"nirf" gorm:"default:''"`
	HasNoINCRA             bool   `json:"has_no_incra" gorm:"default:false"`
	INCRARegistration      string `json:"incra_registration" gorm:"default:''"`
	HasNoStateRegistration bool   `json:"has_no_state_registration" gorm:"default:false"`
	StateRegistration      string `json:"state_registration" gorm:"default:''"`
	HasNoDAP               bool   `json:"has_no_dap" gorm:"default:false"`
	DAP                    string `json:"dap" gorm:"default:''"`
	HasNoCAR               bool   `json:"has_no_car" gorm:"default:false"`
	CAR                    string `json:"car" gorm:"default:''"`

	FullyLeased    bool    `json:"fully_leased" gorm:"default:false"`
	LandValuePerHA float64 `json:"land_value_per_ha" gorm:"type:numeric(14,2);default:0"`

	TotalAreaHA   float64 `json:"total_area_ha" gorm:"type:numeric(12,2);default:0"`
	PlantedAreaHA float64 `json:"planted_area_ha" gorm:"type:numeric(12,2);default:0"`

	DamAreaHA                   float64 `json:"dam_area_ha" gorm:"type:numeric(12,2);default:0"`
	ImprovementsAreaHA          float64 `json:"improvements_area_ha" gorm:"type:numeric(12,2);default:0"`
	RoadsAreaHA                 float64 `json:"roads_area_ha" gorm:"type:numeric(12,2);default:0"`
	APPAreaHA                   float64 `json:"app_area_ha" gorm:"type:numeric(12,2);default:0"`
	LegalReserveAreaHA          float64 `json:"legal_reserve_area_ha" gorm:"type:numeric(12,2);default:0"`
	NativeVegetationAreaHA      float64 `json:"native_vegetation_area_ha" gorm:"type:numeric(12,2);default:0"`
	LivestockAreaNotCoveredHA   float64 `json:"livestock_area_not_covered_ha" gorm:"type:numeric(12,2);default:0"`
	AgricultureAreaNotCoveredHA float64 `json:"agriculture_area_not_covered_ha" gorm:"type:numeric(12,2);default:0"`
	NonAgriculturalAreaHA       float64 `json:"non_agricultural_area_ha" gorm:"type:numeric(12,2);default:0"`

	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Plots        []Plot       `json:"plots,omitempty" gorm:"foreignKey:FarmID"`
	Producer     *Producer    `json:"producer,omitempty" gorm:"foreignKey:FarmID"`
}

func (Farm) TableName() string {
	return "farms"
}

// ProductiveAreaAvailableHA returns the total area minus the non-productive
// area breakdown (dam, improvements, roads, APP, legal reserve, native
// vegetation and other non-agricultural areas), mirroring the "Área
// Produtiva Disponível" calculated field from the reference layout.
func (f *Farm) ProductiveAreaAvailableHA() float64 {
	return f.TotalAreaHA -
		f.DamAreaHA -
		f.ImprovementsAreaHA -
		f.RoadsAreaHA -
		f.APPAreaHA -
		f.LegalReserveAreaHA -
		f.NativeVegetationAreaHA -
		f.LivestockAreaNotCoveredHA -
		f.AgricultureAreaNotCoveredHA -
		f.NonAgriculturalAreaHA
}
