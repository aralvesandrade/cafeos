package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// PlanFeature é um item da lista de features de um plano, com um marcador
// de incluído/não incluído — replica o padrão visual (check/x) usado na
// landing page, agora vindo da tabela em vez de hardcoded no componente.
type PlanFeature struct {
	Label    string `json:"label"`
	Included bool   `json:"included"`
}

// PlanFeatureList é uma lista de PlanFeature persistida como JSON em uma
// coluna jsonb.
type PlanFeatureList []PlanFeature

func (l PlanFeatureList) Value() (driver.Value, error) {
	if l == nil {
		return "[]", nil
	}
	return json.Marshal(l)
}

func (l *PlanFeatureList) Scan(value any) error {
	if value == nil {
		*l = PlanFeatureList{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		s, ok := value.(string)
		if !ok {
			return errors.New("unsupported type for PlanFeatureList")
		}
		b = []byte(s)
	}
	return json.Unmarshal(b, l)
}

type Plan struct {
	ID              string          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name            string          `json:"name" gorm:"not null"`
	Slug            string          `json:"slug" gorm:"uniqueIndex;not null"`
	Description     string          `json:"description" gorm:"default:''"`
	PriceCents      int64           `json:"price_cents" gorm:"not null;default:0"`
	BillingInterval string          `json:"billing_interval" gorm:"not null;default:'monthly'"`
	MaxFarms        int             `json:"max_farms" gorm:"not null;default:0"`
	MaxUsers        int             `json:"max_users" gorm:"not null;default:0"`
	Features        PlanFeatureList `json:"features" gorm:"type:jsonb"`
	Active          bool            `json:"active" gorm:"not null;default:true"`
	Featured        bool            `json:"featured" gorm:"not null;default:false"`
	DisplayOrder    int             `json:"display_order" gorm:"not null;default:0"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
