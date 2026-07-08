-- CafeOS - Migration 002: Farm documentation/address/area breakdown fields + Producer
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the columns
-- added in this change; it is not run by a migration tool.

-- ============================================================
-- FARMS: new columns
-- ============================================================
ALTER TABLE farms
    ADD COLUMN IF NOT EXISTS phone                          VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS activities                      VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS main_crop                       VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS secondary_crop                  VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS state                           VARCHAR(10)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS city                            VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS address                         TEXT         NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS production_system                VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS commercialization_product        VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS has_no_cnpj                     BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS cnpj                             VARCHAR(32)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS has_no_nirf                     BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS nirf                             VARCHAR(32)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS has_no_incra                    BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS incra_registration               VARCHAR(64)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS has_no_state_registration        BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS state_registration               VARCHAR(64)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS has_no_dap                      BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS dap                              VARCHAR(64)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS has_no_car                      BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS car                              VARCHAR(64)  NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS fully_leased                    BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS land_value_per_ha                NUMERIC(14,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS dam_area_ha                      NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS improvements_area_ha             NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS roads_area_ha                    NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS app_area_ha                      NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS legal_reserve_area_ha            NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS native_vegetation_area_ha        NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS livestock_area_not_covered_ha    NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS agriculture_area_not_covered_ha  NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS non_agricultural_area_ha         NUMERIC(12,2) NOT NULL DEFAULT 0;

-- ============================================================
-- PRODUCERS (Produtor responsável pela propriedade, 1:1 com farms)
-- ============================================================
CREATE TABLE IF NOT EXISTS producers (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    farm_id        UUID NOT NULL UNIQUE REFERENCES farms(id) ON DELETE CASCADE,
    cpf            VARCHAR(32)  NOT NULL DEFAULT '',
    name           VARCHAR(255) NOT NULL,
    rg             VARCHAR(32)  NOT NULL DEFAULT '',
    issuing_body   VARCHAR(32)  NOT NULL DEFAULT '',
    gender         VARCHAR(20)  NOT NULL DEFAULT '',
    birth_date     DATE,
    marital_status VARCHAR(32)  NOT NULL DEFAULT '',
    phone          VARCHAR(32)  NOT NULL DEFAULT '',
    email          VARCHAR(255) NOT NULL DEFAULT '',
    education      VARCHAR(64)  NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_producers_tenant ON producers(tenant_id);
