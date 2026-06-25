import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Building2 } from 'lucide-react'

interface Tenant {
  id: string
  name: string
  brand_name: string
  slug: string
  plan: string
  status: string
  created_at: string
}

export function Tenants() {
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Tenant | null>(null)
  const [form, setForm] = useState({ name: '', plan: 'free' })
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<Tenant[]>('/admin/tenants', { admin: true })
      setTenants(data)
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
        await apiRequest(`/admin/tenants/${editing.id}`, { method: 'PUT', body: form, admin: true })
      } else {
        await apiRequest('/admin/tenants', { method: 'POST', body: form, admin: true })
      }
      setDialogOpen(false); setEditing(null)
      await load()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-coffee-green-dark">Tenants</h1>
          <p className="text-sm text-coffee-text-light">Gerenciar clientes da plataforma</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm({ name: '', plan: 'free' }); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Tenant
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
          {tenants.map((t) => (
            <TableRow key={t.id}>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <Building2 className="h-4 w-4 text-coffee-green" />
                  {t.name}
                </div>
              </TableCell>
              <TableCell className="text-coffee-text-light text-sm">{t.slug}</TableCell>
              <TableCell><Badge>{t.plan}</Badge></TableCell>
              <TableCell>
                <Badge variant={t.status === 'active' ? 'success' : 'default'}>
                  {t.status === 'active' ? 'Ativo' : 'Inativo'}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="ghost" size="sm" onClick={() => { setEditing(t); setForm({ name: t.name, plan: t.plan }); setDialogOpen(true) }}>
                  <Pencil className="h-4 w-4" />
                </Button>
              </TableCell>
            </TableRow>
          ))}
          {tenants.length === 0 && (
            <TableRow><TableCell colSpan={5} className="text-center text-coffee-text-light py-8">Nenhum tenant cadastrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Tenant' : 'Novo Tenant'}>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Nome</label>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Plano</label>
            <select className="flex h-10 w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-coffee-text" value={form.plan} onChange={(e) => setForm({ ...form, plan: e.target.value })}>
              <option value="free">Grátis</option>
              <option value="pro">Pro</option>
              <option value="cooperativa">Cooperativa</option>
              <option value="consultoria">Consultoria</option>
            </select>
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
