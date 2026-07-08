# Plano: Refactor de Cadastro de Propriedade, Produtor, Talhão e Indicadores

Origem: `projeto gestão agrícola.pptx` (protótipo de referência com telas de outro
sistema, metodologia SENAR/CEPEA de custo de produção). Objetivo: aproximar o
cadastro atual de `Farm`/`Plot`/`Indicator` do nível de detalhe mostrado no PPTX.

## Estado atual vs. gap

### Farm (`apps/backend/internal/domain/entity/farm.go`)
Hoje: `name, owner, location, total_area_ha, planted_area_ha`.

Falta (do PPTX, slides 1 e 3):
- Documentação: CNPJ, NIRF, número INCRA, inscrição estadual, DAP, CAR (cada um
  com flag "não possui")
- Telefone, atividades (ex: CAFEICULTURA), cultura principal/secundária, UF,
  município, endereço, sistema de produção, produto de comercialização
- Dados gerais: propriedade totalmente arrendada (bool), valor da terra nua/ha
- Divisão de área: açude/represa, benfeitorias, estradas, APP, reserva legal,
  vegetação nativa, área destinada a pecuária não contemplada, área destinada a
  agricultura não contemplada, áreas não agrícolas — todas em ha, somando com
  `total_area_ha` para chegar em "área produtiva disponível" (campo calculado)

### Produtor (não existe hoje)
`Farm.Owner` é hoje uma string livre. PPTX (slide 2) pede uma entidade própria:
CPF, nome, RG, órgão emissor, sexo, data de nascimento, estado civil, telefone,
e-mail, escolaridade. Precisa de nova entidade `Producer` (1:N com `Farm`, ou
1:1 se cada propriedade tiver um produtor responsável — a confirmar com o
usuário).

### Plot (`apps/backend/internal/domain/entity/plot.go`)
Hoje: `name, area_ha, cultivar, planting_year, altitude, soil_type`.

Falta (do PPTX, slides 4 e 7):
- Arrendado (bool), estágio (`formacao` | `producao`), irrigação (enum, ex:
  "não irrigado"), data de ativação, data de plantio, data de desativação,
  consorciada (bool) + cultura secundária, observações
- Tipo (arábica/robusta — hoje é só `cultivar` livre; talvez `cultivar` vire
  variedade/clone e um novo campo `tipo` cubra arábica/robusta)
- Custo de formação/ha (informado), custo de formação total/ha (calculado),
  vida útil (anos), espaçamento linha x planta (m), número de plantas
  (calculado ou informado)
- Divisão de área do talhão: açude, benfeitorias, estradas, APP, reserva legal
  (mesmo padrão da propriedade, em nível de talhão — slide 4)

### Indicadores (`apps/backend/internal/domain/entity/indicator.go` +
`harvest_service.go:buildIndicators`)
Hoje: 6 tipos declarados, só 4 calculados de fato (`producao_total`,
`custo_total`, `sacas_por_hectare`, `custo_por_saca`); `rentabilidade` e
`bienalidade` são código morto.

PPTX (slide 14) pede o conjunto padrão SENAR/CEPEA (~30 indicadores): Renda
Bruta, COE, COT, Custo Total, Preço Médio de Venda, Produção, Área Produção,
Produção por Área, COE/COT/CT por área e por saca, Margem Bruta (total, por
área, por saca), Margem Líquida (total, por área, por saca), Lucro (total, por
área, por saca), Taxa de Remuneração do Capital (com/sem terra), Mão de obra
familiar, Ponto de Cobertura Total (produção e produtividade), Relação
Benefício/Custo, Estoque de Capital (com/sem terra).

**Bloqueio de dados**: para calcular a maioria desses indicadores falta:
1. Preço de venda / receita por safra — hoje só existe `FinancialTransaction`
   genérico (tipo receita/despesa), sem campo de preço/quantidade vendida
   estruturado. Precisa decidir se aproveita `FinancialTransaction` ou cria
   uma entidade de venda.
2. Classificação de custo em categorias operacionais vs. capital/administrativo
   — `CostCenter` hoje só tem `Type: receita|despesa`, sem a taxonomia SENAR
   (Adubação, Colheita, Mão de obra, Mão de obra familiar, Depreciação/Capital
   etc. — lista completa no slide 10 do PPTX). COE exclui mão de obra
   familiar e depreciação; COT inclui mão de obra familiar; CT inclui
   remuneração do capital. Sem essa taxonomia não dá para separar
   automaticamente.
3. Valor da terra nua (agora corrigido — vem do novo campo em `Farm`) para
   Taxa de Remuneração do Capital e Estoque de Capital.

## Proposta de execução (fases independentes, cada uma com sua própria
aprovação de migration/schema antes de implementar)

**Fase 1 — Farm + Producer**
- Novo entity `Producer` (`apps/backend/internal/domain/entity/producer.go`),
  FK `FarmID` (ou tabela própria com N:1, a definir)
- Estender `Farm` com os campos de documentação, endereço e divisão de área
  listados acima
- Migration: adicionar ao `AutoMigrate` em `connection.go` + novo arquivo
  `002_farm_producer_fields.sql` (documentação, já que o `.sql` inicial está
  desatualizado em relação às entities atuais)
