package entity

// SenarCostCategory is a pre-classified despesa cost center category from
// the SENAR/CEPEA coffee production cost methodology. It is a fixed catalog
// (not persisted) offered to organizations when registering new cost centers, so
// that COE/COT/CT can be computed automatically without asking the user to
// classify free-text categories by hand.
type SenarCostCategory struct {
	Name      string    `json:"name"`
	CostGroup CostGroup `json:"cost_group"`
}

// SenarCostCategories lists the 18 despesa categories from the reference
// methodology. Existing organization-created cost centers are unaffected by this
// catalog — it only pre-fills Name/CostGroup for newly created ones.
var SenarCostCategories = []SenarCostCategory{
	{Name: "Adubação via folha", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Adubação via solo", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Arrendamento ou aluguel", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Colheita", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Comercialização", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Condução da lavoura", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Controle de plantas daninhas", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Controle de pragas e doenças", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Despesa administrativa", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Empréstimo - pagamento", CostGroup: CostGroupRemuneracaoCapital},
	{Name: "Formação e Recuperação", CostGroup: CostGroupCapitalDepreciacao},
	{Name: "Gestão Técnica", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Imposto e taxa", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Investimento - aquisição de bem", CostGroup: CostGroupCapitalDepreciacao},
	{Name: "Irrigação", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Mão de obra", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Mão de obra familiar", CostGroup: CostGroupMaoDeObraFamiliar},
	{Name: "Poda e desbrota", CostGroup: CostGroupOperacionalEfetivo},
	{Name: "Pós-colheita", CostGroup: CostGroupOperacionalEfetivo},
}
