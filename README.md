# CafeOS

**A plataforma especialista em cafeicultura**

Plataforma SaaS multi-tenant para gestão operacional, produtiva, financeira e analítica de propriedades cafeeiras.

## Stack

| Camada    | Tecnologia                                              |
| --------- | ------------------------------------------------------- |
| Backend   | Go (REST API, workers, engine de regras)                |
| Frontend  | React + Vite + Tailwind CSS v4 (landing page)           |
| Admin     | React + Vite + Tailwind CSS v4 + Recharts               |
| Mobile    | React Native (Expo) + SQLite         |
| Banco     | PostgreSQL + Redis (GORM ORM)                            |
| Mensageria| RabbitMQ                                                 |
| Infra     | Docker, Kubernetes, ArgoCD, Grafana, Prometheus         |

## Monorepo Structure

```
cafeos/
├── apps/
│   ├── backend/          # Go API (MVP atual)
│   ├── frontend/         # React + Vite (landing page)
│   ├── admin/            # React + Vite (painel administrativo)
│   └── mobile/           # React Native (Expo) + SQLite
├── packages/
│   └── shared/           # Tipos e utilitários compartilhados
├── .specify/             # Spec-kit (constitution, templates, workflows)
├── docs/                 # Documentação técnica
├── scripts/              # Scripts de desenvolvimento
├── plans/                # Planos de arquitetura e features
├── PRD_CONSOLIDADO.md    # Product Requirements Document
├── PRD1.md               # PRD original (SPEC-001)
└── PRD2.md               # PRD original (SPEC-002)
```

## Backend (MVP) — `apps/backend/`

```
internal/
├── domain/
│   ├── entity/           # Entidades de domínio (Farm, Plot, Operation, Harvest, etc.) com tags GORM
│   ├── repository/       # Interfaces de repositório + Transactor
│   └── service/          # Serviços de domínio + Rule Engine
├── api/
│   ├── handler/          # Handlers HTTP (REST)
│   ├── middleware/       # Auth, RBAC, Tenant, CORS
│   └── router.go         # Configuração de rotas
├── event/                # Event bus in-memory (preparado para RabbitMQ/Kafka)
└── infra/
    ├── config/           # Configuração via environment variables
    ├── db/
    │   ├── postgres/     # Conexão GORM + AutoMigrate + Transactor
    │   ├── repository/   # Implementações dos repositórios (GORM + WithTx)
    │   └── migration/    # Migrations SQL (legado)
    └── messaging/        # Integração RabbitMQ (futuro)
```

### Entidades Principais

| Entidade               | Descrição                                      |
| ---------------------- | ---------------------------------------------- |
| `Tenant`               | Cliente multi-tenant (White Label)             |
| `User`                 | Usuário com RBAC (10 perfis)                   |
| `Farm`                 | Propriedade rural (fazenda)                    |
| `Plot` (Talhão)        | Talhão vinculado a uma fazenda                 |
| `Operation`            | Operação agrícola (adubação, pulverização...)  |
| `Harvest` (Safra)      | Safra agrícola por ano                         |
| `HarvestProduction`    | Produção por talhão por safra                  |
| `AgriculturalProduct`  | Insumos agrícolas cadastrados                  |
| `Indicator`            | Indicadores calculados (sacas/ha, custo/saca)  |

### Perfis RBAC

- `platform_owner` — Admin global (cria tenants, planos)
- `tenant_admin` — Admin do tenant (usuários, config)
- `proprietario` — Proprietário rural (indicadores)
- `gerente_agricola` — Gerente agrícola (operações)
- `engenheiro_agronomo` — Recomendações técnicas
- `tecnico_agricola` — Coleta de dados
- `operador_campo` — Executa operações
- `financeiro` — Custos e fluxo de caixa
- `consultor_externo` — Leitura autorizada
- `auditor` — Compliance e rastreabilidade

## Frontend — `apps/frontend/`

Landing page responsiva com paleta de cores cafeicultura:

```
src/
├── components/
│   ├── layout/       # Header, Footer
│   ├── sections/     # Hero, About, Features, CoffeeCycle, Indicators, Plans, TechStack, Roadmap, CTA
│   └── ui/           # Button, Badge (shadcn/ui)
├── lib/utils.ts      # cn() utility
├── App.tsx           # Página principal
├── main.tsx          # Entry point
└── index.css         # Tailwind v4 @theme
```

