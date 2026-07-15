import {
  MapPin,
  Grid3X3,
  Tractor,
  Calendar,
  DollarSign,
  LayoutDashboard,
  Wallet,
  SplitSquareHorizontal,
  Bell,
} from 'lucide-react'

const features = [
  {
    icon: MapPin,
    title: 'Gestão de Fazendas',
    desc: 'Cadastro completo com dados agronômicos, área total e plantada, georreferenciamento.',
  },
  {
    icon: Grid3X3,
    title: 'Gestão de Talhões',
    desc: 'Talhões por fazenda com cultivar, solo, altitude, ano de plantio e área.',
  },
  {
    icon: Tractor,
    title: 'Operações Agrícolas',
    desc: 'Registro de adubação, pulverização, irrigação, poda e colheita com custos.',
  },
  {
    icon: Calendar,
    title: 'Gestão de Safras',
    desc: 'Estimativa e produção realizada por talhão, histórico comparativo entre safras.',
  },
  {
    icon: DollarSign,
    title: 'Custos Agrícolas',
    desc: 'Custos por operação, talhão e safra. Cálculo de custo por hectare e por saca.',
  },
  {
    icon: LayoutDashboard,
    title: 'Dashboard',
    desc: 'Visão consolidada de produção, custos, evolução da safra e operações recentes.',
  },
  {
    icon: Wallet,
    title: 'Orçamento',
    desc: 'Orçado x realizado por centro de custo e safra, com variação e execução.',
  },
  {
    icon: SplitSquareHorizontal,
    title: 'Rateio de Custo',
    desc: 'Distribuição de custos por talhão, proporcional por área ou percentual customizado.',
  },
  {
    icon: Bell,
    title: 'Alertas Inteligentes',
    desc: 'Motor de regras identifica baixa produtividade e custo elevado automaticamente.',
  },
]

export function Features() {
  return (
    <section className="py-20 bg-ink-raised" id="features">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="font-display text-3xl sm:text-4xl font-semibold text-parchment mb-4">
            Funcionalidades
          </h2>
          <p className="text-lg text-muted">
            Tudo que você precisa para gerenciar sua propriedade cafeeira do
            plantio à colheita.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="relative bg-ink border border-rule rounded-sm p-6 hover:border-terreiro/50 transition-colors"
            >
              <div className="absolute top-4 left-4 w-1.5 h-1.5 rounded-full bg-rule" />
              <div className="inline-flex items-center justify-center w-11 h-11 rounded-sm bg-terreiro/10 text-terreiro-light mb-4 mt-2">
                <feature.icon className="h-5 w-5" />
              </div>
              <h3 className="font-display font-semibold text-parchment mb-2">{feature.title}</h3>
              <p className="text-sm text-muted">{feature.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
