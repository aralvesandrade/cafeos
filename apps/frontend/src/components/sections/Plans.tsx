import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Check, X } from 'lucide-react'
import { LeadModal } from '@/components/ui/LeadModal'

const plans = [
  {
    name: 'Grátis',
    price: 'R$ 0',
    period: '/mês',
    desc: 'Para pequenos produtores que estão começando',
    featured: false,
    features: [
      { label: '1 fazenda', included: true },
      { label: '5 talhões', included: true },
      { label: '50 operações/mês', included: true },
      { label: '3 safras', included: true },
      { label: '2 usuários', included: true },
      { label: 'Dashboard básico', included: true },
      { label: 'Relatórios CSV/PDF', included: false },
      { label: 'White Label', included: false },
      { label: 'Suporte prioritário', included: false },
    ],
  },
  {
    name: 'Pro',
    price: 'R$ 97',
    period: '/mês',
    desc: 'Para médios produtores com operações mais complexas',
    featured: true,
    features: [
      { label: '10 fazendas', included: true },
      { label: '50 talhões', included: true },
      { label: 'Operações ilimitadas', included: true },
      { label: 'Safras ilimitadas', included: true },
      { label: '10 usuários', included: true },
      { label: 'Dashboard avançado', included: true },
      { label: 'Relatórios CSV/PDF', included: true },
      { label: 'White Label', included: false },
      { label: 'Suporte por email', included: true },
    ],
  },
  {
    name: 'Cooperativa',
    price: 'R$ 297',
    period: '/mês',
    desc: 'Para cooperativas que gerenciam múltiplos associados',
    featured: false,
    features: [
      { label: 'Fazendas ilimitadas', included: true },
      { label: 'Talhões ilimitados', included: true },
      { label: 'Operações ilimitadas', included: true },
      { label: 'Safras ilimitadas', included: true },
      { label: '50 usuários', included: true },
      { label: 'Dashboard consolidado', included: true },
      { label: 'Benchmarking associados', included: true },
      { label: 'White Label', included: false },
      { label: 'Suporte prioritário', included: true },
    ],
  },
  {
    name: 'Consultoria',
    price: 'R$ 497',
    period: '/mês',
    desc: 'Para consultorias que atendem múltiplos clientes',
    featured: false,
    features: [
      { label: 'Multi-cliente', included: true },
      { label: 'Talhões ilimitados', included: true },
      { label: 'Operações ilimitadas', included: true },
      { label: 'Safras ilimitadas', included: true },
      { label: '30 usuários', included: true },
      { label: 'Dashboard por cliente', included: true },
      { label: 'Relatórios técnicos', included: true },
      { label: 'White Label', included: true },
      { label: 'Suporte dedicado', included: true },
    ],
  },
]

export function Plans() {
  const [showSignup, setShowSignup] = useState(false)

  return (
    <section className="py-20 bg-coffee-beige" id="plans">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-coffee-green-dark mb-4">
            Planos
          </h2>
          <p className="text-lg text-coffee-text-light">
            Escolha o plano ideal para o seu perfil. Do pequeno produtor à
            cooperativa ou consultoria.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={`relative rounded-2xl p-6 flex flex-col ${
                plan.featured
                  ? 'bg-coffee-green-dark text-white ring-2 ring-coffee-green shadow-xl scale-105'
                  : 'bg-white text-coffee-text border border-gray-200'
              }`}
            >
              {plan.featured && (
                <Badge variant="success" className="absolute -top-3 left-1/2 -translate-x-1/2">
                  Mais Popular
                </Badge>
              )}

              <h3 className={`text-xl font-bold mb-1 ${plan.featured ? 'text-white' : 'text-coffee-green-dark'}`}>
                {plan.name}
              </h3>
              <p className={`text-sm mb-4 ${plan.featured ? 'text-coffee-beige/80' : 'text-coffee-text-light'}`}>
                {plan.desc}
              </p>

              <div className="mb-6">
                <span className="text-3xl font-bold">{plan.price}</span>
                <span className={`text-sm ${plan.featured ? 'text-coffee-beige/80' : 'text-coffee-text-light'}`}>
                  {plan.period}
                </span>
              </div>

              <ul className="space-y-3 mb-8 flex-1">
                {plan.features.map((f) => (
                  <li key={f.label} className="flex items-center gap-2 text-sm">
                    {f.included ? (
                      <Check className={`h-4 w-4 flex-shrink-0 ${plan.featured ? 'text-green-300' : 'text-coffee-green'}`} />
                    ) : (
                      <X className={`h-4 w-4 flex-shrink-0 ${plan.featured ? 'text-white/40' : 'text-gray-300'}`} />
                    )}
                    <span className={!f.included && !plan.featured ? 'text-gray-400' : ''}>
                      {f.label}
                    </span>
                  </li>
                ))}
              </ul>

              <Button
                variant={plan.featured ? 'secondary' : 'outline'}
                className="w-full"
                onClick={() => setShowSignup(true)}
              >
                {plan.name === 'Grátis' ? 'Começar Grátis' : 'Assinar'}
              </Button>
            </div>
          ))}
        </div>
      </div>

      <LeadModal open={showSignup} onClose={() => setShowSignup(false)} mode="signup" />
    </section>
  )
}
