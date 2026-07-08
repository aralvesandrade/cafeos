-- CafeOS - Migration 003: Plot (talhão) detail fields
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the columns
-- added in this change; it is not run by a migration tool.

ALTER TABLE plots
    ADD COLUMN IF NOT EXISTS leased                 BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS stage                  VARCHAR(20) NOT NULL DEFAULT 'formacao',
    ADD COLUMN IF NOT EXISTS irrigation              VARCHAR(64) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS activation_date          DATE,
    ADD COLUMN IF NOT EXISTS planting_date            DATE,
    ADD COLUMN IF NOT EXISTS deactivation_date        DATE,
    ADD COLUMN IF NOT EXISTS intercropped            BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS secondary_crop           VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS notes                   TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS crop_type                VARCHAR(64) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS formation_cost_per_ha     NUMERIC(14,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS useful_life_years         INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS row_spacing_m             NUMERIC(6,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS plant_spacing_m           NUMERIC(6,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS plant_count               INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS dam_area_ha               NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS improvements_area_ha      NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS roads_area_ha             NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS app_area_ha               NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS legal_reserve_area_ha     NUMERIC(12,2) NOT NULL DEFAULT 0;
