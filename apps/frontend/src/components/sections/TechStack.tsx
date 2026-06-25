import { Database, Container, Cpu, Brain, Smartphone, Cloud } from 'lucide-react'

const techs = [
  { icon: Container, label: 'Go', desc: 'API REST robusta e performática' },
  { icon: Database, label: 'PostgreSQL + Redis', desc: 'Banco relacional e cache' },
  { icon: Cloud, label: 'Docker + K8s', desc: 'Infraestrutura escalável' },
  { icon: Smartphone, label: 'React + React Native', desc: 'Web e Mobile' },
  { icon: Cpu, label: 'IoT Ready', desc: 'Sensores e telemetria' },
  { icon: Brain, label: 'IA & ML', desc: 'Previsão e recomendação' },
]

export function TechStack() {
  return (
    <section className="py-20 bg-white" id="tech">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-coffee-green-dark mb-4">
            Stack Tecnológica
          </h2>
          <p className="text-lg text-coffee-text-light">
            Construído com tecnologias modernas e preparado para o futuro com
            IoT e Inteligência Artificial.
          </p>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-6">
          {techs.map((tech) => (
            <div
              key={tech.label}
              className="text-center p-6 rounded-xl border border-gray-100 hover:border-coffee-green/30 hover:shadow-md transition-all"
            >
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-coffee-beige text-coffee-green mb-3">
                <tech.icon className="h-6 w-6" />
              </div>
              <h3 className="font-semibold text-sm text-coffee-green-dark">{tech.label}</h3>
              <p className="text-xs text-coffee-text-light mt-1">{tech.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}
