-- CafeOS - Migration 006: Plan (catálogo de planos de assinatura)
-- NOTE: schema is authoritatively managed by GORM AutoMigrate (see
-- internal/infra/db/postgres/connection.go). This file exists for parity
-- with the docker-compose bootstrap and as documentation of the table
-- added in this change; it is not run by a migration tool.
--
-- plans guarda o catálogo de planos de assinatura SaaS (nome, preço,
-- intervalo de cobrança, limites e features). Gerenciado via painel
-- admin (somente platform_owner). organizations.plan_id referencia
-- plans.id; organizations.plan (texto livre legado) é preenchido via
-- backfill a partir do slug correspondente.

CREATE TABLE IF NOT EXISTS plans (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name             VARCHAR(255) NOT NULL,
    slug             VARCHAR(100) NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    price_cents      BIGINT NOT NULL DEFAULT 0,
    billing_interval VARCHAR(16) NOT NULL DEFAULT 'monthly',
    max_farms        INTEGER NOT NULL DEFAULT 0,
    max_users        INTEGER NOT NULL DEFAULT 0,
    features         JSONB NOT NULL DEFAULT '[]',
    active           BOOLEAN NOT NULL DEFAULT true,
    featured         BOOLEAN NOT NULL DEFAULT false,
    display_order    INTEGER NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_plans_slug UNIQUE (slug)
);

ALTER TABLE organizations ADD COLUMN IF NOT EXISTS plan_id UUID;
