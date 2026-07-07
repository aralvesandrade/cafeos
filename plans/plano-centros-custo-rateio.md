# Plano: Centros de Custo e Rateio

## Origem

Arquivo `projeto gestão agrícola.pptx` — apresentação sobre metodologia SENAR de gestão agrícola.

## Gaps identificados vs. sistema atual

| Tópico | Situação | GAP |
|---|---|---|
| Cadastro de áreas (farms/plots) | `Farm` + `Plot` existem | Ok |
| Centros de custo | Não existe entidade | **Ausente** |
| Plano de contas | `FinancialTransaction.Category` string livre | **Ausente** |
| Custo por operação | `Operation.Cost` existe | Sem vínculo com safra |
| Custo de manutenção | `Maintenance.Cost` existe | Não entra nos indicadores |
| Custo de mão de obra | `WorkShift.Cost` existe | Não entra nos indicadores |
| Rateio entre talhões | Não existe | **Ausente** |
| Orçado vs Realizado | Não existe | **Ausente** |
| Regime competência | Despesa registrada na compra | Deve ser no uso |
| Cálculo de safra | Varre tudo do tenant | Deve filtrar por safra |

## Stack

Mesma do backend: Go 1.26 + GORM + PostgreSQL.

## Estrutura do plano

---

### Fase 1 — Entidade CostCenter (Plano de Contas Simples)

#### Backend

Nova entidade `CostCenter`:

| Campo | Tipo | Nota |
|---|---|---|
| ID | UUID PK | |
| TenantID | UUID FK | Index |
| Name | string | Ex: "Adubos", "Mão de Obra", "Frete" |
| Code | string | Código curto único (ex: "ADUBOS") |
| Type | enum | `receita` ou `despesa` |
| Description | string | Opcional |

Seed inicial com centros de custo padrão (base SENAR):
- **Despesas**: Adubos, Defensivos, Combustíveis, Mão de Obra, Frete, Manutenção, Irrigação, Análise de Solo, Outros Insumos, Serviços Terceiros, Energia, Depreciação, Administrativo, Outras Despesas
- **Receitas**: Venda de Café, Venda de Mudas, Outras Receitas

**Arquivos:**
- `apps/backend/internal/domain/entity/cost_center.go` — struct + validação
- `apps/backend/internal/domain/repository/cost_center_repository.go` — interface
- `apps/backend/internal/infra/db/repository/cost_center_repo.go` — GORM impl
- `apps/backend/internal/api/handler/cost_center_handler.go` — CRUD handlers
- `apps/backend/cmd/api/routes.go` — rotas `/api/v1/{tenant_id}/cost-centers`

#### Admin

- Tela CRUD de Centros de Custo (tabela + formulário)
- Substituir campo `Category` nos formulários de `FinancialTransaction`, `Operation`, `Maintenance`, `WorkShift` por select de `CostCenter`

#### Migração

1. Criar tabela `cost_centers`
2. Seed com centros de custo padrão
3. Add FK `cost_center_id` nullable em:
   - `financial_transactions` (drop column `category`)
   - `operations`
   - `maintenances`
   - `work_shifts`
4. Migrar dados existentes: associar ao centro de custo mais provável pela descrição/tipo

---

### Fase 2 — Rateio entre Talhões

Dois modos de alocação de custo:

| Modo | Como funciona | Onde |
|---|---|---|
| Direto | Vinculado a 1 Plot | `Operation.PlotID` (já existe) |
| Rateio | Custo geral distribuído entre N talhões | Nova entidade `CostAllocation` |

#### Entidade `CostAllocation`

| Campo | Tipo | Nota |
|---|---|---|
| ID | UUID PK | |
| TenantID | UUID FK | |
| HarvestID | UUID FK | Vinculado a safra |
| CostCenterID | UUID FK | |
| Description | string | Ex: "Frete geral jan/2026" |
| TotalAmount | float64 | Valor total a ratear |
| Method | enum | `area_proportional`, `custom_percentage` |
| Date | timestamp | |

#### Entidade `CostAllocationItem`

| Campo | Tipo | Nota |
|---|---|---|
| ID | UUID PK | |
| AllocationID | UUID FK | |
| PlotID | UUID FK | |
| Amount | float64 | Valor rateado |
| Percentage | float64 | % aplicada no rateio |

Regras:
- `area_proportional`: backend calcula peso = `Plot.AreaHA / Σ(AreaHA dos plots ativos)`, distribui `TotalAmount` proporcionalmente
- `custom_percentage`: usuário informa % para cada talhão, backend valida se soma = 100%

#### Impacto no cálculo de safra

Tanto `Operation.Cost` (direto) quanto `CostAllocation.TotalAmount` (rateio) entram no custo por talhão da safra.

---

### Fase 3 — Operações Vinculadas a Safra

#### Migração

- Add FK `harvest_id` nullable em `operations`
- Criar índice composto `(tenant_id, harvest_id)`

#### Cálculo de custo em `HarvestService.Finalize`

Antes: `totalCost = Σ(Operation.Cost)` sem filtro.

Depois:

