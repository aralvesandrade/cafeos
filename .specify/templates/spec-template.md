# Feature Specification: [FEATURE NAME]

**Feature Branch**: `[###-feature-name]`
**Created**: [DATE]
**Status**: Draft
**Input**: User description: "$ARGUMENTS"

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - [Brief Title] (Priority: P1)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently - e.g., "Can be fully tested by [specific action] and delivers [specific value]"]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]
2. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 2 - [Brief Title] (Priority: P2)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

### User Story 3 - [Brief Title] (Priority: P3)

[Describe this user journey in plain language]

**Why this priority**: [Explain the value and why it has this priority level]

**Independent Test**: [Describe how this can be tested independently]

**Acceptance Scenarios**:

1. **Given** [initial state], **When** [action], **Then** [expected outcome]

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

- What happens when [boundary condition]?
- How does system handle [error scenario]?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST [specific capability]
- **FR-002**: System MUST [specific capability]
- **FR-003**: Users MUST be able to [key interaction]
- **FR-004**: System MUST [data requirement]
- **FR-005**: System MUST [behavior]

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

### API Contract *(include if feature has API endpoints)*

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| POST | `/api/v1/{organization_id}/[resource]` | JWT + RBAC | [description] |
| GET | `/api/v1/{organization_id}/[resource]` | JWT + RBAC | [description] |
| GET | `/api/v1/{organization_id}/[resource]/{id}` | JWT + RBAC | [description] |
| PUT | `/api/v1/{organization_id}/[resource]/{id}` | JWT + RBAC | [description] |
| DELETE | `/api/v1/{organization_id}/[resource]/{id}` | JWT + RBAC | [description] |

**Admin routes** (platform_owner only):

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/api/v1/admin/[resource]` | [description] |

**Auth routes**:

| Method | Route | Description |
|--------|-------|-------------|
| POST | `/auth/login` | [description] |

### RBAC Requirements *(include if feature has permission rules)*

| Role | Access |
|------|--------|
| platform_owner | [full/admin access] |
| organization_admin | [organization config access] |
| proprietario | [read indicators, approve] |
| gerente_agricola | [manage operations] |
| engenheiro_agronomo | [technical recommendations] |
| tecnico_agricola | [data collection] |
| operador_campo | [execute operations] |
| financeiro | [financial transactions] |
| consultor_externo | [read-only authorized] |
| auditor | [compliance view] |

### Database Migrations *(include if feature adds/changes entities)*

- **New entities**: [list new GORM entity files to create in `apps/backend/internal/domain/entity/`]
- **New repository interfaces**: [list interfaces in `apps/backend/internal/domain/repository/`]
- **New repository implementations**: [list GORM impls in `apps/backend/internal/infra/db/repository/`]
- **Existing entity changes**: [list field additions/removals]

### Mobile Sync Considerations *(include if feature affects mobile)*

- **Offline queue**: Which operations need sync_queue entries?
- **Conflict resolution**: How are conflicts handled?
- **Batch size**: Max 50 per sync batch.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: [Measurable metric]
- **SC-002**: [Measurable metric]
- **SC-003**: [User satisfaction metric]
- **SC-004**: [Business metric]

## Assumptions

- [Assumption about target users]
- [Assumption about scope boundaries]
- [Assumption about data/environment]
- [Dependency on existing system/service]
