package main

import (
	"fmt"
	"log"

	"github.com/aralvesandrade/cafeos/internal/domain/entity"
	"github.com/aralvesandrade/cafeos/internal/infra/config"
	"github.com/aralvesandrade/cafeos/internal/infra/db/postgres"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	tenantName = "CafeOS Padrão"
	tenantSlug = "cafeos"
	adminName  = "Administrador"
	adminEmail = "admin@cafeos.com.br"
	adminPass  = "admin123"
	adminRole  = "platform_owner"
)

func main() {
	cfg := config.Load()
	db, err := postgres.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	defer sqlDB.Close()

	if err := seed(db); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	fmt.Println("Seed concluído com sucesso!")
	fmt.Printf("  Tenant: %s (%s)\n", tenantName, tenantSlug)
	fmt.Printf("  Admin:  %s / %s\n", adminEmail, adminPass)
}

func seed(db *gorm.DB) error {
	var tenant entity.Tenant
	err := db.Where("slug = ?", tenantSlug).First(&tenant).Error
	if err == gorm.ErrRecordNotFound {
		tenant = entity.Tenant{
			Name:      tenantName,
			Slug:      tenantSlug,
			BrandName: tenantName,
			Plan:      "pro",
		}
		if err := db.Create(&tenant).Error; err != nil {
			return fmt.Errorf("create tenant: %w", err)
		}
		fmt.Println("  ✓ Tenant criado")
	} else if err != nil {
		return fmt.Errorf("check tenant: %w", err)
	} else {
		fmt.Println("  ✓ Tenant já existe")
	}

	var userCount int64
	db.Model(&entity.User{}).Where("email = ? AND tenant_id = ?", adminEmail, tenant.ID).Count(&userCount)
	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		user := entity.User{
			TenantID:     tenant.ID,
			Name:         adminName,
			Email:        adminEmail,
			PasswordHash: string(hash),
			Role:         adminRole,
			IsActive:     true,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("create user: %w", err)
		}
		fmt.Println("  ✓ Usuário admin criado")
	} else {
		fmt.Println("  ✓ Usuário admin já existe")
	}

	return nil
}
