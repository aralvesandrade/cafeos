import { useEffect, useState } from 'react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Check, X } from 'lucide-react'
import { LeadModal } from '@/components/ui/LeadModal'
import { fetchPlans, type Plan } from '@/lib/api'

function formatPrice(cents: number) {
  return (cents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL', minimumFractionDigits: cents % 100 === 0 ? 0 : 2 })
}

export function Plans() {
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [showSignup, setShowSignup] = useState(false)
  const [selectedPlan, setSelectedPlan] = useState('')

  useEffect(() => {
    fetchPlans()
      .then(setPlans)
      .catch(() => setPlans([]))
      .finally(() => setLoading(false))
  }, [])

  return (
    <section className="py-20 bg-ink" id="plans">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center max-w-3xl mx-auto mb-16">
          <h2 className="font-display text-3xl sm:text-4xl font-semibold text-parchment mb-4">
            Planos
          </h2>
          <p className="text-lg text-muted">
            Escolha o plano ideal para o seu perfil. Do pequeno produtor à
            cooperativa ou consultoria.
          </p>
        </div>

        {loading ? (
          <div className="text-center text-muted">Carregando planos...</div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {plans.map((plan) => (
              <div
                key={plan.id}
                className={`relative rounded-sm p-6 flex flex-col border ${
                  plan.featured
                    ? 'bg-ink-raised border-terreiro shadow-xl scale-105'
                    : 'bg-ink-raised border-rule'
                }`}
              >
                {plan.featured && (
                  <Badge variant="success" className="absolute -top-3 left-1/2 -translate-x-1/2">
                    Mais Popular
                  </Badge>
                )}

                <h3 className="font-display text-xl font-semibold mb-1 text-parchment">
                  {plan.name}
                </h3>
                <p className="text-sm mb-4 text-muted">
                  {plan.description}
                </p>

                <div className="mb-6 font-mono">
                  <span className="text-3xl text-gold font-medium">{formatPrice(plan.price_cents)}</span>
                  <span className="text-sm text-muted">
                    /{plan.billing_interval === 'yearly' ? 'ano' : 'mês'}
                  </span>
                </div>

                <ul className="space-y-3 mb-8 flex-1">
                  {plan.features.map((feature) => (
                    <li key={feature.label} className="flex items-center gap-2 text-sm">
                      {feature.included ? (
                        <Check className="h-4 w-4 flex-shrink-0 text-terreiro-light" />
                      ) : (
                        <X className="h-4 w-4 flex-shrink-0 text-rule" />
                      )}
                      <span className={!feature.included ? 'text-muted/50' : 'text-muted'}>
                        {feature.label}
                      </span>
                    </li>
                  ))}
                </ul>

                <Button
                  variant={plan.featured ? 'secondary' : 'outline'}
                  className="w-full"
                  onClick={() => { setSelectedPlan(plan.name); setShowSignup(true) }}
                >
                  Assinar
                </Button>
              </div>
            ))}
          </div>
        )}
      </div>

      <LeadModal open={showSignup} onClose={() => setShowSignup(false)} mode="signup" plan={selectedPlan} />
    </section>
  )
}
