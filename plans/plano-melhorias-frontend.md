# Plano: melhorar informações mostradas no admin

## Contexto

Depois de implementar o escopo por fazenda (Farm/Plot/Financial/Stock/Fleet) e as páginas dedicadas de fazenda/talhão, mapeei todo o admin (`apps/admin/src`) e o backend (`apps/backend`) pra ver o que já existe mas não aparece na tela, e o que dá pra melhorar. Achei bem mais gordura do que esperava: formulários capturam dezenas de campos que nunca voltam pra lista/detalhe, indicadores inteiros (COE/COT/CT) são calculados e gravados no banco mas nenhuma rota os devolve, e duas features de backend completas (Budget, Cost Allocation) não têm UI nenhuma.

Documento de proposta, organizado em fases por esforço x valor, pra priorização — não é pra implementar tudo de uma vez.

## Achados (resumo)

| Achado | Tipo | Onde |
|---|---|---|
| Farm/Plot: form de cadastro captura ~35/~20 campos (documentação, divisão de área, dados do produtor); lista e detalhe mostram só um punhado | Dado capturado, não mostrado | `Farms.tsx`, `FarmDetail.tsx`, `Plots.tsx`, `PlotDetail.tsx` |
| Select de status no formulário de edição do Financeiro não funciona (`onChange` é no-op) | Bug funcional | `Financial.tsx:142` |
| Campo `notes` capturado no estado mas sem input no formulário (Estoque/movimentação, Frota/manutenção) | Bug funcional | `Stock.tsx`, `Fleet.tsx` |
| Status de Veículo, `is_active` de Trabalhador, status de Usuário/Organização: aparecem na lista mas não tem controle pra mudar | Dado capturado, não editável | `Fleet.tsx`, `Labor.tsx`, `Users.tsx`, `Organizations.tsx` |
| Operações: página existe, mas não tem criar/editar em lugar nenhum do admin | Feature ausente | `Operations.tsx` |
| Indicadores COE/COT/CT, sacas/ha, custo/saca, rentabilidade: calculados e gravados na tabela `indicators` ao finalizar safra, nenhuma rota devolve isso | Backend pronto, sem API/UI | `harvest_service.go:392-438`, `dashboard_handler.go` |
| Budget (orçado): CRUD completo no backend, zero linha de frontend | Feature backend sem UI | `budget_handler.go`, `budget_service.go` |
| Cost Allocation (rateio por talhão): CRUD + cálculo por talhão completo no backend, zero linha de frontend; e a resposta nem inclui o nome do talhão (tipo `allocationItemResponse` existe mas nunca é usado) | Feature backend sem UI + bug menor | `cost_allocation_handler.go:31-36` |
| Rollup de custo por equipe/trabalhador (Labor) e por veículo (Fleet): não existe, só lista bruta | Feature ausente | `Labor.tsx`, `Fleet.tsx` |
| Rule Engine (baixa produtividade, alto custo/saca) + `AlertGenerated`: construído mas nunca chamado por nada — código órfão | Código morto | `rule_engine.go`, `event/handler.go` |

## Fases propostas (esforço crescente)

### Fase 1 — Consertos rápidos (bugs de UI já existentes)
- Financeiro: consertar o select de status no formulário de edição (hoje é decorativo).
- Estoque (movimentação) e Frota (manutenção): adicionar o input de `notes` que falta.
- Veículo, Trabalhador, Usuário, Organização: adicionar o controle que falta pra editar status/ativo (hoje só dá pra ver, não mudar, exceto recriando).

### Fase 2 — Mostrar o que já é capturado (Farm/Plot)
- `FarmDetail.tsx`: adicionar seções pra documentação (CNPJ/NIRF/INCRA/inscrição estadual/DAP/CAR), divisão de área (açude, benfeitorias, estradas, APP, reserva legal, vegetação nativa, etc.), dados completos do produtor (hoje só nome/CPF/telefone/email aparecem; RG, órgão emissor, sexo, nascimento, estado civil, escolaridade ficam de fora).
- `PlotDetail.tsx`: adicionar a divisão de área do talhão (açude, benfeitorias, estradas, APP, reserva legal) e a data de desativação.
- Considerar mover fazenda/talhão pra layout em abas (Geral / Documentação / Áreas / Produtor) já que o formulário é grande — mesmo espírito da página dedicada que acabamos de criar.

### Fase 3 — CRUD de Operações
- `Operations.tsx` ganha criar/editar (hoje só lista) — reaproveitar o padrão de `OperationHandler.Create` já existente no backend, só falta o formulário no admin.

### Fase 4 — Indicadores de safra na UI
- Backend: novo endpoint (ou campo extra em `GetProduction`/`GetByID` de Harvest) devolvendo os indicadores já calculados (`IndSacasHA`, `IndCustoSaca`, `IndRentabilidade`\*, COE/COT/CT e variantes por área/saca) — hoje ficam presos na tabela `indicators`.
- `HarvestDetail.tsx`: mostrar esses indicadores pra safras finalizadas (cards: sacas/ha, custo/saca, COE, COT, CT).
- Dashboard: usar COE/COT/CT pra enriquecer os cards hoje limitados a produção/custo total.
- \*Nota: `IndRentabilidade` e `IndBienalidade` estão declarados mas **nunca calculados** em lugar nenhum (confirmado) — ou implementamos o cálculo ou removemos as constantes mortas, não faz sentido expor um indicador que nunca tem valor real.

### Fase 5 — Budget (orçado x realizado)
- Backend: `Budget` hoje só tem `PlannedAmount`, sem comparação com gasto real — precisa somar o realizado (via Operation/Maintenance/WorkShift/FinancialTransaction/CostAllocation do mesmo centro de custo+safra, mesma lógica que `calculateTotalCostWithRepos` já usa) e devolver a variância.
- Admin: nova página `Budget.tsx` (CRUD + visão orçado x realizado por centro de custo/safra), com link a partir de `HarvestDetail`.

### Fase 6 — Cost Allocation (rateio por talhão)
- Backend: corrigir a resposta pra usar o tipo `allocationItemResponse` já declarado (mas nunca usado) enriquecendo `PlotName` — hoje só devolve `PlotID` cru.
- Admin: nova página `CostAllocations.tsx` — criar rateio (proporcional por área ou percentual customizado), listar com a tabela de valores por talhão.

### Fase 7 — Rollups de custo
- `Labor.tsx`: total de horas/custo por trabalhador e por equipe.
- `Fleet.tsx`: total de custo de manutenção por veículo (útil pra decidir troca de equipamento).

### Fase 8 — Reativar Rule Engine / Alertas (avaliar se vale)
- Hoje `RuleEngine` (baixa produtividade, alto custo/saca) e o evento `AlertGenerated` existem mas nada os aciona — nenhuma tela, nenhuma tabela de alerta.
- Se topar: chamar `RuleEngine.Evaluate` ao finalizar safra, persistir o alerta gerado (entidade nova + migration), endpoint `/alerts`, e um sino de notificação no admin.
- Maior escopo que as fases anteriores — só entra se você achar que vale a pena reconectar em vez de deixar como está.

## Verificação (quando cada fase for implementada, não agora)
- Backend: `go build`/`go vet`/`go test` a cada fase que tocar Go.
- Admin: `npm run build`/`npm run lint` a cada fase.
- Testar manualmente o fluxo afetado (formulário salva o campo, tela mostra o valor certo).
