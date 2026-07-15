# Refatoração: cadastro de Roles e Modules

## Contexto

O sistema de permissões (introduzido no commit `6666b8d`) tem uma peça
genuinamente configurável — a tabela `role_permissions` (matriz papel×módulo→
acesso) — mas as duas dimensões dessa matriz continuam hardcoded em código,
duplicadas em dois lugares cada:

- **Roles**: 10 constantes Go (`entity.UserRole`, `user.go:7-18`) + union type
  espelhado em TS (`apps/admin/src/lib/roles.ts`). `User.Role` é string solta,
  sem FK, sem validação em `user_handler.go:88,147` (aceita qualquer string).
- **Modules**: 8 constantes Go (`entity.Module`, `permission.go:28-48`) + union
  TS (`apps/admin/src/lib/permissions.tsx:5`), amarradas 1:1 a telas do admin
  e usadas em ~70 declarações de rota via `RequireModule(svc, ModuleX, need)`
  no `router.go`.

Decisão de escopo (confirmada com o usuário):
- **Roles viram customizáveis por organização** — admin da organização pode
  criar papéis próprios além dos padrão. `platform_owner` e
  `organization_admin` continuam papéis de sistema fixos (não editáveis/
  deletáveis, checados por string em `RequireRole`).
- **Modules ganham cadastro/tela**, mas a lista de módulos usada nas rotas
  continua fixa em código — cada módulo já corresponde a uma tela+rotas
  reais; criar um módulo novo por tela ainda exige deploy. O cadastro serve
  para visualizar/nomear/ordenar módulos e ser a fonte única (banco) em vez
  de duas listas hardcoded divergentes.

Resultado: elimina duplicação Go/TS, dá times de organização o poder de criar
papéis sob medida (ex: "colhedor_chefe") sem precisar de deploy no backend, e
mantém a superfície de módulos estável (sem o risco de tornar `router.go`
dinâmico).

## Mudanças — Backend

### 1. Nova entidade `Role` (`domain/entity/role.go`)
```go
type Role struct {
    ID             string  // uuid
    OrganizationID *string // NULL = papel de sistema (global), preenchido = papel custom da org
    Key            string  // slug estável, ex: "colhedor_chefe" — usado em RolePermission.Role e User.Role
    Name           string  // label exibido
    IsSystem       bool    // true para platform_owner e organization_admin — não editável/deletável
    CreatedAt, UpdatedAt time.Time
}
```
Unique index em `(organization_id, key)` (tratando NULL como um valor lógico
via índice parcial ou constante sentinela — seguir o padrão já usado em
`idx_org_role_module`).

Seed: as 10 roles atuais viram linhas com `organization_id = NULL,
is_system = true` para `platform_owner`/`organization_admin`, e
`is_system = false` para as outras 8 (organização pode editar/apagar essas
oito se não estiverem em uso — são só o "kit inicial", não mais papéis de
sistema).

### 2. Nova entidade `Module` (`domain/entity/module.go`)
Substitui as constantes por linhas de tabela: `Key`, `Name`, `Order`. Seed
único, global (sem `organization_id` — módulo é característica da aplicação,
não da organização). CRUD limitado a leitura/rename/reorder pelo admin —
**não** permite criar módulo novo arbitrário (rotas continuam amarradas às 8
chaves existentes), então o handler expõe `List` e `Update` (nome/ordem), não
`Create`/`Delete`.

`entity.Module` deixa de ser `type Module string` com constantes — vira
`type ModuleKey string` (as mesmas 8 strings, ainda usadas como chave
type-safe em `router.go`/`RequireModule`), e a tabela `modules` guarda os
metadados (nome, ordem) para essas mesmas chaves. Ou seja: as chaves
continuam fixas em código (necessário, pois amarram rotas), mas texto/ordem
exibidos vêm do banco.

### 3. `RolePermission` (`entity/permission.go:52-60`)
Troca `Role UserRole` → `RoleID string` (FK para `roles.id`), mantém
`Module ModuleKey`. Ajustar unique index para `(organization_id, role_id,
module)`.

### 4. `User.Role` (`entity/user.go:26`)
Troca `Role UserRole` (string solta) → `RoleID string` (FK para `roles.id`).
JWT: avaliar se carrega `role_id` + `role_key` (key necessária para
`RequireRole`/`RequireModule` sem round-trip ao banco a cada request) —
recomendado manter `role_key` no claim, já que é isso que os middlewares
comparam hoje.

### 5. Repository/Service novos
- `domain/repository/role_repository.go` + impl GORM — seguir exatamente o
  padrão de `permission_repository.go` (`Upsert`, `ListByOrganization`,
  incluindo papéis de sistema globais na listagem).
- `domain/service/role_service.go` — CRUD com regras: não permite editar/
  apagar `IsSystem`, não permite apagar role em uso por algum `User` ou
  `RolePermission`, valida `Key` único por organização.
