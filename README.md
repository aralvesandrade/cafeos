# CafeOS

**A plataforma especialista em cafeicultura**

Plataforma SaaS multi-tenant para gestão operacional, produtiva, financeira e analítica de propriedades cafeeiras.

## Stack

| Camada    | Tecnologia                                              |
| --------- | ------------------------------------------------------- |
| Backend   | Go (REST API, workers, engine de regras)                |
| Frontend  | React + Vite + Tailwind CSS v4                          |
| Mobile    | React Native _(futuro)_                                 |
| Banco     | PostgreSQL + Redis                                      |
| Mensageria| RabbitMQ / Kafka _(futuro)_                             |
| Infra     | Docker, Kubernetes, ArgoCD, Grafana, Prometheus         |

## Monorepo Structure

```
cafeos/
├── apps/
│   ├── backend/          # Go API (MVP atual)
│   ├── frontend/         # React + Vite (landing page)
│   └── mobile/           # React Native (futuro)
├── packages/
│   └── shared/           # Tipos e utilitários compartilhados
├── infra/
│   └── dev/              # Docker Compose para desenvolvimento
├── docs/                 # Documentação técnica
├── scripts/              # Scripts de desenvolvimento
├── PRD_CONSOLIDADO.md    # Product Requirements Document
├── PRD1.md               # PRD original (SPEC-001)
└── PRD2.md               # PRD original (SPEC-002)
```

## Backend (MVP) — `apps/backend/`

```
internal/
├── domain/
│   ├── entity/           # Entidades de domínio (Farm, Plot, Operation, Harvest, etc.)
│   ├── repository/       # Interfaces de repositório
│   └── service/          # Serviços de domínio + Rule Engine
├── api/
│   ├── handler/          # Handlers HTTP (REST)
│   ├── middleware/       # Auth, RBAC, Tenant, CORS
│   └── router.go         # Configuração de rotas
├── event/                # Event bus in-memory (preparado para RabbitMQ/Kafka)
└── infra/
    ├── config/           # Configuração via environment variables
    ├── db/
    │   ├── postgres/     # Conexão PostgreSQL
    │   ├── repository/   # Implementações dos repositórios
    │   └── migration/    # Migrations SQL
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

### Swagger

Documentação interativa disponível em `http://localhost:8080/swagger/index.html` (com backend rodando).

### API REST

Todas as rotas sob `/api/v1/{tenant_id}`:

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
| GET    | `/plots/{plot_id}/operations`  | Operações por talhão         |
| POST   | `/harvests`                    | Criar safra                  |
| GET    | `/harvests`                    | Listar safras                |
| GET    | `/harvests/{id}`               | Detalhe da safra             |
| PUT    | `/harvests/{id}/finalize`      | Finalizar safra              |
| POST   | `/harvests/{id}/production`    | Registrar produção           |
| GET    | `/harvests/{id}/production`    | Produção da safra            |
| GET    | `/dashboard`                   | Dashboard consolidado        |

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

## Desenvolvimento Local

### Pré-requisitos

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose

### Comandos

```bash
# Subir ambiente + API
./scripts/dev.sh up

# Apenas infraestrutura (PostgreSQL + Redis)
./scripts/dev.sh up        # depois rode o backend separadamente

# Parar ambiente
./scripts/dev.sh down

# Rodar migrations
./scripts/dev.sh db:migrate

# Resetar banco
./scripts/dev.sh db:reset

# Rodar testes backend
./scripts/dev.sh test
# ou diretamente:
cd apps/backend && go test ./... -v

# Frontend (outro terminal)
cd apps/frontend && npm run dev
```

### Variáveis de Ambiente

| Variável       | Default                                      |
| -------------- | -------------------------------------------- |
| `SERVER_PORT`  | `8080`                                       |
| `DATABASE_URL` | `postgres://cafeos:cafeos@localhost:5432/cafeos?sslmode=disable` |
| `REDIS_URL`    | `redis://localhost:6379`                     |
| `JWT_SECRET`   | `dev-secret-change-in-production`            |

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

### Fase 2
- [ ] Financeiro (contas, fluxo de caixa, planejamento)
- [ ] Estoque (insumos, validade, consumo)
- [ ] Frota e equipes

### Fase 3
- [ ] Mobile offline (React Native)
- [ ] Cooperativas e consultorias
- [ ] Integrações externas

### Fase 4
- [ ] IoT (sensores, estações meteorológicas)
- [ ] IA (previsão de safra, recomendação, detecção de doenças)
- [ ] Analytics avançado

## Licença

Proprietária — todos os direitos reservados.
