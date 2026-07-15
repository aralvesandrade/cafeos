-- CafeOS - Migration 005: RolePermission (matriz de acesso por papel)
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the table
-- added in this change; it is not run by a migration tool.
--
-- role_permissions guarda, por organização, o nível de acesso
-- (none|read|write) de cada UserRole em cada módulo de tela (dashboard,
-- farms, operations, harvests, resources, financial, permissions).
-- Configurável via tela de Permissões (organization_admin/proprietario);
-- linhas são seedadas com valores padrão na criação da organização.

CREATE TABLE IF NOT EXISTS role_permissions (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NOT NULL,
    role            VARCHAR(32) NOT NULL,
    module          VARCHAR(32) NOT NULL,
    access          VARCHAR(16) NOT NULL DEFAULT 'none',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT idx_org_role_module UNIQUE (organization_id, role, module)
);
