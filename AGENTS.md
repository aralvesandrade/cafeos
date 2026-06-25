# CafeOS — Agent Guide

## Stack

- Backend: Go 1.26 + GORM + PostgreSQL + Redis + RabbitMQ
- Admin: React 19 + TypeScript + Vite + Tailwind v4 + Recharts
- Frontend: React + Vite + Tailwind v4 (landing page)
- Mobile: Expo (React Native) + SQLite

## Project Structure

```
apps/
├── backend/          # Go API + Worker
│   ├── cmd/api/      # API entry point (:5001)
│   ├── cmd/worker/   # RabbitMQ worker
│   ├── internal/
│   │   ├── api/handler/   # HTTP handlers
│   │   ├── api/middleware/ # Auth, RBAC, CORS
│   │   ├── domain/entity/  # GORM entities
│   │   ├── domain/service/ # Business logic
│   │   ├── domain/repository/ # Repo interfaces
│   │   ├── infra/db/repository/ # GORM implementations
│   │   ├── infra/messaging/     # RabbitMQ pub/sub
│   │   └── event/               # In-memory event bus
│   └── docs/         # Swagger
├── admin/            # Admin panel (:5174)
│   └── src/pages/    # Dashboard, Farms, Plots, Operations, Harvests, Financial, Stock, Fleet, Labor, Tenants, Users
├── frontend/         # Landing page
└── mobile/           # Expo app (:8081)
    └── src/
        ├── api/      # HTTP client
        ├── db/       # SQLite schema
        ├── sync/     # Offline sync engine
        └── screens/  # Login, Operations, PendingSync
```

## Ports

| Service | Port |
|---------|------|
| API | 5001 |
| Admin panel | 5174 |
| Mobile (Expo web) | 8081 |
| PostgreSQL | 5432 |
| Redis | 6379 |
| RabbitMQ | 5672 |
| RabbitMQ UI | 15672 |

## RBAC Roles

- `platform_owner` — full access + admin (tenants, users)
- `tenant_admin` — tenant config
- `proprietario` — farm owner
- `gerente_agricola` — agricultural manager
- `engenheiro_agronomo` — agronomic engineer
- `tecnico_agricola` — agricultural technician
- `operador_campo` — field operator
- `financeiro` — financial
- `consultor_externo` — external consultant
- `auditor` — auditor

## Routes

### Multi-tenant (`/api/v1/{tenant_id}`)
- CRUD: farms, plots, operations, harvests, financial, stock/*, fleet/*, labor/*
- GET: dashboard, agricultural-products, operations/recent
- POST: sync

### Admin (`/api/v1/admin`, platform_owner only)
- CRUD: tenants, users

### Auth
- POST /auth/login

## Key Conventions

- Entities: `domain/entity/`, GORM tags, UUID PK, TenantID FK
- Repos: interface in `domain/repository/`, impl in `infra/db/repository/`
- Handlers: `api/handler/`, use `writeJSON()` and `writeError()`
- Services: optional business logic layer
- Every handler gets tenantID from context via `middleware.TenantIDKey`
- Admin routes use `RequireRole(entity.RolePlatformOwner)`
- Frontend API calls use `apiRequest()` from `lib/api.ts`
- Admin API calls use `{ admin: true }` to skip tenant_id URL prefix

## Mobile Offline Sync

1. App saves to local SQLite
2. Enqueues in sync_queue with status "pending"
3. On connectivity (NetInfo), sends batch (max 50) to POST /sync
4. Backend publishes to RabbitMQ
5. Worker consumes and persists to PostgreSQL
6. DLQ after 3 retries

## Seed Credentials

| User | Email | Password | Role |
|------|-------|----------|------|
| Admin | admin@cafeos.com.br | admin123 | platform_owner |
| Proprietário | joao@cafeos.com.br | 123456 | proprietario |
| Gerente | maria@cafeos.com.br | 123456 | gerente_agricola |
| Engenheiro | carlos@cafeos.com.br | 123456 | engenheiro_agronomo |
| Operador | ana@cafeos.com.br | 123456 | operador_campo |
