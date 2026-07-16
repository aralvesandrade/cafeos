import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Field } from '@/components/ui/field'
import { RequiredLegend } from '@/components/ui/required-legend'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, CreditCard } from 'lucide-react'

interface PlanFeature {
  label: string
  included: boolean
}

interface Plan {
  id: string
  name: string
  slug: string
  description: string
  price_cents: number
  billing_interval: string
  max_farms: number
  max_users: number
  features: PlanFeature[]
  active: boolean
  featured: boolean
  display_order: number
}

// Uma feature por linha; prefixo "-" marca item não incluído (aparece
// riscado/com X na landing page). Ex.: "Relatórios CSV/PDF" (incluído),
// "-White Label" (não incluído).
function parseFeatures(text: string): PlanFeature[] {
  return text
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      if (line.startsWith('-')) {
        return { label: line.slice(1).trim(), included: false }
      }
      return { label: line, included: true }
    })
}

function featuresToText(features: PlanFeature[]): string {
  return (features || []).map((f) => (f.included ? f.label : `-${f.label}`)).join('\n')
}

const emptyForm = {
  name: '',
  slug: '',
  description: '',
  price: '0',
  billing_interval: 'monthly',
  max_farms: '0',
  max_users: '0',
  features: '',
  active: 'true',
  featured: 'false',
  display_order: '0',
}

function formatPrice(cents: number) {
  return (cents / 100).toLocaleString('pt-BR', { style: 'currency', currency: 'BRL' })
}

export function Plans() {
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Plan | null>(null)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<Plan[]>('/admin/plans', { admin: true })
      setPlans(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const buildBody = () => ({
    name: form.name,
    slug: form.slug,
    description: form.description,
    price_cents: Math.round(parseFloat(form.price.replace(',', '.') || '0') * 100),
    billing_interval: form.billing_interval,
    max_farms: parseInt(form.max_farms, 10) || 0,
    max_users: parseInt(form.max_users, 10) || 0,
    features: parseFeatures(form.features),
    active: form.active === 'true',
    featured: form.featured === 'true',
    display_order: parseInt(form.display_order, 10) || 0,
  })

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      const body = buildBody()
      if (editing) {
        await apiRequest(`/admin/plans/${editing.id}`, { method: 'PUT', body, admin: true })
      } else {
        await apiRequest('/admin/plans', { method: 'POST', body, admin: true })
      }
      setDialogOpen(false); setEditing(null)
      await load()
      toast.success(editing ? 'Plano atualizado' : 'Plano criado')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao salvar plano')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover plano?', variant: 'danger' }))) return
    try {
      await apiRequest(`/admin/plans/${id}`, { method: 'DELETE', admin: true })
      await load()
      toast.success('Plano removido')
    } catch (err) { console.error(err); toast.error('Erro ao remover plano') }
  }

  const openEdit = (p: Plan) => {
    setEditing(p)
    setForm({
      name: p.name,
      slug: p.slug,
      description: p.description,
      price: (p.price_cents / 100).toString(),
      billing_interval: p.billing_interval,
      max_farms: p.max_farms.toString(),
      max_users: p.max_users.toString(),
      features: featuresToText(p.features),
      active: p.active ? 'true' : 'false',
      featured: p.featured ? 'true' : 'false',
      display_order: p.display_order.toString(),
    })
    setDialogOpen(true)
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Planos</h1>
          <p className="text-sm text-muted-foreground">Gerenciar planos de assinatura da plataforma</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm(emptyForm); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Plano
        </Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Slug</TableHeader>
            <TableHeader>Preço</TableHeader>
            <TableHeader>Intervalo</TableHeader>
            <TableHeader>Limites</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {plans.map((p) => (
            <TableRow key={p.id}>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <CreditCard className="h-4 w-4 text-primary" />
                  {p.name}
                </div>
              </TableCell>
              <TableCell className="text-muted-foreground text-sm">{p.slug}</TableCell>
              <TableCell>{formatPrice(p.price_cents)}</TableCell>
              <TableCell className="text-muted-foreground text-sm">{p.billing_interval === 'yearly' ? 'Anual' : 'Mensal'}</TableCell>
              <TableCell className="text-muted-foreground text-sm">
                {p.max_farms === 0 ? '∞' : p.max_farms} fazendas / {p.max_users === 0 ? '∞' : p.max_users} usuários
              </TableCell>
              <TableCell>
                <div className="flex gap-1">
                  <Badge variant={p.active ? 'success' : 'default'}>{p.active ? 'Ativo' : 'Inativo'}</Badge>
                  {p.featured && <Badge variant="success">Destaque</Badge>}
                </div>
              </TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Button variant="ghost" size="sm" onClick={() => openEdit(p)}>
                    <Pencil className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDelete(p.id)}>
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
          {plans.length === 0 && (
            <TableRow><TableCell colSpan={7} className="text-center text-muted-foreground py-8">Nenhum plano cadastrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Plano' : 'Novo Plano'}>
        <form onSubmit={handleSave} className="space-y-4">
          <Field label="Nome" required>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </Field>
          <Field label="Slug" required>
            <Input value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} placeholder="ex: pro" required />
          </Field>
          <Field label="Descrição">
            <Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
          </Field>
          <div className="grid grid-cols-2 gap-4">
            <Field label="Preço (R$)">
              <Input type="number" step="0.01" min="0" value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} />
            </Field>
            <Field label="Intervalo">
              <Select value={form.billing_interval} onChange={(e) => setForm({ ...form, billing_interval: e.target.value })}>
                <option value="monthly">Mensal</option>
                <option value="yearly">Anual</option>
              </Select>
            </Field>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <Field label="Máx. fazendas (0 = ilimitado)">
              <Input type="number" min="0" value={form.max_farms} onChange={(e) => setForm({ ...form, max_farms: e.target.value })} />
            </Field>
            <Field label="Máx. usuários (0 = ilimitado)">
              <Input type="number" min="0" value={form.max_users} onChange={(e) => setForm({ ...form, max_users: e.target.value })} />
            </Field>
          </div>
          <Field label='Features (uma por linha; prefixe com "-" para marcar como não incluída)'>
            <textarea
              className="flex min-h-[120px] w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/20 focus-visible:border-ring"
              value={form.features}
              onChange={(e) => setForm({ ...form, features: e.target.value })}
              placeholder={'Relatórios CSV/PDF\n-White Label\nSuporte prioritário'}
            />
          </Field>
          <div className="grid grid-cols-2 gap-4">
            <Field label="Status">
              <Select value={form.active} onChange={(e) => setForm({ ...form, active: e.target.value })}>
                <option value="true">Ativo</option>
                <option value="false">Inativo</option>
              </Select>
            </Field>
            <Field label="Ordem de exibição">
              <Input type="number" value={form.display_order} onChange={(e) => setForm({ ...form, display_order: e.target.value })} />
            </Field>
          </div>
          <Field label='Destaque ("Mais Popular" na landing page)'>
            <Select value={form.featured} onChange={(e) => setForm({ ...form, featured: e.target.value })}>
              <option value="false">Não</option>
              <option value="true">Sim</option>
            </Select>
          </Field>
          <div className="flex items-center justify-between gap-3 pt-2">
            <RequiredLegend />
            <div className="flex gap-3">
              <Button type="button" variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
              <Button type="submit" disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
            </div>
          </div>
        </form>
      </Dialog>
    </div>
  )
}
