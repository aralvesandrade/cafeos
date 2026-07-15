package entity

import "time"

type Producer struct {
	ID             string  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	OrganizationID string  `json:"organization_id" gorm:"type:uuid;not null;index"`
	FarmID         string  `json:"farm_id" gorm:"type:uuid;not null;uniqueIndex"`
	UserID         *string `json:"user_id" gorm:"type:uuid;index"`

	CPF           string     `json:"cpf" gorm:"default:''"`
	Name          string     `json:"name" gorm:"not null"`
	RG            string     `json:"rg" gorm:"default:''"`
	IssuingBody   string     `json:"issuing_body" gorm:"default:''"`
	Gender        string     `json:"gender" gorm:"default:''"`
	BirthDate     *time.Time `json:"birth_date"`
	MaritalStatus string     `json:"marital_status" gorm:"default:''"`
	Phone         string     `json:"phone" gorm:"default:''"`
	Email         string     `json:"email" gorm:"default:''"`
	Education     string     `json:"education" gorm:"default:''"`

	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Organization Organization `json:"-" gorm:"foreignKey:OrganizationID"`
	Farm         Farm         `json:"-" gorm:"foreignKey:FarmID"`
	User         *User        `json:"-" gorm:"foreignKey:UserID"`
}

func (Producer) TableName() string {
	return "producers"
}
