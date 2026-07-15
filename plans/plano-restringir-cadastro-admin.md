# Restringir cadastro de fazendas/usuários ao organization_admin

## Contexto

Hoje qualquer organização criada não tem usuário nenhum — quem cria os
primeiros usuários (via `POST /admin/users`, platform_owner, ou
`POST /{organization_id}/users`, quem tiver write no módulo `users`) escolhe
o papel manualmente, sem nenhuma regra. E a matriz de permissão padrão
(`defaultMatrix` em `permission_service.go:34-45,97-108`) dá **write no
módulo `farms`** pra `proprietario`, `gerente_agricola`,
`engenheiro_agronomo` e `tecnico_agricola` além de `organization_admin`, e
**write no módulo `users`** também pra `proprietario` — ou seja, hoje vários
papéis já podem cadastrar fazenda e cadastrar outro usuário, não só o
administrador da organização.

O pedido: o primeiro usuário de uma organização deve necessariamente ter o
papel `organization_admin` (system role fixa, não editável — ver
`entity/role.go:5-14`), e só `organization_admin` (+ `platform_owner`, que já
é onisciente) deve conseguir cadastrar fazendas e cadastrar/vincular novos
usuários **por padrão**. Essa restrição é o *default* semeado pra
organizações novas — continua editável depois pela tela de Permissões, como
qualquer módulo já funciona hoje (não é uma trava rígida no código).

O vínculo fazenda↔usuário (papel `proprietario`/`operador`/`financeiro` por
fazenda, campo `producers[]`) já foi implementado numa mudança anterior — o
cadastro de fazenda em si já aceita esses vínculos no mesmo request; não há
nada adicional a fazer ali além de garantir que só quem tem write em
`farms` consiga chamar esse endpoint (que já é o caso, `RequireModule` em
`router.go`).

## Mudanças — Backend

### 1. Primeiro usuário da organização é sempre `organization_admin`
- `domain/repository/user_repository.go` (+ impl GORM): novo método
  `CountByOrganization(organizationID string) (int64, error)`.
- `api/handler/user_handler.go`: nova função auxiliar
  `resolveRoleIDForCreate(organizationID, roleID, roleKey string) (string,
  error)` — se `h.repo.CountByOrganization(organizationID)` retornar `0`,
  ignora `roleID`/`roleKey` recebidos e resolve direto pra
  `entity.SystemRoleOrganizationAdmin` via `roleSvc.FindByKey`; caso
  contrário, delega pro `resolveRoleID` já existente (linhas 44-62,
  inalterado). Usado tanto em `Create` (linha ~97, admin/cross-tenant)
  quanto em `CreateForOrg` (linha ~282, org-scoped) — cobre os dois
  caminhos por onde um usuário pode ser o primeiro da organização.
- `UpdateForOrg`/`Update` continuam usando `resolveRoleID` sem essa trava —
  a regra só vale na criação do primeiro usuário, não impede trocar o papel
  de alguém depois.

### 2. `defaultMatrix` — `domain/service/permission_service.go:18-120`
Ajusta as linhas `ModuleFarms` e `ModuleUsers`:

| Role | farms (antes → depois) | users (antes → depois) |
|---|---|---|
| platform_owner | write → write | write → write |
| organization_admin | write → write | write → write |
| proprietario | write → **read** | write → **read** |
| gerente_agricola | write → **read** | none → none |
| engenheiro_agronomo | write → **read** | none → none |
| tecnico_agricola | write → **read** | none → none |
| operador_campo | read → read | none → none |
| financeiro | read → read | none → none |
| consultor_externo | read → read | none → none |
| auditor | read → read | none → none |

Só `organization_admin`/`platform_owner` ficam com write nos dois módulos;
os demais que tinham write em `farms` descem pra `read` (continuam vendo a
listagem/detalhe, só perdem criar/editar/excluir — inclui o botão
"Adicionar vínculo" da tela de fazenda, já que é a mesma tela). Continua
100% editável depois pela tela de Permissões (`PUT /{organization_id}/permissions`),
sem trava adicional de código — é só o ponto de partida semeado.

**Nota de escopo**: como já é o padrão neste projeto (ver refatoração
anterior de roles), essa mudança de `defaultMatrix` só afeta organizações
**novas** (seed roda em `OrganizationHandler.Create` →
`permSvc.SeedDefaults`) — organizações já existentes mantêm a matriz que já
foi seedada pra elas (não há backfill automático de "defaults mudaram").

## Verificação

- `go build ./... && go vet ./... && go test ./...`.
- Reset + seed do banco local, criar uma organização nova via
  `POST /admin/organizations` (platform_owner) e depois criar o primeiro
  usuário via `POST /{organization_id}/users` enviando `role_id`/`role` de
  outro papel (ex: `operador_campo`) — confirmar que o usuário criado sai
  com `role: organization_admin` mesmo assim.
- Confirmar que o segundo usuário criado na mesma organização já respeita o
  papel enviado normalmente (a trava só vale pro primeiro).
- Checar `GET /{organization_id}/permissions/me` logado como
  `proprietario`/`gerente_agricola` numa organização nova: `farms` deve vir
  `read`, `users` deve vir `read`/`none` conforme a tabela acima — e no
  admin, o botão de criar/editar fazenda e usuário deve sumir pra esses
  papéis (via `useModuleAccess`, já reativo, sem mudança de frontend
  necessária).
