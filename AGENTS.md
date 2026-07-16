# CafeOS — Agent Guide

## Agent Rules

- Commits and push NEVER automatic. Ask user before each commit/push.
- Plans must be saved in `/plans/` folder at project root, not in `.opencode/` or elsewhere.
- When user says "commitar" or "push", confirm and execute.
- When planning features, create a markdown file in `/plans/` first and present to user for approval before implementing.
- Any destructive operation (delete, reset, migration) must be confirmed by user first.

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
│   └── src/pages/    # Dashboard, Farms, Plots, Operations, Harvests, Financial, Stock, Fleet, Labor, Organizations, Users
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

- `platform_owner` — full access + admin (organizations, users)
- `organization_admin` — organization config
- `proprietario` — farm owner
- `gerente_agricola` — agricultural manager
- `engenheiro_agronomo` — agronomic engineer
- `tecnico_agricola` — agricultural technician
- `operador_campo` — field operator
- `financeiro` — financial
- `consultor_externo` — external consultant
- `auditor` — auditor

## Routes

### Multi-tenant (`/api/v1/{organization_id}`)
- CRUD: farms, plots, operations, harvests, financial, stock/*, fleet/*, labor/*
- GET: dashboard, agricultural-products, operations/recent
- POST: sync

### Admin (`/api/v1/admin`, platform_owner only)
- CRUD: organizations, users

### Auth
- POST /auth/login

## Key Conventions

- Entities: `domain/entity/`, GORM tags, UUID PK, OrganizationID FK
- Repos: interface in `domain/repository/`, impl in `infra/db/repository/`
- Handlers: `api/handler/`, use `writeJSON()` and `writeError()`
- Services: optional business logic layer
- Every handler gets organizationID from context via `middleware.OrganizationIDKey`
- Admin routes use `RequireRole(entity.RolePlatformOwner)`
- Frontend API calls use `apiRequest()` from `lib/api.ts`
- Admin API calls use `{ admin: true }` to skip organization_id URL prefix

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
| Gerente Agrícola | maria@cafeos.com.br | 123456 | gerente_agricola |
| Engenheiro | carlos@cafeos.com.br | 123456 | engenheiro_agronomo |
| Operador | ana@cafeos.com.br | 123456 | operador_campo |
| Admin Organização | fernanda@cafeos.com.br | 123456 | organization_admin |
| Consultor | rodrigo@cafeos.com.br | 123456 | consultor_externo |

<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan
<!-- SPECKIT END -->

## Planos

| Plano | Status |
|-------|--------|
| `plans/plano-mobile.md` | Ativo |
| `plans/plano-flutter-migration.md` | Cancelado |
