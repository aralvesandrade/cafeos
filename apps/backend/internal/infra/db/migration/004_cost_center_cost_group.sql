-- CafeOS - Migration 004: CostCenter SENAR/CEPEA cost group classification
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the columns
-- added in this change; it is not run by a migration tool.
--
-- cost_group classifies a despesa cost center into the SENAR/CEPEA cost
-- hierarchy (operacional_efetivo | mao_de_obra_familiar | capital_depreciacao
-- | remuneracao_capital), used to compute COE/COT/CT harvest indicators.
-- Existing cost centers are left with cost_group = '' (unclassified) and are
-- excluded from COE/COT/CT until reclassified — no existing data is altered
-- or deleted by this migration. The fixed SENAR category catalog itself
-- (entity.SenarCostCategories) is a Go-level constant, not a DB table.

ALTER TABLE cost_centers
    ADD COLUMN IF NOT EXISTS cost_group VARCHAR(32) NOT NULL DEFAULT '';
