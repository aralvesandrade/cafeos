import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Sprout, CheckCircle, ArrowRight } from 'lucide-react'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Tabs } from '@/components/ui/tabs'
import { RequiredLegend } from '@/components/ui/required-legend'
import { Button } from '@/components/ui/button'
import { registerSignup, fetchPlans, type Plan } from '@/lib/api'

function formatPrice(cents: number) {
  return (cents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL', minimumFractionDigits: cents % 100 === 0 ? 0 : 2 })
}

const ESTADOS = [
  'AC', 'AL', 'AP', 'AM', 'BA', 'CE', 'DF', 'ES', 'GO', 'MA', 'MT', 'MS',
  'MG', 'PA', 'PB', 'PR', 'PE', 'PI', 'RJ', 'RN', 'RS', 'RO', 'RR', 'SC',
  'SP', 'SE', 'TO',
]

const steps = [
  { key: 'dono', label: '1. Seus dados' },
  { key: 'fazenda', label: '2. Sua fazenda' },
]

export function SignupPage() {
  const [searchParams] = useSearchParams()
  const planSlug = searchParams.get('plano') || ''
  const [selectedPlan, setSelectedPlan] = useState<Plan | null>(null)
  const [activeTab, setActiveTab] = useState('dono')
  const [form, setForm] = useState({
    name: '',
    email: '',
    password: '',
    phone: '',
    farm_name: '',
    state: '',
    city: '',
    main_crop: '',
    total_area_ha: '',
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
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => setForm({ ...form, [field]: e.target.value })

  const handleAvancar = () => {
    if (!form.name || !form.email || !form.password) return
    setActiveTab('fazenda')
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (activeTab !== 'fazenda') {
      handleAvancar()
      return
    }
    if (!form.farm_name) return

    setSubmitting(true)
    setError('')
    try {
      await registerSignup({
        name: form.name,
        email: form.email,
        password: form.password,
        phone: form.phone,
        farm_name: form.farm_name,
        state: form.state,
        city: form.city,
        main_crop: form.main_crop,
        total_area_ha: form.total_area_ha ? Number(form.total_area_ha) : 0,
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
                Recebemos os dados da sua fazenda. Nossa equipe pode entrar em
                contato pelo telefone informado para ajudar você a começar.
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
                Leva menos de dois minutos — conte pra gente sobre você e sua fazenda.
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
                <Tabs tabs={steps} activeTab={activeTab}>
                  {(tab) =>
                    tab === 'dono' ? (
                      <div className="space-y-4">
                        <Field label="Nome" required>
                          <Input value={form.name} onChange={set('name')} required />
                        </Field>
                        <Field label="Email" required>
                          <Input type="email" value={form.email} onChange={set('email')} required />
                        </Field>
                        <Field label="Senha" required>
                          <Input type="password" value={form.password} onChange={set('password')} required minLength={6} />
                        </Field>
                        <Field label="Telefone">
                          <Input value={form.phone} onChange={set('phone')} placeholder="(DDD) 99999-9999" />
                        </Field>
                      </div>
                    ) : (
                      <div className="space-y-4">
                        <Field label="Nome da fazenda" required>
                          <Input value={form.farm_name} onChange={set('farm_name')} required />
                        </Field>
                        <div className="grid grid-cols-2 gap-4">
                          <Field label="Estado">
                            <Select value={form.state} onChange={set('state')}>
                              <option value="">Selecione</option>
                              {ESTADOS.map((uf) => (
                                <option key={uf} value={uf}>{uf}</option>
                              ))}
                            </Select>
                          </Field>
                          <Field label="Cidade">
                            <Input value={form.city} onChange={set('city')} />
                          </Field>
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                          <Field label="Área total (ha)">
                            <Input type="number" min="0" step="0.01" value={form.total_area_ha} onChange={set('total_area_ha')} />
                          </Field>
                          <Field label="Cultura principal">
                            <Input value={form.main_crop} onChange={set('main_crop')} placeholder="Café arábica" />
                          </Field>
                        </div>
                      </div>
                    )
                  }
                </Tabs>

                {error && <p className="text-sm text-destructive">{error}</p>}

                <div className="flex items-center justify-between gap-3 pt-2">
                  <RequiredLegend />
                  <div className="flex gap-3">
                    {activeTab === 'fazenda' && (
                      <Button type="button" variant="outline" onClick={() => setActiveTab('dono')}>
                        Voltar
                      </Button>
                    )}
                    {activeTab === 'dono' ? (
                      <Button type="submit" variant="primary" className="gap-2">
                        Avançar
                        <ArrowRight className="h-4 w-4" />
                      </Button>
                    ) : (
                      <Button type="submit" variant="primary" disabled={submitting}>
                        {submitting ? 'Enviando...' : 'Concluir cadastro'}
                      </Button>
                    )}
                  </div>
                </div>
              </form>
            </>
          )}
        </div>
      </main>
    </div>
  )
}
