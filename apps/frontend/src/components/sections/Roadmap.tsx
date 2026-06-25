const phases = [
  {
    phase: 'Fase 1',
    title: 'MVP',
    status: 'completed',
    items: [
      'Gestão de fazendas e talhões',
      'Operações agrícolas',
      'Gestão de safras',
      'Custos agrícolas',
      'Dashboard inicial',
    ],
  },
  {
    phase: 'Fase 2',
    title: 'Financeiro e Operacional',
    status: 'current',
    items: [
      'Financeiro (contas, fluxo de caixa)',
      'Estoque de insumos',
      'Gestão de frota',
      'Gestão de equipes',
    ],
  },
  {
    phase: 'Fase 3',
    title: 'Expansão',
    status: 'upcoming',
    items: [
      'Mobile offline (React Native)',
      'Cooperativas e consultorias',
      'Pós-colheita',
      'Integrações externas',
    ],
  },
  {
    phase: 'Fase 4',
    title: 'Inteligência',
    status: 'upcoming',
    items: [
      'IoT (sensores e telemetria)',
      'IA para previsão de safra',
      'Recomendação de adubação',
      'Detecção de doenças',
    ],
  },
]

export function Roadmap() {
  return (
    <section className="py-20 bg-coffee-beige" id="roadmap">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-coffee-green-dark mb-4">
            Roadmap
          </h2>
          <p className="text-lg text-coffee-text-light">
            Conheça nossa jornada de evolução. O CafeOS está sendo construído
            em fases para entregar valor contínuo.
          </p>
        </div>

        <div className="relative">
          <div className="hidden lg:block absolute top-1/2 left-0 right-0 h-0.5 bg-gray-200 -translate-y-1/2" />

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {phases.map((phase) => (
              <div
                key={phase.phase}
                className={`relative rounded-xl p-6 ${
                  phase.status === 'completed'
                    ? 'bg-coffee-green-dark text-white'
                    : phase.status === 'current'
                      ? 'bg-coffee-green text-white'
                      : 'bg-white border border-gray-200'
                }`}
              >
                <span
                  className={`text-xs font-bold uppercase tracking-wider ${
                    phase.status === 'upcoming' ? 'text-coffee-text-light' : 'text-coffee-beige/80'
                  }`}
                >
                  {phase.phase}
                </span>
                <h3
                  className={`text-lg font-bold mt-1 mb-3 ${
                    phase.status === 'upcoming' ? 'text-coffee-green-dark' : 'text-white'
                  }`}
                >
                  {phase.title}
                </h3>

                <ul className="space-y-2">
                  {phase.items.map((item) => (
                    <li
                      key={item}
                      className={`text-sm flex items-start gap-2 ${
                        phase.status === 'upcoming' ? 'text-coffee-text-light' : 'text-white/80'
                      }`}
                    >
                      <span className="mt-0.5">•</span>
                      {item}
                    </li>
                  ))}
                </ul>

                {phase.status === 'completed' && (
                  <span className="absolute top-2 right-2 text-xs bg-white/20 px-2 py-1 rounded-full">
                    ✅ Concluído
                  </span>
                )}
                {phase.status === 'current' && (
                  <span className="absolute top-2 right-2 text-xs bg-white/20 px-2 py-1 rounded-full">
                    🔄 Em andamento
                  </span>
                )}
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  )
}
