# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

**Language/Version**: Go 1.26 / TypeScript 6.0
**Primary Dependencies**:
  - Backend: GORM, golang-jwt/v5, pgx, RabbitMQ AMQP, uuid
  - Admin: React 19, Vite 8, Tailwind v4, Recharts, Radix UI, React Router v7
  - Frontend: React 19, Vite 8, Tailwind v4
  - Mobile: Expo 56, React Native 0.85, expo-sqlite, expo-secure-store, NetInfo
**Storage**: PostgreSQL (GORM ORM), Redis, local SQLite (mobile)
**Testing**: `go test ./apps/backend/...` (backend)
**Target Platform**: Web (admin + frontend), iOS + Android (mobile)
**Project Type**: Monorepo — Go API + React SPA + React Native
**Performance Goals**: [domain-specific]
**Constraints**: Mobile offline-first, multi-tenant isolation, RBAC
**Scale/Scope**: Multi-tenant SaaS, 10 RBAC roles

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[Gates determined based on constitution file]

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

**Module map — CafeOS monorepo**:

```text
apps/
├── backend/
│   ├── cmd/
│   │   ├── api/             # API entry point (:5001)
│   │   ├── worker/          # RabbitMQ worker
│   │   └── seed/            # Database seeder
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/     # HTTP handlers (writeJSON/writeError)
│   │   │   ├── middleware/  # Auth, RBAC, Organization, CORS
│   │   │   └── router.go    # Route config
│   │   ├── domain/
│   │   │   ├── entity/      # GORM entities (UUID PK, OrganizationID FK)
│   │   │   ├── repository/  # Go interfaces
│   │   │   └── service/     # Business logic
│   │   ├── event/           # In-memory event bus
│   │   └── infra/
│   │       ├── config/      # Env config
│   │       ├── db/
│   │       │   ├── postgres/     # GORM connection
│   │       │   ├── repository/   # GORM impls
│   │       │   └── migration/    # SQL migrations
│   │       └── messaging/   # RabbitMQ pub/sub
│   └── docs/                # Swagger
├── admin/
│   └── src/
│       ├── components/      # UI, layout, feature components
│       ├── pages/           # Route pages (PascalCase.tsx)
│       ├── lib/             # api.ts, auth.tsx, utils.ts
│       └── router.tsx       # React Router v7 config
├── frontend/
│   └── src/
│       ├── components/      # layout/, sections/, ui/
│       ├── lib/             # utils.ts
│       └── assets/          # Static assets
└── mobile/
    └── src/
        ├── api/             # HTTP client
        ├── db/              # SQLite schema + migrations
        ├── sync/            # Offline sync engine
        ├── hooks/           # Network status hooks
        ├── screens/         # PascalCase + Screen.tsx
        └── navigation/      # Bottom tabs
```

**Structure Decision**: Uses the existing monorepo `apps/` layout. New code goes into the appropriate module directory. No structural changes needed.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [violation] | [current need] | [why simpler not sufficient] |
