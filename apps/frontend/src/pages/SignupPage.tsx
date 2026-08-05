import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Sprout, CheckCircle } from 'lucide-react'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { RequiredLegend } from '@/components/ui/required-legend'
import { Button } from '@/components/ui/button'
import { registerSignup, fetchPlans, type Plan } from '@/lib/api'

function formatPrice(cents: number) {
  return (cents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL', minimumFractionDigits: cents % 100 === 0 ? 0 : 2 })
}

export function SignupPage() {
  const [searchParams] = useSearchParams()
  const planSlug = searchParams.get('plano') || ''
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null)
  const [form, setForm] = useState({
    name: '',
    email: '',
    password: '',
  })
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [done, setDone] = useState(false)

  useEffect(() => {
    if (!planSlug) return
    fetchPlans()
      .then((plans) => setSelectedPlan(plans.find((p) => p.slug === planSlug) ?? null))
      .catch(() => setSelectedPlan(null))
  }, [planSlug])

  const set = (field: keyof typeof form) => (
    e: React.ChangeEvent<HTMLInputElement>
  ) => setForm({ ...form, [field]: e.target.value })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    setSubmitting(true)
    setError('')
    try {
      await registerSignup({
        name: form.name,
        email: form.email,
        password: form.password,
        plan_slug: planSlug || undefined,
      })
      setDone(true)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'falha ao realizar cadastro')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <header className="border-b border-border">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center">
          <Link to="/" className="flex items-center gap-2">
            <Sprout className="h-7 w-7 text-primary" />
            <span className="text-xl font-display font-semibold text-foreground">CafeOS</span>
          </Link>
        </div>
      </header>

      <main className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="w-full max-w-xl bg-card rounded-lg shadow-lg p-6 sm:p-8">
          {done ? (
            <div className="text-center py-6">
              <CheckCircle className="h-12 w-12 text-leaf mx-auto mb-4" />
              <h1 className="font-display text-2xl font-semibold text-card-foreground mb-2">
                Cadastro realizado!
              </h1>
              <p className="text-sm text-card-foreground/60 mb-6">
                Recebemos seu cadastro. Faça login para começar a configurar sua fazenda.
              </p>
              <Link to="/">
                <Button variant="primary">Voltar ao início</Button>
              </Link>
            </div>
          ) : (
            <>
              <h1 className="font-display text-2xl font-semibold text-card-foreground mb-1">
                Cadastre-se no CafeOS
              </h1>
              <p className="text-sm text-card-foreground/60 mb-6">
                Leva menos de um minuto — conte pra gente sobre você.
              </p>

              {selectedPlan && (
                <div className="flex items-center justify-between rounded-sm border border-primary/40 bg-primary/10 px-4 py-3 mb-6">
                  <div>
                    <p className="text-xs text-muted-foreground uppercase tracking-wide">Plano selecionado</p>
                    <p className="text-sm font-medium text-card-foreground">{selectedPlan.name}</p>
                  </div>
                  <p className="font-mono text-gold font-medium">
                    {formatPrice(selectedPlan.price_cents)}
                    <span className="text-xs text-muted-foreground">/{selectedPlan.billing_interval === 'yearly' ? 'ano' : 'mês'}</span>
                  </p>
                </div>
              )}

              <form onSubmit={handleSubmit} className="space-y-4">
                <Field label="Nome" required>
                  <Input value={form.name} onChange={set('name')} required />
                </Field>
                <Field label="Email" required>
                  <Input type="email" value={form.email} onChange={set('email')} required />
                </Field>
                <Field label="Senha" required>
                  <Input type="password" value={form.password} onChange={set('password')} required minLength={6} />
                </Field>

                {error && <p className="text-sm text-destructive">{error}</p>}

                <div className="flex items-center justify-between gap-3 pt-2">
                  <RequiredLegend />
                  <Button type="submit" variant="primary" disabled={submitting}>
                    {submitting ? 'Enviando...' : 'Concluir cadastro'}
                  </Button>
                </div>
              </form>
            </>
          )}
        </div>
      </main>
    </div>
  )
}