```bash
cd apps/frontend
npm run dev      # Desenvolvimento
npm run build    # Build produção
```

## Admin Panel — `apps/admin/`

Painel administrativo com autenticação JWT, dashboard com gráficos e CRUD completo:

```
src/
├── components/
│   ├── layout/       # Sidebar, Header, AdminLayout, AuthLayout
│   ├── ui/           # Button, Badge, Card, Table, Dialog, Input, Select
│   ├── dashboard/    # StatsCards, ProductionChart, CostChart, RecentOperations
│   └── farms/        # FarmList, FarmForm
├── lib/
│   ├── utils.ts      # cn() utility
│   ├── api.ts        # Cliente HTTP com JWT
│   └── auth.tsx      # AuthContext + hook
├── pages/            # Login, Dashboard, Farms, Plots, Operations, Harvests, Tenants, Users, NotFound
├── router.tsx        # React Router DOM v7 (nested layouts, lazy routes)
├── App.tsx
├── main.tsx
└── index.css         # Tailwind v4 @theme
```

```bash
cd apps/admin
npm run dev      # Desenvolvimento (http://localhost:5174)
npm run build    # Build produção
```

### Swagger

Documentação interativa disponível em `http://localhost:5001/swagger/index.html` (com backend rodando).

### API REST

Rotas multi-tenant sob `/api/v1/{tenant_id}`:

| Método | Rota                           | Descrição                    |
| ------ | ------------------------------ | ---------------------------- |
| GET    | `/health`                      | Health check                 |
| POST   | `/farms`                       | Criar fazenda                |
| GET    | `/farms`                       | Listar fazendas              |
| GET    | `/farms/{id}`                  | Detalhe da fazenda           |
| PUT    | `/farms/{id}`                  | Atualizar fazenda            |
| DELETE | `/farms/{id}`                  | Remover fazenda              |
| POST   | `/plots`                       | Criar talhão                 |
| GET    | `/plots`                       | Listar talhões               |
| GET    | `/farms/{farm_id}/plots`       | Listar talhões por fazenda   |
| GET    | `/plots/{id}`                  | Detalhe do talhão            |
| PUT    | `/plots/{id}`                  | Atualizar talhão             |
| DELETE | `/plots/{id}`                  | Remover talhão               |
| POST   | `/operations`                  | Registrar operação           |
| GET    | `/operations`                  | Listar operações             |
| GET    | `/operations/recent`           | Operações recentes           |
| GET    | `/operations/{id}`             | Detalhe da operação          |
| DELETE | `/operations/{id}`             | Remover operação             |
| GET    | `/plots/{plot_id}/operations`  | Operações por talhão         |
| POST   | `/harvests`                    | Criar safra                  |
| GET    | `/harvests`                    | Listar safras                |
| GET    | `/harvests/{id}`               | Detalhe da safra             |
| PUT    | `/harvests/{id}/finalize`      | Finalizar safra              |
| POST   | `/harvests/{id}/production`    | Registrar produção           |
| GET    | `/harvests/{id}/production`    | Produção da safra            |
| GET    | `/dashboard`                   | Dashboard consolidado        |
| POST   | `/financial`                   | Criar transação financeira    |
| GET    | `/financial`                   | Listar transações             |
| GET    | `/financial/{id}`              | Detalhe da transação          |
| PUT    | `/financial/{id}`              | Atualizar transação           |
| DELETE | `/financial/{id}`              | Remover transação             |
| GET    | `/agricultural-products`       | Listar produtos agrícolas     |
| POST   | `/stock/items`                 | Criar item de estoque         |
| GET    | `/stock/items`                 | Listar itens                  |
| PUT    | `/stock/items/{id}`            | Atualizar item                |
| DELETE | `/stock/items/{id}`            | Remover item                  |
| POST   | `/stock/movements`             | Registrar movimentação        |
| GET    | `/stock/movements`             | Listar movimentações          |
| POST   | `/fleet/vehicles`              | Criar veículo                 |
| GET    | `/fleet/vehicles`              | Listar veículos               |
| PUT    | `/fleet/vehicles/{id}`         | Atualizar veículo             |
| DELETE | `/fleet/vehicles/{id}`         | Remover veículo               |
| POST   | `/fleet/maintenance`           | Registrar manutenção          |
| GET    | `/fleet/maintenance`           | Listar manutenções            |
| DELETE | `/fleet/maintenance/{id}`      | Remover manutenção            |
| POST   | `/labor/teams`                 | Criar equipe                  |
| GET    | `/labor/teams`                 | Listar equipes                |
| PUT    | `/labor/teams/{id}`            | Atualizar equipe              |
| DELETE | `/labor/teams/{id}`            | Remover equipe                 |
| POST   | `/labor/workers`               | Criar trabalhador             |
| GET    | `/labor/workers`               | Listar trabalhadores          |
| PUT    | `/labor/workers/{id}`          | Atualizar trabalhador         |
| DELETE | `/labor/workers/{id}`          | Remover trabalhador           |
| POST   | `/labor/shifts`                | Registrar apontamento         |
| GET    | `/labor/shifts`                | Listar apontamentos           |
| DELETE | `/labor/shifts/{id}`           | Remover apontamento           |
| POST   | `/sync`                        | Sincronizar lote offline      |

