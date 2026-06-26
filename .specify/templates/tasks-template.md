---

description: "Task list template for feature implementation"

---

# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, or module like BACKEND, ADMIN, MOBILE)
- Include exact file paths in descriptions

## Path Conventions — CafeOS Monorepo

| Module | Root | Key dirs |
|--------|------|----------|
| Backend | `apps/backend/` | `cmd/api/`, `cmd/worker/`, `internal/api/handler/`, `internal/api/middleware/`, `internal/domain/entity/`, `internal/domain/service/`, `internal/domain/repository/`, `internal/infra/db/repository/`, `internal/infra/db/postgres/`, `internal/infra/messaging/`, `internal/event/` |
| Admin | `apps/admin/` | `src/pages/`, `src/components/`, `src/lib/` |
| Frontend | `apps/frontend/` | `src/components/sections/`, `src/components/layout/`, `src/components/ui/`, `src/lib/` |
| Mobile | `apps/mobile/` | `src/screens/`, `src/api/`, `src/db/`, `src/sync/`, `src/hooks/`, `src/navigation/` |

## Test Commands

- Backend: `go test ./apps/backend/...`
- Admin: `npm run build` and `npm run lint` in `apps/admin/`
- Frontend: `npm run build` and `npm run lint` in `apps/frontend/`
- All JS: `npm run build` individually per app

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize Go module or install JS dependencies
- [ ] T003 [P] Configure linting and formatting tools

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T004 [BACKEND] Create/update GORM entities in `apps/backend/internal/domain/entity/`
- [ ] T005 [BACKEND] Create repository interfaces in `apps/backend/internal/domain/repository/`
- [ ] T006 [BACKEND] Implement GORM repositories in `apps/backend/internal/infra/db/repository/`
- [ ] T007 [BACKEND] Create/update service layer in `apps/backend/internal/domain/service/`
- [ ] T008 [BACKEND] Create/update HTTP handlers in `apps/backend/internal/api/handler/`
- [ ] T009 [BACKEND] Register routes in `apps/backend/internal/api/router.go`
- [ ] T010 [ADMIN] Create/update pages in `apps/admin/src/pages/`
- [ ] T011 [ADMIN] Create/update API client calls in `apps/admin/src/lib/api.ts`
- [ ] T012 [MOBILE] Create/update screens in `apps/mobile/src/screens/`
- [ ] T013 [MOBILE] Create/update API client in `apps/mobile/src/api/client.ts`
- [ ] T014 [MOBILE] Create/update local schema in `apps/mobile/src/db/schema.ts`
- [ ] T015 [MOBILE] Create/update sync engine in `apps/mobile/src/sync/engine.ts`
- [ ] T016 [FRONTEND] Create/update sections in `apps/frontend/src/components/sections/`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - [Title] (Priority: P1) 🎯 MVP

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Backend Tasks

- [ ] T017 [P] [US1] Create [Entity] in `apps/backend/internal/domain/entity/[entity].go`
- [ ] T018 [P] [US1] Create repository interface in `apps/backend/internal/domain/repository/[entity]_repository.go`
- [ ] T019 [P] [US1] Implement GORM repository in `apps/backend/internal/infra/db/repository/[entity]_repository.go`
- [ ] T020 [US1] Implement service in `apps/backend/internal/domain/service/[entity]_service.go`
- [ ] T021 [US1] Implement handler in `apps/backend/internal/api/handler/[entity]_handler.go`
- [ ] T022 [US1] Register routes in `apps/backend/internal/api/router.go`
- [ ] T023 [US1] Add tests in `apps/backend/internal/domain/service/[entity]_service_test.go`

### Admin Tasks (if applicable)

- [ ] T024 [P] [US1] Create page at `apps/admin/src/pages/[Page].tsx`
- [ ] T025 [US1] Add API call in `apps/admin/src/lib/api.ts`
- [ ] T026 [US1] Add route in `apps/admin/src/router.tsx`

### Mobile Tasks (if applicable)

- [ ] T027 [P] [US1] Create screen at `apps/mobile/src/screens/[Screen].tsx`
- [ ] T028 [US1] Add API client method in `apps/mobile/src/api/client.ts`
- [ ] T029 [US1] Add sync handling in `apps/mobile/src/sync/engine.ts`
- [ ] T030 [US1] Add local DB schema in `apps/mobile/src/db/schema.ts`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - [Title] (Priority: P2)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Backend Tasks

- [ ] T031 [P] [US2] Create/update [Entity] in `apps/backend/internal/domain/entity/[entity].go`
- [ ] T032 [P] [US2] Create/update repository in `apps/backend/internal/infra/db/repository/[entity]_repository.go`
- [ ] T033 [US2] Implement/update service in `apps/backend/internal/domain/service/[entity]_service.go`
- [ ] T034 [US2] Implement/update handler in `apps/backend/internal/api/handler/[entity]_handler.go`
- [ ] T035 [US2] Add/update tests in `apps/backend/internal/domain/service/[entity]_service_test.go`

### Frontend Tasks (if applicable)

- [ ] T036 [P] [US2] Create/update section at `apps/frontend/src/components/sections/[Section].tsx`
- [ ] T037 [US2] Wire up in `apps/frontend/src/App.tsx`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - [Title] (Priority: P3)

**Goal**: [Brief description of what this story delivers]

**Independent Test**: [How to verify this story works on its own]

### Backend Tasks

- [ ] T038 [P] [US3] Create/update [Entity]
- [ ] T039 [P] [US3] Create/update repository interface
- [ ] T040 [US3] Implement/update GORM repository
- [ ] T041 [US3] Implement/update service
- [ ] T042 [US3] Implement/update handler

### Admin Tasks (if applicable)

- [ ] T043 [P] [US3] Create/update page and components
- [ ] T044 [US3] Add/update API calls and routing

**Checkpoint**: All user stories should now be independently functional

---

## Phase N: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] TXXX [P] Update Swagger docs in `apps/backend/docs/`
- [ ] TXXX [P] Run `go test ./apps/backend/...` and fix failures
- [ ] TXXX [P] Run `npm run build` in `apps/admin/` and `apps/frontend/`
- [ ] TXXX [P] Run `npm run lint` (oxlint) in `apps/admin/` and `apps/frontend/`
- [ ] TXXX Code cleanup and refactoring
- [ ] TXXX Update `AGENTS.md` if new routes/entities were added
- [ ] TXXX Run `scripts/dev.sh db:migrate` to apply DB changes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### Module Dependencies

- **Backend tasks**: Must be complete before Admin/Frontend tasks that consume them
- **Admin tasks**: Depend on backend API endpoints being available
- **Mobile tasks**: Depend on backend sync endpoint being available
- **Frontend tasks**: Independent of Admin; may depend on backend

### Within Each User Story

- Entities before repositories
- Repositories before services
- Services before handlers
- Handlers before frontend/admin integration
- Tests (if included) alongside service implementation

### Parallel Opportunities

- Backend, Admin, Frontend, and Mobile foundational tasks marked [P] can run in parallel
- Entity and repository interface tasks can run in parallel
- Once foundational phase completes, all user stories can start in parallel

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (backend first, then admin/mobile)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: Backend User Story 1
   - Developer B: Admin User Story 1
   - Developer C: Mobile User Story 1
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Backend test command: `go test -v ./apps/backend/internal/domain/service/...`
- Full backend test: `go test ./apps/backend/...`
