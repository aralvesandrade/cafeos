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
    <section className="py-20 bg-white" id="about">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-coffee-green-dark mb-4">
            Por que o CafeOS?
          </h2>
          <p className="text-lg text-coffee-text-light">
            A cafeicultura enfrenta desafios únicos que sistemas genéricos não resolvem.
            O CafeOS nasceu para ser a plataforma especialista que o cafeicultor merece.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {challenges.map((item) => (
            <div
              key={item.title}
              className="text-center p-6 rounded-xl bg-coffee-beige hover:bg-coffee-beige-dark transition-colors"
            >
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-coffee-green/10 text-coffee-green mb-4">
                <item.icon className="h-6 w-6" />
              </div>
              <h3 className="font-semibold text-coffee-green-dark mb-2">{item.title}</h3>
              <p className="text-sm text-coffee-text-light">{item.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