Rotas admin (`platform_owner` apenas, prefixo `/api/v1/admin`):

| Método | Rota                          | Descrição              |
| ------ | ----------------------------- | ---------------------- |
| GET    | `/api/v1/admin/tenants`       | Listar tenants         |
| POST   | `/api/v1/admin/tenants`       | Criar tenant           |
| GET    | `/api/v1/admin/tenants/{id}`  | Detalhe do tenant      |
| PUT    | `/api/v1/admin/tenants/{id}`  | Atualizar tenant       |
| DELETE | `/api/v1/admin/tenants/{id}`  | Remover tenant         |
| GET    | `/api/v1/admin/users`         | Listar usuários        |
| POST   | `/api/v1/admin/users`         | Criar usuário          |
| PUT    | `/api/v1/admin/users/{id}`    | Atualizar usuário      |
| DELETE | `/api/v1/admin/users/{id}`    | Remover usuário        |

### Engine de Regras

Motor de regras configurável para alertas automáticos:

- **Baixa produtividade**: alerta se < 25 sacas/hectare
- **Custo elevado**: alerta se > R$ 400/saca
- Regras customizáveis via `RuleEngine.AddRule()`

### Eventos

Sistema orientado a eventos (in-memory, preparado para fila):

| Evento                    | Disparo                           |
| ------------------------- | --------------------------------- |
| `OperationRegistered`     | Ao registrar operação agrícola    |
| `HarvestFinalized`        | Ao finalizar safra                |
| `IndicatorUpdated`        | Ao recalcular indicadores         |
| `AlertGenerated`          | Quando regra é acionada           |

### Transações

O HarvestService utiliza o `Transactor` para garantir atomicidade na finalização de safras (atualização + cálculo de indicadores). Demais services operam com repositories individuais. Uso:

```go
transactor.RunInTx(func(repos repository.TransactionProvider) error {
    repos.Harvest().Update(harvest)
    repos.Indicator().Create(indicator)
    return nil
})
```

## Mobile — `apps/mobile/`

App React Native (Expo) offline-first para operações de campo:

```
src/
├── api/            # Cliente HTTP (JWT)
├── db/             # SQLite local + migrações
├── sync/           # Engine de sincronização (lote 50)
├── hooks/          # Network status, operações offline
├── screens/        # Login, Operations, PendingSync
└── navigation/     # Bottom tabs
```

```bash
cd apps/mobile
npm install
npx expo start
```

## Worker — `cmd/worker`

Processo separado que consome filas RabbitMQ e persiste no PostgreSQL:

```bash
# Requer RabbitMQ rodando
cd apps/backend && go run ./cmd/worker/main.go
```

## Desenvolvimento Local

### Pré-requisitos

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### Serviços

| Serviço    | Porta | Acesso                     |
|-----------|-------|----------------------------|
| PostgreSQL | 5432  | `cafeos:cafeos@localhost`  |
| Redis      | 6379  | `redis://localhost`        |
| RabbitMQ   | 5672  | `amqp://cafeos:cafeos@localhost` |
| RabbitMQ UI| 15672 | `http://localhost:15672`   |

