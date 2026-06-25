-- CafeOS - Initial Schema
-- Migration 001: Core entities for MVP

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- TENANTS
-- ============================================================
CREATE TABLE IF NOT EXISTS tenants (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    brand_name  VARCHAR(255) NOT NULL DEFAULT '',
    logo_url    TEXT NOT NULL DEFAULT '',
    primary_color VARCHAR(50) NOT NULL DEFAULT '#2E7D32',
    plan        VARCHAR(50) NOT NULL DEFAULT 'free',
    domain      VARCHAR(255) NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- USERS
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(50) NOT NULL DEFAULT 'operador_campo',
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id);

-- ============================================================
-- FARMS
-- ============================================================
CREATE TABLE IF NOT EXISTS farms (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id      UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name           VARCHAR(255) NOT NULL,
    owner          VARCHAR(255) NOT NULL DEFAULT '',
    location       TEXT NOT NULL DEFAULT '',
    total_area_ha  NUMERIC(12,2) NOT NULL DEFAULT 0,
    planted_area_ha NUMERIC(12,2) NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_farms_tenant ON farms(tenant_id);

-- ============================================================
-- PLOTS (Talhões)
-- ============================================================
CREATE TABLE IF NOT EXISTS plots (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    farm_id       UUID NOT NULL REFERENCES farms(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    area_ha       NUMERIC(12,2) NOT NULL DEFAULT 0,
    cultivar      VARCHAR(255) NOT NULL DEFAULT '',
    planting_year INTEGER NOT NULL DEFAULT 0,
    altitude      INTEGER NOT NULL DEFAULT 0,
    soil_type     VARCHAR(100) NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_plots_tenant ON plots(tenant_id);
CREATE INDEX idx_plots_farm ON plots(farm_id);

-- ============================================================
-- AGRICULTURAL PRODUCTS
-- ============================================================
CREATE TABLE IF NOT EXISTS agricultural_products (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id  UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    type       VARCHAR(50) NOT NULL DEFAULT 'outro',
    unit       VARCHAR(50) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agri_products_tenant ON agricultural_products(tenant_id);

-- ============================================================
-- OPERATIONS
-- ============================================================
CREATE TABLE IF NOT EXISTS operations (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id    UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    plot_id      UUID NOT NULL REFERENCES plots(id) ON DELETE CASCADE,
    type         VARCHAR(50) NOT NULL,
    date         TIMESTAMPTZ NOT NULL,
    responsible  VARCHAR(255) NOT NULL DEFAULT '',
    product_used VARCHAR(255) NOT NULL DEFAULT '',
    quantity     NUMERIC(12,2) NOT NULL DEFAULT 0,
    cost         NUMERIC(12,2) NOT NULL DEFAULT 0,
    notes        TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_operations_tenant ON operations(tenant_id);
CREATE INDEX idx_operations_plot ON operations(plot_id);
CREATE INDEX idx_operations_date ON operations(date);

-- ============================================================
-- HARVESTS (Safras)
-- ============================================================
CREATE TABLE IF NOT EXISTS harvests (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id           UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    year                INTEGER NOT NULL,
    description         TEXT NOT NULL DEFAULT '',
    estimated_production NUMERIC(12,2) NOT NULL DEFAULT 0,
    status              VARCHAR(50) NOT NULL DEFAULT 'planejada',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, year)
);

CREATE INDEX idx_harvests_tenant ON harvests(tenant_id);

-- ============================================================
-- HARVEST PRODUCTION
-- ============================================================
CREATE TABLE IF NOT EXISTS harvest_productions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id   UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    harvest_id  UUID NOT NULL REFERENCES harvests(id) ON DELETE CASCADE,
    plot_id     UUID NOT NULL REFERENCES plots(id) ON DELETE CASCADE,
    quantity    NUMERIC(12,2) NOT NULL DEFAULT 0,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes       TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_hp_tenant ON harvest_productions(tenant_id);
CREATE INDEX idx_hp_harvest ON harvest_productions(harvest_id);
CREATE INDEX idx_hp_plot ON harvest_productions(plot_id);

-- ============================================================
-- INDICATORS
-- ============================================================
CREATE TABLE IF NOT EXISTS indicators (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id     UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    harvest_id    UUID NOT NULL REFERENCES harvests(id) ON DELETE CASCADE,
    plot_id       UUID REFERENCES plots(id) ON DELETE SET NULL,
    type          VARCHAR(100) NOT NULL,
    value         NUMERIC(14,2) NOT NULL DEFAULT 0,
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_indicators_tenant ON indicators(tenant_id);
CREATE INDEX idx_indicators_harvest ON indicators(harvest_id);
CREATE INDEX idx_indicators_type ON indicators(type);
