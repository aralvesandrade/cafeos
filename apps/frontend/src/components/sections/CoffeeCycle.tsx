import { Flower2, Circle, Ellipsis, Leaf, Cherry } from 'lucide-react'

const phases = [
  { icon: Flower2, label: 'Florada', months: 'Set–Out', color: '#e8dcc0', desc: 'Floração do cafeeiro, definição do potencial produtivo' },
  { icon: Circle, label: 'Chumbinho', months: 'Nov–Dez', color: '#8fae72', desc: 'Formação inicial dos frutos' },
  { icon: Ellipsis, label: 'Granação', months: 'Jan–Mar', color: '#5c7a52', desc: 'Desenvolvimento e enchimento dos grãos' },
  { icon: Leaf, label: 'Maturação', months: 'Abr–Jun', color: '#d2a44c', desc: 'Amadurecimento dos frutos para colheita' },
  { icon: Cherry, label: 'Colheita', months: 'Mai–Ago', color: '#c1552f', desc: 'Colheita seletiva ou total dos frutos maduros' },
]

export function CoffeeCycle() {
  return (
    <section className="py-20 bg-background" id="coffee-cycle">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="font-display text-3xl sm:text-4xl font-semibold text-foreground mb-4">
            Ciclo do Café
          </h2>
          <p className="text-lg text-muted-foreground">
            Acompanhe cada fase do ciclo agronômico do cafeeiro, do calendário
            real da arábica brasileira, com suporte especializado da plataforma.
          </p>
        </div>

        {/* Barra de calendário: gradiente segue a cor real do fruto em cada fase */}
        <div className="h-1.5 rounded-full overflow-hidden flex mb-10">
          {phases.map((phase) => (
            <div key={phase.label} className="flex-1" style={{ backgroundColor: phase.color }} />
          ))}
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 lg:grid-cols-5 gap-6">
          {phases.map((phase, index) => (
            <div key={phase.label} className="flex flex-col items-center text-center">
              <div
                className="w-14 h-14 rounded-full flex items-center justify-center mb-4 border"
                style={{ borderColor: phase.color, backgroundColor: `${phase.color}1a` }}
              >
                <phase.icon className="h-6 w-6" style={{ color: phase.color }} />
              </div>
              <span className="text-[10px] font-mono tracking-widest text-muted-foreground mb-1">
                FASE {index + 1} · {phase.months.toUpperCase()}
              </span>
              <h3 className="font-display font-semibold text-foreground text-sm mb-1">
                {phase.label}
              </h3>
              <p className="text-xs text-muted-foreground">{phase.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
