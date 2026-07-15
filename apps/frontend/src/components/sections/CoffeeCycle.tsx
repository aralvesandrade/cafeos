import { Flower2, Circle, Ellipsis, Leaf, Cherry } from 'lucide-react'

const phases = [
  { icon: Flower2, label: 'Florada', desc: 'Floração do cafeeiro, definição potencial produtivo' },
  { icon: Circle, label: 'Chumbinho', desc: 'Formação inicial dos frutos' },
  { icon: Ellipsis, label: 'Granação', desc: 'Desenvolvimento e enchimento dos grãos' },
  { icon: Leaf, label: 'Maturação', desc: 'Amadurecimento dos frutos para colheita' },
  { icon: Cherry, label: 'Colheita', desc: 'Colheita seletiva ou total dos frutos maduros' },
]

export function CoffeeCycle() {
  return (
    <section className="py-20 bg-white" id="coffee-cycle">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-coffee-green-dark mb-4">
            Ciclo do Café
          </h2>
          <p className="text-lg text-coffee-text-light">
            Acompanhe cada fase do ciclo agronômico do cafeeiro com suporte
            especializado da plataforma.
          </p>
        </div>

        <div className="relative">
          <div className="hidden lg:block absolute top-1/2 left-0 right-0 h-0.5 bg-coffee-green/20 -translate-y-1/2" />

          <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-5 gap-6">
            {phases.map((phase, index) => (
              <div key={phase.label} className="relative flex flex-col items-center text-center">
                <div className="relative z-10 w-16 h-16 rounded-full bg-coffee-beige flex items-center justify-center mb-4 border-2 border-coffee-green/30">
                  <phase.icon className="h-7 w-7 text-coffee-green" />
                </div>
                <span className="text-xs font-bold text-coffee-green mb-1">
                  Fase {index + 1}
                </span>
                <h3 className="font-semibold text-coffee-green-dark text-sm mb-1">
                  {phase.label}
                </h3>
                <p className="text-xs text-coffee-text-light">{phase.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
