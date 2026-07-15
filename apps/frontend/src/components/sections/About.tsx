import { Sprout, TrendingUp, Shield, Cpu } from 'lucide-react'

const challenges = [
  {
    icon: Sprout,
    title: 'Falta de visibilidade',
    desc: 'Sem dados centralizados, o produtor perde o controle sobre custos, produtividade e operações do campo.',
  },
  {
    icon: TrendingUp,
    title: 'Baixa produtividade',
    desc: 'Sem indicadores precisos, fica difícil identificar gargalos e tomar decisões baseadas em dados.',
  },
  {
    icon: Shield,
    title: 'Rastreabilidade',
    desc: 'Certificações e compliance exigem histórico completo das operações, insumos e produção.',
  },
  {
    icon: Cpu,
    title: 'Tecnologia defasada',
    desc: 'Planilhas e sistemas genéricos não atendem às necessidades específicas da cafeicultura.',
  },
]

export function About() {
  return (
    <section className="py-20 bg-ink" id="about">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="font-display text-3xl sm:text-4xl font-semibold text-parchment mb-4">
            Por que o CafeOS?
          </h2>
          <p className="text-lg text-muted">
            A cafeicultura enfrenta desafios únicos que sistemas genéricos não resolvem.
            O CafeOS nasceu para ser a plataforma especialista que o cafeicultor merece.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {challenges.map((item) => (
            <div
              key={item.title}
              className="relative bg-ink-raised border border-rule rounded-sm p-6 hover:border-terreiro/50 transition-colors"
            >
              <div className="absolute top-4 left-4 w-1.5 h-1.5 rounded-full bg-rule" />
              <div className="inline-flex items-center justify-center w-11 h-11 rounded-sm bg-terreiro/10 text-terreiro-light mb-4 mt-2">
                <item.icon className="h-5 w-5" />
              </div>
              <h3 className="font-display font-semibold text-parchment mb-2">{item.title}</h3>
              <p className="text-sm text-muted">{item.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
