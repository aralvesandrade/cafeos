-- CafeOS - Migration 007: Roles e Modules (cadastro configurável de papéis
-- e catálogo de módulos)
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the tables
-- added in this change; it is not run by a migration tool.
--
-- roles substitui as antigas 10 constantes UserRole hardcoded em código.
-- platform_owner e organization_admin são papéis de sistema
-- (organization_id NULL, is_system true), compartilhados por todas as
-- organizações e não editáveis/deletáveis. Todo outro papel pertence a uma
-- única organização: os oito papéis padrão (proprietario, gerente_agricola,
-- ...) são semeados como "kit inicial" editável na criação da organização,
-- e admins podem criar papéis adicionais (ex: "colhedor_chefe") pela tela
-- de Papéis.
CREATE TABLE IF NOT EXISTS roles (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID NULL,
    key             VARCHAR(64) NOT NULL,
    name            VARCHAR(120) NOT NULL,
    is_system       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- Partial unique indexes: a key must be unique among an organization's own
-- roles, and unique among the (organization-less) global system roles.
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_org_key
    ON roles (organization_id, key) WHERE organization_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_system_key
    ON roles (key) WHERE organization_id IS NULL;

-- modules is a global catalog of the application's fixed screens/areas.
-- Route wiring in router.go still references the module key directly
-- (adding a module still requires code + deploy) — this table only stores
-- display metadata (name, order) so the admin UI has a single source of
-- truth instead of duplicating labels in Go and TS.
CREATE TABLE IF NOT EXISTS modules (
    key         VARCHAR(50) PRIMARY KEY,
    name        VARCHAR(120) NOT NULL,
    "order"     INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- role_permissions.role_id replaces the old role_permissions.role string
-- column, now pointing at roles.id.
ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS role_id UUID;
-- users.role_id replaces the old users.role string column.
ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id UUID;

-- Data backfill (run once, in order, against a database still carrying the
-- legacy `role`/`role` VARCHAR columns):
--   1. INSERT INTO roles (id, key, name, is_system) VALUES
--        (uuid_generate_v4(), 'platform_owner', 'Platform Owner', true),
--        (uuid_generate_v4(), 'organization_admin', 'Admin da Organização', true);
--   2. For each organization, INSERT INTO roles (id, organization_id, key, name)
--      one row per legacy UserRole key (proprietario, gerente_agricola,
--      engenheiro_agronomo, tecnico_agricola, operador_campo, financeiro,
--      consultor_externo, auditor).
--   3. UPDATE role_permissions rp SET role_id = r.id FROM roles r
--        WHERE r.key = rp.role AND (r.organization_id = rp.organization_id OR r.organization_id IS NULL);
--   4. UPDATE users u SET role_id = r.id FROM roles r
--        WHERE r.key = u.role AND (r.organization_id = u.organization_id OR r.organization_id IS NULL);
--   5. ALTER TABLE role_permissions ALTER COLUMN role_id SET NOT NULL, DROP COLUMN role;
--   6. ALTER TABLE users ALTER COLUMN role_id SET NOT NULL, DROP COLUMN role;
-- (In this codebase the equivalent seeding is done by
-- RoleService.SeedSystemRolesIfMissing/SeedDefaultsIfMissing at API boot —
-- see internal/api/router.go — since there was no production data to
-- migrate at the time this feature shipped.)
