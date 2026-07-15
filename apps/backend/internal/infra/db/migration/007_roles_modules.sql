-- CafeOS - Migration 007: Roles e Modules (catálogo global configurável de
-- papéis e módulos)
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the tables
-- added in this change; it is not run by a migration tool.
--
-- roles substitui as antigas 10 constantes UserRole hardcoded em código.
-- Assim como modules, é um catálogo GLOBAL — não organization-scoped —
-- compartilhado por toda a plataforma. platform_owner e organization_admin
-- são papéis de sistema (is_system true) e não são editáveis/deletáveis;
-- todo outro papel pode ser renomeado, criado ou removido (se não estiver
-- em uso) pela tela de Papéis. O que varia por organização é o ACESSO de
-- cada papel a cada módulo (ver role_permissions), não o conjunto de
-- papéis disponíveis.
CREATE TABLE IF NOT EXISTS roles (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key         VARCHAR(64) NOT NULL UNIQUE,
    name        VARCHAR(120) NOT NULL,
    is_system   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

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
-- column, now pointing at roles.id. The table itself stays
-- organization-scoped — it's the per-org access matrix over the global
-- roles/modules catalogs above, hence the three FKs.
ALTER TABLE role_permissions ADD COLUMN IF NOT EXISTS role_id UUID;
ALTER TABLE role_permissions
    ADD CONSTRAINT fk_role_permissions_organization FOREIGN KEY (organization_id) REFERENCES organizations(id),
    ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES roles(id),
    ADD CONSTRAINT fk_role_permissions_module_catalog FOREIGN KEY (module) REFERENCES modules(key);

-- users.role_id replaces the old users.role string column.
ALTER TABLE users ADD COLUMN IF NOT EXISTS role_id UUID;
ALTER TABLE users
    ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles(id);

-- Data backfill (run once, in order, against a database still carrying the
-- legacy `role` VARCHAR columns):
--   1. INSERT INTO roles (id, key, name, is_system) VALUES
--        (uuid_generate_v4(), 'platform_owner', 'Platform Owner', true),
--        (uuid_generate_v4(), 'organization_admin', 'Admin da Organização', true),
--        (uuid_generate_v4(), 'proprietario', 'Proprietário', false),
--        (uuid_generate_v4(), 'gerente_agricola', 'Gerente Agrícola', false),
--        (uuid_generate_v4(), 'engenheiro_agronomo', 'Engenheiro Agrônomo', false),
--        (uuid_generate_v4(), 'tecnico_agricola', 'Técnico Agrícola', false),
--        (uuid_generate_v4(), 'operador_campo', 'Operador de Campo', false),
--        (uuid_generate_v4(), 'financeiro', 'Financeiro', false),
--        (uuid_generate_v4(), 'consultor_externo', 'Consultor Externo', false),
--        (uuid_generate_v4(), 'auditor', 'Auditor', false);
--   2. UPDATE role_permissions rp SET role_id = r.id FROM roles r WHERE r.key = rp.role;
--   3. UPDATE users u SET role_id = r.id FROM roles r WHERE r.key = u.role;
--   4. ALTER TABLE role_permissions ALTER COLUMN role_id SET NOT NULL, DROP COLUMN role;
--   5. ALTER TABLE users ALTER COLUMN role_id SET NOT NULL, DROP COLUMN role;
-- (In this codebase the equivalent seeding is done by
-- RoleService.SeedDefaultsIfMissing at API boot — see
-- internal/api/router.go — since there was no production data to migrate
-- at the time this feature shipped.)
