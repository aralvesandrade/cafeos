# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

CafeOS — multi-tenant SaaS for coffee farm management (operational, production, financial, analytics). Portuguese domain language throughout the codebase (entities, comments, docs) — keep new code consistent with that.

## Monorepo Layout

```
apps/
├── backend/   # Go REST API + RabbitMQ worker (MVP core)
├── admin/     # React + Vite admin panel (:5174)
├── frontend/  # React + Vite landing page
└── mobile/    # Expo (React Native) + SQLite offline-first app (:8081)
packages/shared/  # reserved for shared types/utils — currently empty
```

Module boundary rule: `apps/backend` never imports from the other apps; admin/frontend/mobile talk to the backend only over HTTP (mobile via the `/sync` batch endpoint).

## Commands

Local infra + services are orchestrated through `scripts/dev.sh` (run from repo root):

```bash
./scripts/dev.sh up          # docker compose (Postgres/Redis/RabbitMQ) + API in background
./scripts/dev.sh down        # stop containers
./scripts/dev.sh api         # API only, foreground (go run ./cmd/api)
./scripts/dev.sh worker      # RabbitMQ sync worker (go run ./cmd/worker/main.go)
./scripts/dev.sh admin       # admin panel dev server
./scripts/dev.sh mobile      # expo start --web
./scripts/dev.sh db:migrate  # apply schema via GORM AutoMigrate (no SQL files)
./scripts/dev.sh db:reset    # drop+recreate public schema, re-migrate
./scripts/dev.sh db:seed     # go run ./cmd/seed
./scripts/dev.sh test        # go test ./... -v (backend)
```

Docker compose file is `docker-compose.yml` at repo root (not `infra/dev/` despite what `dev.sh` output text implies).

Backend (from `apps/backend/`):
```bash
go build ./...
go test ./...                              # all tests
go test ./internal/domain/service/... -run TestName -v   # single test
go vet ./...
```

Frontend/Admin (from `apps/frontend/` or `apps/admin/`):
```bash
npm run dev      # vite dev server
npm run build     # tsc -b && vite build
npm run lint      # oxlint
```
No JS/TS test framework is configured — don't add one without asking the user first.

Mobile (from `apps/mobile/`):
```bash
npx expo start          # or npm start
npx expo start --web
```

## Backend Architecture (`apps/backend/`)

Strict layered dependency flow, inward only:

```
api/handler/ → domain/service/ → domain/repository/ (interfaces)
                                  infra/db/repository/ (GORM impl)
```

- `domain/entity/` — GORM entities. UUID PK, `OrganizationID` FK on every organization-scoped entity. No external deps.
- `domain/repository/` — Go interfaces only, no GORM imports. Includes `Transactor` for cross-repo atomicity.
- `domain/service/` — business logic + the Rule Engine; depends on repository interfaces, never on `infra/db/repository` implementations directly.
- `infra/db/repository/` — GORM implementations (`WithTx` variants for transactions).
- `api/handler/` — thin: parse request → call service → `writeJSON()`/`writeError()`.
- `api/middleware/` — auth (JWT), RBAC (`RequireRole`), organization resolution (`middleware.OrganizationIDKey`), CORS.
- `event/` — in-memory event bus, decoupled from messaging infra (RabbitMQ is separate, in `infra/messaging/`).
- `cmd/api` — API entrypoint (:5001). `cmd/worker` — RabbitMQ consumer that persists to Postgres. `cmd/seed` — seed data.

Every handler pulls the organization ID from context via `middleware.OrganizationIDKey`; routes live under `/api/v1/{organization_id}/...`. Admin-only routes (`/api/v1/admin/...`) require `platform_owner` via `RequireRole`.

Transactions: use `Transactor.RunInTx` when a flow touches multiple repos atomically (e.g. `HarvestService` updates harvest + recalculates indicators in one tx). Most other services operate on individual repositories without a tx wrapper.

### RBAC roles (10)
`platform_owner`, `organization_admin`, `proprietario`, `gerente_agricola`, `engenheiro_agronomo`, `tecnico_agricola`, `operador_campo`, `financeiro`, `consultor_externo`, `auditor`.

### Rule Engine
Configurable alerting (`RuleEngine.AddRule()`) — e.g. low productivity (<25 sacas/ha), high cost (>R$400/saca) — emits `AlertGenerated` events.

### Events
In-memory bus emits: `OperationRegistered`, `HarvestFinalized`, `IndicatorUpdated`, `AlertGenerated`.

## Mobile Offline Sync (`apps/mobile/`)

1. Writes land in local SQLite first; enqueued into `sync_queue` with status `pending`.
2. On connectivity (NetInfo), client POSTs a batch (max 50) to `/sync`.
3. Backend publishes each item to RabbitMQ.
4. `cmd/worker` consumes and persists to PostgreSQL; failed items go to a DLQ after 3 retries.

Structure: `src/api/` (HTTP client), `src/db/` (SQLite schema/migrations), `src/sync/` (sync engine), `src/screens/`, `src/navigation/`.

## Naming Conventions

| Scope | Convention | Example |
|-------|-----------|---------|
| Go files | `snake_case.go` | `farm_handler.go` |
| Go exported types | PascalCase | `FarmHandler` |
| Go JSON tags | `snake_case` | `json:"total_area_ha"` |
| TS/TSX components (admin/frontend) | PascalCase | `FarmDetail.tsx` |
| TS utility files | camelCase | `api.ts` |
| Mobile screens | PascalCase + `Screen.tsx` | `LoginScreen.tsx` |
| Branch names | `feat/*`, `fix/*`, `docs/*`, `chore/*`, `style/*` |
| Commits | Conventional Commits: `type: description` |

## Agent Rules (from AGENTS.md)

- Commits and pushes are NEVER automatic — always ask before each commit/push.
- Feature plans go in `/plans/` at repo root (not `.opencode/` or elsewhere); present the plan for approval before implementing.
- Any destructive operation (delete, reset, migration) needs explicit user confirmation first.

## Seed Credentials (local dev only)

| User | Email | Password | Role |
|------|-------|----------|------|
| Admin | admin@cafeos.com.br | admin123 | platform_owner |
| Proprietário (2 fazendas) | joao@cafeos.com.br | 123456 | proprietario |
| Proprietário (1 fazenda) | maria@cafeos.com.br | 123456 | proprietario |
| Engenheiro | carlos@cafeos.com.br | 123456 | engenheiro_agronomo |
| Operador | ana@cafeos.com.br | 123456 | operador_campo |
| Admin Organização | fernanda@cafeos.com.br | 123456 | organization_admin |
| Consultor | rodrigo@cafeos.com.br | 123456 | consultor_externo |

## Spec-kit

Project uses spec-kit (brownfield extension) for AI-assisted development. Constitution lives at `.specify/memory/constitution.md` (source of truth for architecture rules, superset of this file) — check it if this file and the code disagree. Relevant commands: `/speckit.specify`, `/speckit.clarify`, `/speckit.plan`, `/speckit.tasks`, `/speckit.checklist`, `/speckit.analyze`.
