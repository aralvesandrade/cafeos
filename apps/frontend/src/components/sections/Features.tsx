import { MapPin, Grid3X3, Tractor, Calendar, DollarSign, LayoutDashboard } from 'lucide-react'

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
]

export function Features() {
  return (
    <section className="py-20 bg-coffee-beige" id="features">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-coffee-green-dark mb-4">
            Funcionalidades do MVP
          </h2>
          <p className="text-lg text-coffee-text-light">
            Tudo que você precisa para gerenciar sua propriedade cafeeira do plantio à colheita.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="bg-white rounded-xl p-6 shadow-sm hover:shadow-md transition-shadow border border-gray-100"
            >
              <div className="inline-flex items-center justify-center w-10 h-10 rounded-lg bg-coffee-green text-white mb-4">
                <feature.icon className="h-5 w-5" />
              </div>
              <h3 className="font-semibold text-coffee-green-dark mb-2">{feature.title}</h3>
              <p className="text-sm text-coffee-text-light">{feature.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
