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
    <section className="py-20 bg-coffee-green-dark text-white" id="indicators">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            Indicadores Estratégicos
          </h2>
          <p className="text-lg text-coffee-beige/80">
            Acompanhe os principais indicadores da sua produção e tome decisões
            baseadas em dados precisos.
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
          {indicators.map((indicator) => (
            <div
              key={indicator.label}
              className="text-center p-6 rounded-xl bg-white/10 backdrop-blur-sm"
            >
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-coffee-green/30 mb-4">
                <indicator.icon className="h-6 w-6" />
              </div>
              <span className="text-xs font-medium text-coffee-beige/80 uppercase tracking-wider">
                {indicator.value}
              </span>
              <h3 className="font-semibold text-lg mt-1 mb-2">{indicator.label}</h3>
              <p className="text-sm text-coffee-beige/70">{indicator.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