### Comandos

```bash
# Subir infra + API (background)
./scripts/dev.sh up

# Parar todos containers
./scripts/dev.sh down

# Serviços individuais (cada um em terminal separado):
./scripts/dev.sh api       # API Go na :5001
./scripts/dev.sh worker    # Worker RabbitMQ
./scripts/dev.sh admin     # Admin panel na :5174
./scripts/dev.sh mobile    # App mobile na :8081

# Utilitários
./scripts/dev.sh db:migrate  # Rodar migrations
./scripts/dev.sh db:reset    # Resetar banco
./scripts/dev.sh db:seed     # Seed dados iniciais
./scripts/dev.sh test        # Testes backend
```

### Variáveis de Ambiente

| Variável       | Default                                      |
| -------------- | -------------------------------------------- |
| `SERVER_PORT`  | `5001`                                       |
| `DATABASE_URL` | `postgres://cafeos:cafeos@localhost:5432/cafeos?sslmode=disable` |
| `REDIS_URL`    | `redis://localhost:6379`                     |
| `JWT_SECRET`   | `dev-secret-change-in-production`            |
| `RABBITMQ_URL` | `amqp://cafeos:cafeos@localhost:5672/`       |

## Roadmap

### Fase 1 — MVP ✅ (atual)
- [x] Gestão de fazendas e talhões
- [x] Operações agrícolas
- [x] Gestão de safras
- [x] Custos agrícolas
- [x] Dashboard inicial
- [x] Engine de regras
- [x] API REST multi-tenant
- [x] Event system
- [x] RBAC com 10 perfis + autorização por rota
- [x] Gestão de tenants e usuários (platform_owner)
- [x] Login com acesso rápido por perfil (dev)

### Fase 2 ✅ (implementado)
- [x] Financeiro (contas a pagar/receber, categorias)
- [x] Estoque (insumos, validade, movimentações)
- [x] Frota (veículos, manutenções preventivas/corretivas)
- [x] Mão de Obra (equipes, trabalhadores, apontamento de horas)

### Fase 3 ✅ (MVP)
- [x] Mobile offline (React Native + SQLite + sync engine)
- [x] RabbitMQ para fila de sincronização
- [x] Worker para processar registros offline (cmd/worker)
- [ ] Cooperativas e consultorias (pendente)
- [ ] Integrações externas (pendente)

### Fase 4
- [ ] IoT (sensores, estações meteorológicas)
- [ ] IA (previsão de safra, recomendação, detecção de doenças)
- [ ] Analytics avançado

## Desenvolvimento com IA

O arquivo `AGENTS.md` na raiz do projeto contém o resumo técnico para agentes de IA (stack, rotas, convenções, credenciais de seed).

## Spec-kit

Projeto configurado com [spec-kit](https://github.com/anomalyco/spec-kit) (`extensão brownfield`) para desenvolvimento assistido por IA. O fluxo segue:

| Comando | Descrição |
|---------|-----------|
| `/speckit.brownfield.scan` | Escanear projeto existente e detectar stack + convenções |
| `/speckit.brownfield.bootstrap` | Gerar constitution + templates personalizados |
| `/speckit.brownfield.validate` | Validar configuração contra o projeto |
| `/speckit.brownfield.migrate` | Reverter engenharia de specs para features existentes |
| `/speckit.specify` | Criar especificação de nova feature |
| `/speckit.clarify` | Esclarecer requisitos da spec |
| `/speckit.plan` | Gerar plano de implementação |
| `/speckit.tasks` | Gerar lista de tarefas |
| `/speckit.checklist` | Gerar checklist de verificação |
| `/speckit.analyze` | Analisar impacto de mudanças |

Configuração em `.specify/`:
- `memory/constitution.md` — Regras do projeto (stack, módulos, convenções, RBAC)
- `templates/` — Templates de spec, plan, tasks, checklist
- `extensions/brownfield/` — Extensão de brownfield migration
- `extensions/git/` — Extensão de automação git

## Licença

Proprietária — todos os direitos reservados.
