import { TrendingUp, TrendingDown, Wallet, Landmark } from 'lucide-react'

const indicators = [
  {
    icon: TrendingUp,
    label: 'Sacas por Hectare',
    value: 'sacas/ha',
    desc: 'Produtividade da lavoura por área plantada',
  },
  {
    icon: TrendingDown,
    label: 'Custo por Saca',
    value: 'R$/saca',
    desc: 'Custo operacional total dividido pela produção',
  },
  {
    icon: Wallet,
    label: 'COE',
    value: 'R$',
    desc: 'Custo Operacional Efetivo: desembolsos diretos da safra',
  },
  {
    icon: Landmark,
    label: 'COT',
    value: 'R$',
    desc: 'Custo Operacional Total: COE + mão de obra familiar e depreciação',
  },
]

export function Indicators() {
  return (
    <section className="py-20 bg-card" id="indicators">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-14">
          <h2 className="font-display text-3xl sm:text-4xl font-semibold text-foreground mb-4">
            Indicadores Estratégicos
          </h2>
          <p className="text-lg text-muted-foreground">
            Acompanhe os principais indicadores da sua produção e tome decisões
            baseadas em dados precisos.
          </p>
        </div>

        {/* Ficha de pesagem — mesma linguagem do manifesto do Hero, em destaque */}
        <div className="relative bg-background border border-border rounded-sm p-8 sm:p-10">
          <div className="absolute -top-3 right-6 sm:right-10 bg-gold text-gold-foreground text-[10px] font-mono uppercase tracking-widest px-3 py-1 rounded-sm rotate-2 shadow-md">
            Safra 2025 · Certificado CafeOS
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 mt-2">
            {indicators.map((indicator, i) => (
              <div
                key={indicator.label}
                className={`font-mono ${i > 0 ? 'sm:pl-8 sm:border-l sm:border-border' : ''}`}
              >
                <div className="flex items-center gap-2 text-[10px] tracking-widest text-muted-foreground mb-2">
                  <indicator.icon className="h-3.5 w-3.5 text-primary" />
                  {indicator.label.toUpperCase()}
                </div>
                <div className="text-2xl text-gold font-medium mb-1">{indicator.value}</div>
                <p className="text-xs text-muted-foreground font-sans">{indicator.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
