# CafeOS Constitution

## Core Principles

### I. Monorepo Module Boundaries

Code lives in `apps/` modules. Each module is independent:
- `apps/backend/` — Go REST API + Worker. Never import from `apps/admin/`, `apps/frontend/`, or `apps/mobile/`.
- `apps/admin/` — React admin panel. Communicates with backend via HTTP only.
- `apps/frontend/` — React landing page. Communicates with backend via HTTP only.
- `apps/mobile/` — Expo React Native app. Communicates with backend via sync endpoint only.
- `packages/shared/` — Shared types/utilities. Must have zero runtime dependencies on other modules.

### II. Backend Layered Architecture

Backend follows strict layering. Dependencies flow inward only:

```
api/handler/ → domain/service/ → domain/repository/ (interface)
                                  infra/db/repository/ (impl GORM)
```

- `domain/entity/` — GORM entities, UUID PK, TenantID FK. No external dependencies.
- `domain/repository/` — Go interfaces only. No GORM imports.
- `domain/service/` — Business logic. Depends on repository interfaces, not implementations.
- `infra/db/repository/` — GORM implementations of repository interfaces.
- `api/handler/` — HTTP handlers. Thin layer: parse request → call service → write response.
- `api/middleware/` — Auth, RBAC, Tenant resolution, CORS.
- `event/` — In-memory event bus. Decoupled from messaging infra.

### III. Naming Conventions

| Scope | Convention | Example |
|-------|-----------|---------|
| Go files (backend) | `snake_case.go` | `farm_handler.go`, `plot_service.go` |
| Go exported types | PascalCase | `FarmHandler`, `PlotService` |
| Go JSON tags | `snake_case` | `json:"total_area_ha"` |
| TS/TSX files (admin/frontend) | PascalCase | `Farms.tsx`, `FarmDetail.tsx` |
| TS utility files | camelCase | `api.ts`, `utils.ts` |
| Mobile screens | PascalCase + `Screen.tsx` | `LoginScreen.tsx` |
| Mobile utilities | camelCase | `client.ts`, `engine.ts` |
| Branch names | `feat/*`, `fix/*`, `docs/*`, `chore/*`, `style/*` |
| Commit messages | Conventional Commits: `type: description` |

### IV. Testing Requirements

- **Go tests**: Write `*_test.go` alongside source file in the same package. Use standard `testing` package.
- **Test location**: `domain/service/*_test.go`, `event/*_test.go`.
- **Coverage**: Backend service tests required for new business logic. Run `go test ./apps/backend/...`.
- **JS/TS tests**: Not yet configured. No test framework detected. Do not add JS tests without user confirmation.

### V. Multi-Tenant Rules

- Every entity has `TenantID` UUID foreign key.
- Every handler extracts tenant ID from context via `middleware.TenantIDKey`.
- Admin routes (`/api/v1/admin/...`) are `platform_owner` only.
- Multi-tenant routes are under `/api/v1/{tenant_id}/...`.
- Sync endpoint (`POST /sync`) handles offline mobile batches (max 50).

### VI. API Conventions

- Responses use `writeJSON()` and `writeError()` helpers.
- JSON request/response format throughout.
- Swagger docs at `apps/backend/docs/`.
- Health check at `GET /health`.

### VII. Quality Gates

- `go build ./apps/backend/...` must pass.
- `npm run build` in `apps/admin/` and `apps/frontend/` must pass.
- `npm run lint` (oxlint) in `apps/admin/` and `apps/frontend/` must pass.
- No secrets in commits. `.env` and `apps/backend/seed` are gitignored.

### VIII. Mobile Offline-First

- App writes to local SQLite first, enqueues in `sync_queue` with status `pending`.
- On connectivity, sends batch (max 50) to `POST /sync`.
- Backend publishes to RabbitMQ; worker consumes and persists to PostgreSQL.
- DLQ after 3 retries.

### IX. RBAC Roles

10 roles enforced via middleware: `platform_owner`, `tenant_admin`, `proprietario`, `gerente_agricola`, `engenheiro_agronomo`, `tecnico_agricola`, `operador_campo`, `financeiro`, `consultor_externo`, `auditor`.

### X. Infrastructure

- Local dev: Docker Compose (PostgreSQL 16, Redis 7, RabbitMQ 4) via `infra/dev/docker-compose.yml`.
- Dev script: `scripts/dev.sh` manages all services.
- CI/CD: Not configured.

## Governance

This constitution supersedes all other practices. Amendments require documentation in `plans/` and user approval. All PRs/reviews must verify compliance with these rules. Complexity must be justified via Complexity Tracking in plan. Use `AGENTS.md` for runtime development guidance.

**Version**: 1.0.0 | **Ratified**: 2026-06-26 | **Last Amended**: 2026-06-26
