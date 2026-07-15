import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Building2 } from 'lucide-react'

interface Organization {
  id: string
  name: string
  brand_name: string
  slug: string
  plan: string
  plan_id: string | null
  status: string
  created_at: string
}

interface Plan {
  id: string
  name: string
  slug: string
}

export function Organizations() {
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [plans, setPlans] = useState<Plan[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Organization | null>(null)
  const [form, setForm] = useState({ name: '', plan_id: '', status: 'active' })
  const [saving, setSaving] = useState(false)
  const toast = useToast()

  const load = useCallback(async () => {
    try {
      const [organizationData, planData] = await Promise.all([
        apiRequest<Organization[]>('/admin/organizations', { admin: true }),
        apiRequest<Plan[]>('/admin/plans', { admin: true }),
      ])
      setOrganizations(organizationData)
      setPlans(planData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      if (editing) {
        await apiRequest(`/admin/organizations/${editing.id}`, { method: 'PUT', body: form, admin: true })
      } else {
        await apiRequest('/admin/organizations', { method: 'POST', body: form, admin: true })
      }
      setDialogOpen(false); setEditing(null)
      await load()
      toast.success(editing ? 'Organização atualizada' : 'Organização criada')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao salvar organização')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Organizações</h1>
          <p className="text-sm text-muted-foreground">Gerenciar clientes da plataforma</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm({ name: '', plan_id: plans[0]?.id || '', status: 'active' }); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Nova Organização
        </Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Slug</TableHeader>
            <TableHeader>Plano</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {organizations.map((o) => (
            <TableRow key={o.id}>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <Building2 className="h-4 w-4 text-primary" />
                  {o.name}
                </div>
              </TableCell>
              <TableCell className="text-muted-foreground text-sm">{o.slug}</TableCell>
              <TableCell><Badge>{plans.find((p) => p.id === o.plan_id)?.name || o.plan}</Badge></TableCell>
              <TableCell>
                <Badge variant={o.status === 'active' ? 'success' : 'default'}>
                  {o.status === 'active' ? 'Ativo' : 'Inativo'}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="ghost" size="sm" onClick={() => { setEditing(o); setForm({ name: o.name, plan_id: o.plan_id || plans[0]?.id || '', status: o.status }); setDialogOpen(true) }}>
                  <Pencil className="h-4 w-4" />
                </Button>
              </TableCell>
            </TableRow>
          ))}
          {organizations.length === 0 && (
            <TableRow><TableCell colSpan={5} className="text-center text-muted-foreground py-8">Nenhuma organização cadastrada.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Organização' : 'Nova Organização'}>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Nome</label>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Plano</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.plan_id} onChange={(e) => setForm({ ...form, plan_id: e.target.value })}>
              {plans.map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))}
            </select>
          </div>
          {editing && (
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">Status</label>
              <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })}>
                <option value="active">Ativo</option>
                <option value="inactive">Inativo</option>
              </select>
            </div>
          )}
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