- Atualizar `FarmService`, `farm_handler.go` (DTOs de create/update), rotas de
  producer (`/api/v1/{tenant_id}/farms/{farm_id}/producer` ou embutido no
  payload de Farm — a definir)
- Atualizar `apps/admin/src/components/farms/FarmForm.tsx` e `Farms.tsx`/
  `FarmDetail.tsx` com os novos campos, organizados em seções (Propriedade /
  Documentação / Produtor / Divisão de área) como no PPTX

**Fase 2 — Plot**
- Estender `Plot` com os campos de estágio/irrigação/datas/custo de
  formação/espaçamento
- Migration equivalente
- Atualizar `PlotService` (incluir validação de datas e de estágio),
  `plot_handler.go`, e criar `PlotForm.tsx` no admin (hoje o form de talhão
  está inline em `Plots.tsx` — vale extrair como foi feito para Farm)
- Corrigir de passagem o bug encontrado: `Plots.tsx` envia `area` no payload
  mas o backend espera `area_ha`

**Fase 3 — Indicadores SENAR/CEPEA**
- Requer decisão prévia do usuário sobre os 3 bloqueios de dados acima (preço
  de venda, taxonomia de custo, e uso do valor da terra nua da Fase 1)
- Adicionar taxonomia de categoria a `CostCenter` (novo campo `SenarCategory`
  ou tabela de categorias fixas com as 18 opções do slide 10)
- Novas constantes em `IndicatorType` + lógica em `buildIndicators`
  (`harvest_service.go`) para os indicadores que já têm dado suficiente
  (Produção, Área Produção, Produção por Área, Custo Total, COT/CT por
  área/saca — falta só granularidade de custo, não dado novo)
- Indicadores que dependem de receita/preço de venda (Renda Bruta, Margem
  Bruta/Líquida, Lucro, Relação Benefício/Custo) ficam bloqueados até resolver
  o ponto 1 acima

## Decisões (confirmadas com o usuário em 2026-07-07)

1. **Producer é 1:1 com Farm** — um responsável por propriedade, FK simples
   `Producer.FarmID` (unique).
2. **Campos de documentação são todos opcionais** — CNPJ, NIRF, INCRA,
   inscrição estadual, DAP, CAR seguem o padrão do PPTX (flag "não possui" +
   campo livre, nenhum obrigatório no banco).
3. **Escopo desta rodada = Fases 1 e 2 apenas** (Farm+Producer, depois Plot).
   Fase 3 (indicadores SENAR/CEPEA) fica para um plano futuro, após o
   cadastro base estar pronto — requer decisão separada sobre como registrar
   venda/receita por safra com preço.

## Próximo passo

Detalhar a Fase 1 (schema de `Producer`, campos exatos de `Farm`, DTOs de
request/response, migration) e implementar após revisão desse detalhamento.

## Fase 3 — decisões e execução (2026-07-07)

Decisões confirmadas antes de implementar:

1. **Receita/venda**: `FinancialTransaction` será estendida (não uma entidade
   nova) com `harvest_id`, `quantity_sacas`, `unit_price` opcionais — ainda
   **não implementado**, fica para a 2ª rodada (indicadores que dependem de
   receita).
2. **Taxonomia de custo**: migração para as 18 categorias fixas do PPTX via
   novo campo `CostCenter.CostGroup` (`operacional_efetivo` |
   `mao_de_obra_familiar` | `capital_depreciacao` | `remuneracao_capital`).
   Centros de custo já existentes **não são alterados nem apagados** — ficam
   com `cost_group` vazio (não classificado) e são excluídos do cálculo de
   COE/COT/CT até reclassificação manual. O catálogo fixo
   (`entity.SenarCostCategories`) é uma constante Go, não uma tabela.
3. **Escopo**: Fatia 1 apenas (indicadores que não dependem de receita).

**Implementado (fatia 1):**
- `CostCenter.CostGroup` + catálogo `SenarCostCategories` (19 categorias do
  slide 10 do PPTX, com sua classificação COE/COT/CT já mapeada)
- Endpoint `GET /cost-centers/senar-categories` + picker no admin
  (`CostCenters.tsx`) que pré-preenche nome + classificação
- `HarvestService.calculateCostByGroupWithRepos`: soma custos (Operation,
  Maintenance, WorkShift, FinancialTransaction despesa, CostAllocation) por
  `CostGroup`, olhando o `CostCenterID` de cada um
- Novos indicadores: `area_producao`, `coe`, `coe_por_area`, `coe_por_saca`,
  `cot`, `cot_por_area`, `cot_por_saca`, `ct_producao`,
  `ct_producao_por_area`, `ct_producao_por_saca` — calculados em
  `HarvestService.Finalize` junto com os indicadores legados (que
  permanecem inalterados para não quebrar `RuleEngine`/`DashboardHandler`)
- Testado ponta a ponta via API real (cost center classificado → despesa
  vinculada → finalize → indicadores corretos no banco)

**Ainda bloqueado (fica para depois, requer nova decisão):**
- Renda Bruta, Preço Médio de Venda, Margem Bruta/Líquida, Lucro, Relação
  Benefício/Custo, Taxa de Remuneração do Capital, Estoque de Capital —
  dependem da extensão de `FinancialTransaction` com dados de venda (item 1
  acima) ainda não implementada.