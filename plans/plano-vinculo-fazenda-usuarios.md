# Vínculo Fazenda ↔ Usuários (múltiplos papéis por fazenda)

## Contexto

Hoje `Organization` e `User` já têm hierarquia clara (usuário pertence a uma
organização, com um papel global via `RoleID`). Mas `Farm` não tem vínculo
real com usuários específicos — existe uma entity `Producer` (1:1 com Farm,
reforçado por `uniqueIndex` em `farm_id`) que guarda dados cadastrais soltos
(CPF, RG, nome...) e tem um campo `UserID *string` opcional que **nunca é
exposto na UI nem no payload de request** (`producerRequest` em
`farm_handler.go:32-43` não tem `user_id`). Na prática, todo `Producer`
criado pela aplicação tem `UserID = nil` — só o seed popula esse campo
manualmente.

Isso já quebra silenciosamente a única feature que depende dele: usuários
com papel global `proprietario` deveriam só enxergar as fazendas que
possuem (via `FarmService.OwnedFarmIDs`/`ListByOwner`/`IsOwner`, usadas em
9+ handlers — harvest, operation, plot, financial, stock, fleet, dashboard),
mas como `UserID` nunca é setado no fluxo real, essa restrição só funciona
para os usuários semeados manualmente.

O pedido: uma fazenda deve poder ter **múltiplos usuários vinculados**, cada
um com um papel específico *daquela fazenda* (ex: Fazenda XPTO tem João como
proprietário, Ana como operador de campo, Carlos do financeiro) — um mesmo
usuário pode estar vinculado a várias fazendas. Decisões já validadas com o
usuário:

- `Producer` deixa de ser 1:1 e vira o vínculo N:N Farm↔User (não cria uma
  entity nova separada).
- O "papel na fazenda" reaproveita o catálogo global de `Role` (mesma tabela
  `roles` da refatoração anterior) via `RoleID` — é só informativo aqui, não
  altera o `RoleID` global do usuário nem a permissão de sistema dele.
- A restrição de visão continua valendo **só** para o papel global
  `proprietario` (comportamento atual preservado) — só troca a fonte de
  dados de "1 producer" para "N vínculos", sem estender a restrição a
  operador/financeiro.

## Mudanças — Backend

### 1. `entity/producer.go`
- `FarmID`: tira `uniqueIndex`, vira `index` normal (permite N vínculos por
  fazenda).
- `UserID`: de `*string` opcional vira `string` obrigatório (`not null`) —
  todo vínculo agora exige um usuário real da organização.
- Novo campo `RoleID string` (`not null`, FK pra `Role`) — o papel do
  usuário *nesta fazenda*.
- Unique index composto em `(farm_id, user_id)` — um usuário não pode ter
  dois vínculos na mesma fazenda.
- Campos legais (CPF/RG/órgão emissor/etc.) continuam existindo mas passam
  a ser conceitualmente só relevantes quando o vínculo é do papel
  `proprietario` (produtor legal p/ CAR/INCRA) — sem enforcement rígido no
  banco, só orientação de uso na UI.
- Adiciona `Role Role` (`gorm:"foreignKey:RoleID"`) na entity, junto da já
  existente `User *User` (que vira `User User`, não mais ponteiro).

### 2. `domain/repository/producer_repository.go` + impl GORM
Troca de "1 por fazenda" pra "N por fazenda":
- `GetByFarmID` → `ListByFarmID(farmID string) ([]*entity.Producer, error)`
- Novo `ListByUserID(userID string) ([]*entity.Producer, error)` — usado
  pra resolver direto quais fazendas um usuário possui, sem precisar
  carregar todas as fazendas da organização e filtrar em memória (melhora
  o que `FarmService.ListByOwner` faz hoje em `farm_service.go:79-91`).
- `DeleteByFarmID` continua igual (já deleta todas as linhas daquele
  farm_id, útil tanto pra "remover 1 vínculo" via delete+recreate quanto
  pra apagar a fazenda inteira).

### 3. `domain/service/farm_service.go`
- `Create(farm, producers []*entity.Producer)` — lista em vez de ponteiro
  único; cada producer precisa de `UserID`+`RoleID` válidos.
- `SetProducers(farmID string, producers []*entity.Producer) error` —
  substitui `UpsertProducer`: apaga todos os vínculos da fazenda e recria
  a lista enviada (padrão "PUT substitui o estado", igual ao que
  `PermissionHandler.Update` já faz pra `role_permissions`).
- `ListByOwner`/`OwnedFarmIDs`: passam a usar `producerRepo.ListByUserID`
  diretamente (filtrando por `organization_id`), em vez de
  `ListByOrganization` + filtro em memória.