```go
var operationCost float64
db.Where("harvest_id = ?", harvest.ID).Model(&Operation{}).Select("COALESCE(SUM(cost), 0)").Scan(&operationCost)

var maintenanceCost float64
db.Where("tenant_id = ? AND date BETWEEN ? AND ?", harvest.TenantID, startDate, endDate).
  Model(&Maintenance{}).Select("COALESCE(SUM(cost), 0)").Scan(&maintenanceCost)

var laborCost float64
db.Where("tenant_id = ? AND date BETWEEN ? AND ?", harvest.TenantID, startDate, endDate).
  Model(&WorkShift{}).Select("COALESCE(SUM(cost), 0)").Scan(&laborCost)

var allocationCost float64
db.Where("harvest_id = ?", harvest.ID).Model(&CostAllocation{}).
  Select("COALESCE(SUM(total_amount), 0)").Scan(&allocationCost)

var financialCost float64
db.Where("tenant_id = ? AND type = 'despesa' AND date BETWEEN ? AND ?", harvest.TenantID, startDate, endDate).
  Model(&FinancialTransaction{}).Select("COALESCE(SUM(amount), 0)").Scan(&financialCost)

totalCost := operationCost + maintenanceCost + laborCost + allocationCost + financialCost
```

---

### Fase 4 — Orçado vs Realizado

#### Entidade `Budget`

| Campo | Tipo | Nota |
|---|---|---|
| ID | UUID PK | |
| TenantID | UUID FK | |
| HarvestID | UUID FK | Orçamento por safra |
| CostCenterID | UUID FK | |
| PlannedAmount | float64 | Valor orçado |
| Description | string | Opcional |

- Unique constraint: `(tenant_id, harvest_id, cost_center_id)`

#### Relatório

Query base:

```sql
SELECT
  cc.id,
  cc.name,
  b.planned_amount,
  COALESCE(actual_table.total_cost, 0) AS actual_cost,
  COALESCE(b.planned_amount - actual_table.total_cost, -b.planned_amount) AS variance
FROM cost_centers cc
LEFT JOIN budgets b ON b.cost_center_id = cc.id AND b.harvest_id = ?
LEFT JOIN (
  -- união de todos os custos realizados
  SELECT cost_center_id, SUM(cost_value) AS total_cost
  FROM (
    SELECT cost_center_id, cost AS cost_value FROM operations WHERE harvest_id = ?
    UNION ALL
    SELECT cost_center_id, cost AS cost_value FROM work_shifts ws
      JOIN operations o ON o.id = ws.operation_id WHERE o.harvest_id = ?
    -- etc
  ) sub
  GROUP BY cost_center_id
) actual_table ON actual_table.cost_center_id = cc.id
WHERE cc.tenant_id = ?
ORDER BY cc.name
```

#### Dashboard

- Gráfico de barras: orçado vs realizado por centro de custo
- Gráfico de pizza: distribuição de custos por centro
- Tabela com variação (R$ e %)
- Filtro por safra

---

### Fase 5 — Metodologia SENAR (Regime de Competência)

Regra de negócio:
- **Compra de insumo** → entrada em estoque (`StockMovement.type = "in"`)
- **Uso do insumo** → operação agrícola (`Operation`) + baixa em estoque (`StockMovement.type = "out"`)
- **Despesa financeira** → registrada no momento do uso, não da compra

Camada de validação:
- Ao criar `FinancialTransaction` com type `despesa` e centro de custo = insumo, verificar se há baixa de estoque correspondente (warning, não block)
- Documentar no campo `notes` ou em property `accrual_date`

---

## Resumo de arquivos novos/alterados

### Novos
| Arquivo | Descrição |
|---|---|
| `domain/entity/cost_center.go` | Entidade CostCenter |
| `domain/entity/cost_allocation.go` | Entidades CostAllocation + CostAllocationItem |
| `domain/entity/budget.go` | Entidade Budget |
| `domain/repository/cost_center_repository.go` | Interface CostCenterRepository |
| `domain/repository/cost_allocation_repository.go` | Interface CostAllocationRepository |
| `domain/repository/budget_repository.go` | Interface BudgetRepository |
| `infra/db/repository/cost_center_repo.go` | GORM impl CostCenter |
| `infra/db/repository/cost_allocation_repo.go` | GORM impl CostAllocation |
| `infra/db/repository/budget_repo.go` | GORM impl Budget |
| `api/handler/cost_center_handler.go` | CRUD handler |
| `api/handler/cost_allocation_handler.go` | CRUD handler |
| `api/handler/budget_handler.go` | CRUD handler |

### Alterados
| Arquivo | Mudança |
|---|---|
| `domain/entity/operation.go` | Add `HarvestID`, `CostCenterID` |
| `domain/entity/financial.go` | Substituir `Category` por `CostCenterID` |
| `domain/entity/fleet.go` (Maintenance) | Add `CostCenterID` |
| `domain/entity/labor.go` (WorkShift) | Add `CostCenterID` |
| `domain/entity/harvest.go` | Add `StartDate`, `EndDate` (para filtro) |
| `domain/service/harvest_service.go` | Novo cálculo de custo incluindo todas as fontes |
| `cmd/api/routes.go` | Novas rotas |

## Prioridade de implementação

1. **Fase 1** — CostCenter (base para tudo)
2. **Fase 3** — HarvesterID em Operation + filtro de datas (corrige indicadores)
3. **Fase 4** — Budget (orçado vs realizado)
4. **Fase 2** — Rateio entre talhões
5. **Fase 5** — Regime de competência (validações)