- `domain/service/module_service.go` — `List`/`UpdateMeta` (nome/ordem), sem
  create/delete.
- `permission_service.go`: `defaultMatrix`/`allRoles`/`validRole` deixam de
  ser constantes — `validRole` passa a consultar `role_service`/repo;
  `defaultMatrix` (seed) permanece como está mas chaveado por `Key` de role
  em vez de `UserRole`, aplicado só às roles de kit inicial.

### 6. Handlers novos
- `role_handler.go`: `GET/POST/PUT/DELETE /{organization_id}/roles` —
  seguir estrutura de `permission_handler.go` (thin, chama service).
- `module_handler.go`: `GET /{organization_id}/modules`,
  `PUT /{organization_id}/modules/{key}` (nome/ordem).
- `user_handler.go:88,147`: validar `req.Role`/`req.RoleID` contra
  `role_service` em vez de aceitar qualquer string (fecha o gap atual).

### 7. Rotas (`router.go`)
Novo bloco de rotas para roles/modules, protegido por
`RequireModule(permSvc, ModulePermissions, AccessWrite)` (reaproveita módulo
"permissions" existente — gestão de papéis é parte do mesmo domínio de
governança de acesso).

### 8. Migration
`007_roles_modules.sql` — documentação do schema (GORM AutoMigrate segue
sendo a fonte real, mesmo padrão do `005_role_permissions.sql`):
- `roles` (id, organization_id nullable, key, name, is_system, timestamps)
- `modules` (key PK, name, order, timestamps) — seed dos 8 módulos atuais
- `role_permissions.role_id` substitui `role_permissions.role`
- `users.role_id` substitui `users.role`
- Migração de dados: popular `roles` com as 10 atuais, depois back-fill
  `role_id` em `role_permissions`/`users` a partir da string antiga, por
  último dropar colunas string.

### 9. `RequireRole` (`middleware/rbac.go`)
Sem mudança de assinatura — continua comparando string (`role_key` do JWT)
contra lista fixa passada na rota. Usado só para `platform_owner` (rotas
`/admin/...`), que continua papel de sistema fixo.

## Mudanças — Frontend (admin)

- `lib/roles.ts`: remove union type + `ROLE_LABELS`/`ALL_ROLES` hardcoded;
  passam a vir de `GET /{org}/roles` (novo hook `useRoles()`, mesmo padrão de
  `usePermissions()` em `lib/permissions.tsx`).
- `lib/permissions.tsx:5,8-19`: `MODULE_LABELS`/`ALL_MODULES` passam a vir de
  `GET /{org}/modules`; `ROUTE_MODULE` (mapeamento rota→módulo) continua
  hardcoded no front (é puramente de navegação, não precisa ir ao banco).
- Nova tela `pages/Roles.tsx` (CRUD de papéis customizados da organização),
  seguindo o padrão visual/estrutural de `pages/Permissions.tsx` e
  `pages/TeamUsers.tsx`.
- `pages/Permissions.tsx`: passa a montar a matriz a partir dos roles/modules
  vindos da API em vez das constantes locais.
- Sidebar/rota nova para "Papéis" (`/roles`), reaproveitando o guard de
  módulo `permissions` já usado em `/permissions`.

## Ordem de execução sugerida

1. Entities + migration (roles, modules) — sem quebrar nada ainda (colunas
   antigas continuam existindo em paralelo).
2. Repository + service (role, module) + testes unitários (seguir
   `permission_service_test.go` como referência de estilo).
3. Back-fill de dados (script/migration) — popular `roles`/`role_id`.
4. Handlers + rotas novas.
5. Trocar `RolePermission.Role`/`User.Role` para `RoleID`, remover coluna
   string antiga, atualizar `permission_service.go`, `user_handler.go`,
   `router.go` (SeedDefaults, JWT claims).
6. Frontend: hooks `useRoles`/módulos dinâmicos, tela `Roles.tsx`, ajustar
   `Permissions.tsx`.
7. `swag init` (obrigatório após mudança de handlers/entities — regra do
   projeto).

## Verificação

- `go test ./...` (cobrir `role_service_test.go`, `module_service_test.go`
  novos + `permission_service_test.go` atualizado).
- `go build ./...` e `go vet ./...`.
- Rodar `./scripts/dev.sh db:reset && db:migrate && db:seed` e validar login
  com os 7 usuários seed (tabela de credenciais no CLAUDE.md) — confirmar que
  cada um ainda enxerga os módulos/telas esperados no admin.
- Fluxo manual no admin: criar papel customizado numa organização, atribuir a
  um usuário novo, configurar acesso na tela de Permissões, logar como esse
  usuário e confirmar que o menu lateral reflete o acesso configurado.
- Confirmar que `platform_owner` continua acessando `/admin/...` mesmo após a
  migração (papel de sistema não deve quebrar).