- `IsOwner(farmID, userID)`: verifica se existe alguma linha
  `Producer{FarmID: farmID, UserID: userID}` (usa `ListByFarmID` + busca,
  ou novo `ExistsByFarmAndUser` no repo — mais direto).
- `Delete(id)`: continua chamando `DeleteByFarmID` antes de apagar a farm.

### 4. `api/handler/farm_handler.go`
- `producerRequest` ganha `UserID string json:"user_id"` (obrigatório) e
  `RoleID string json:"role_id"` (obrigatório).
- `createFarmRequest.Producer *producerRequest` → `Producers
  []producerRequest` (plural).
- `Update`: `input.Producers []producerRequest` chama
  `svc.SetProducers(existing.ID, producers)`.
- `canAccessFarm`/`List` (linhas 225-260): lógica não muda — continuam
  checando só o papel global `proprietario` do usuário logado; só a
  implementação por trás (`IsOwner`/`ListByOwner`) muda de fonte.

### 5. Migration de documentação
Nova `internal/infra/db/migration/008_farm_user_links.sql` documentando:
tira o `UNIQUE` de `producers.farm_id`, adiciona `role_id UUID NOT NULL`
com FK pra `roles(id)`, adiciona `UNIQUE(farm_id, user_id)`, `user_id`
passa a `NOT NULL` com FK pra `users(id)` (já implícito, só reforça).
Segue o padrão dos arquivos anteriores (é documentação — schema real vem
do `AutoMigrate`, ver `internal/infra/db/postgres/connection.go`).

### 6. `cmd/seed/main.go`
Ajustar o seed (hoje popula `Producer` com `UserID` direto — ver
`cmd/seed/main.go:~245`) pra also setar `RoleID` (resolver via `roleRepo`,
igual já é feito pra `User.RoleID`), e trocar pra lista de producers por
fazenda se fizer sentido demonstrar múltiplos vínculos num exemplo (ex:
Fazenda Recanto Verde com João como proprietário + Ana como operador).

## Mudanças — Frontend (admin)

### `components/farms/FarmForm.tsx`
- `FarmData.producer: ProducerData` → `FarmData.producers: ProducerLink[]`,
  onde `ProducerLink` tem `user_id`, `role_id` + os campos legais atuais
  (cpf, name, rg, etc.), todos opcionais exceto `user_id`/`role_id`.
- Seção "Dados do Produtor" vira "Vínculos da Fazenda": lista de linhas,
  cada uma com `Select` de usuário (buscar via `GET /users`, mesmo padrão
  de `TeamUsers.tsx`) e `Select` de papel (via `useRoles()`, já existe em
  `lib/roles.tsx`), botão "Adicionar vínculo" e remover linha. Campos
  legais (CPF/RG/etc.) só aparecem/expandem quando o papel selecionado for
  `proprietario` (comparar `role.key`).

### `pages/FarmEdit.tsx`
- `farmToFormData`/`formDataToPayload`: adaptar de `producer` singular pra
  `producers` (array), mapeando `farm.producers` (retornado pelo backend)
  pra `ProducerLink[]` e vice-versa.

## Ordem de execução sugerida

1. Entity + migration (Producer: tira uniqueIndex, User/RoleID obrigatórios,
   novo unique composto).
2. Repository (`ListByFarmID`, `ListByUserID`, `ExistsByFarmAndUser`).
3. Service (`Create` com lista, `SetProducers`, `ListByOwner`/`OwnedFarmIDs`/
   `IsOwner` usando `ListByUserID`).
4. Handler (`producerRequest` com `user_id`/`role_id`, `Producers` plural).
5. Seed atualizado.
6. `swag init` (regra do projeto — mudou entity/handler).
7. Frontend: `FarmForm.tsx` (seção de vínculos com select de usuário+papel),
   `FarmEdit.tsx` (mapear array).

## Verificação

- `go build ./... && go vet ./... && go test ./...`.
- Reset + seed do banco local (`docker compose exec postgres psql ... DROP
  SCHEMA` + `go run ./cmd/seed`), confirmar no `\d producers` que o unique
  index mudou de `(farm_id)` pra `(farm_id, user_id)` e que `role_id`/FK
  existem.
- Fluxo manual no admin: criar fazenda com 2 vínculos (proprietário +
  operador, usuários diferentes), confirmar que aparecem na listagem;
  logar como o usuário vinculado como `proprietario` e confirmar que só
  essa fazenda aparece; logar como o usuário vinculado como `operador` (se
  o papel global dele não for `proprietario`) e confirmar que continua
  vendo todas as fazendas da organização (restrição não se estendeu).
