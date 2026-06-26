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
- Mobile: Flutter 3.x + drift + dio + riverpod (em migração do Expo)

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
├── mobile/           # Expo app (deprecated, migrando p/ Flutter)
│   └── src/
│       ├── api/      # HTTP client
│       ├── db/       # SQLite schema
│       ├── sync/     # Offline sync engine
│       └── screens/  # Login, Operations, PendingSync
└── mobile_flutter/   # Flutter app (:8081)
    └── lib/
        ├── api/      # Dio client + secure storage
        ├── db/       # Drift database (4 tabelas)
        ├── models/   # Operation, Plot, Farm
        ├── repos/    # CRUD local
        ├── services/ # Auth, Sync, Offline
        ├── screens/  # Login, Home, Operations, PendingSync
        ├── router/   # GoRouter
        └── shared/   # Theme + widgets
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

Apps (`mobile/` Expo, `mobile_flutter/` Flutter) compartilham mesmo contrato:

1. App salva em SQLite local (expo-sqlite / drift)
2. Enfileira em sync_queue com status "pending"
3. Na conexão, envia batch (máx 50) p/ POST /api/v1/{tenant_id}/sync
4. Backend publica no RabbitMQ por event_type
5. Worker consome e persiste no PostgreSQL
6. DLQ após 3 retentativas

### Flutter — Riverpod Providers

| Provider | Tipo | Uso |
|----------|------|-----|
| `databaseProvider` | Provider<AppDatabase> | Drift DB instance |
| `authServiceProvider` | Provider<AuthService> | Login, JWT storage |
| `offlineServiceProvider` | Provider<OfflineService> | CRUD offline + enqueue |
| `syncServiceProvider` | Provider<SyncService> | SyncAll batch |
| `loginControllerProvider` | StateNotifierProvider | Login form state |
| `operationsControllerProvider` | StateNotifierProvider | Operations list + form |
| `pendingSyncControllerProvider` | StateNotifierProvider | Sync queue viewer |

### Flutter — Estrutura de Telas v1

| Tela | Rota | Provider |
|------|------|----------|
| Login | /login | `loginControllerProvider` |
| Home (placeholder) | / (redirect p/ /operations) | — |
| Operations | /operations | `operationsControllerProvider` |
| PendingSync | /pending-sync | `pendingSyncControllerProvider` |

## Seed Credentials

| User | Email | Password | Role |
|------|-------|----------|------|
| Admin | admin@cafeos.com.br | admin123 | platform_owner |
| Proprietário | joao@cafeos.com.br | 123456 | proprietario |
| Gerente | maria@cafeos.com.br | 123456 | gerente_agricola |
| Engenheiro | carlos@cafeos.com.br | 123456 | engenheiro_agronomo |
| Operador | ana@cafeos.com.br | 123456 | operador_campo |

<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan
<!-- SPECKIT END -->

## Planos

| Plano | Status |
|-------|--------|
| `plans/plano-mobile.md` | Substituído por Flutter |
| `plans/plano-flutter-migration.md` | Ativo — migração Expo → Flutter |
