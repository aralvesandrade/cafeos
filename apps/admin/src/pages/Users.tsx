import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, User, Building2 } from 'lucide-react'

interface AppUser {
  id: string
  tenant_id: string
  name: string
  email: string
  role: string
  status: string
}

interface Tenant {
  id: string
  name: string
}

const roleLabels: Record<string, string> = {
  platform_owner: 'Platform Owner',
  tenant_admin: 'Admin',
  proprietario: 'Proprietário',
  gerente_agricola: 'Gerente Agrícola',
  engenheiro_agronomo: 'Eng. Agrônomo',
  tecnico_agricola: 'Técnico Agrícola',
  operador_campo: 'Operador de Campo',
  financeiro: 'Financeiro',
  consultor_externo: 'Consultor',
  auditor: 'Auditor',
}

export function Users() {
  const [users, setUsers] = useState<AppUser[]>([])
  const [tenants, setTenants] = useState<Tenant[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<AppUser | null>(null)
  const [form, setForm] = useState({ name: '', email: '', password: '', role: '', tenant_id: '' })
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const [userData, tenantData] = await Promise.all([
        apiRequest<AppUser[]>('/users', { admin: true }),
        apiRequest<Tenant[]>('/tenants', { admin: true }),
      ])
      setUsers(userData)
      setTenants(tenantData)
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
        await apiRequest(`/users/${editing.id}`, { method: 'PUT', body: { name: form.name, email: form.email, role: form.role }, admin: true })
      } else {
        await apiRequest('/users', { method: 'POST', body: { name: form.name, email: form.email, password: form.password, role: form.role, tenant_id: form.tenant_id }, admin: true })
      }
      setDialogOpen(false); setEditing(null)
      await load()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Remover usuário?')) return
    try {
      await apiRequest(`/users/${id}`, { method: 'DELETE', admin: true })
      await load()
    } catch (err) { console.error(err) }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-coffee-green-dark">Usuários</h1>
          <p className="text-sm text-coffee-text-light">Gerenciar usuários do sistema</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm({ name: '', email: '', password: '', role: '', tenant_id: '' }); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Usuário
        </Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Email</TableHeader>
            <TableHeader>Perfil</TableHeader>
            <TableHeader>Tenant</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((u) => {
            const tenantName = tenants.find((t) => t.id === u.tenant_id)?.name || u.tenant_id
            return (
              <TableRow key={u.id}>
                <TableCell className="font-medium">
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-coffee-green" />
                    {u.name}
                  </div>
                </TableCell>
                <TableCell>{u.email}</TableCell>
                <TableCell><Badge>{roleLabels[u.role] || u.role}</Badge></TableCell>
                <TableCell className="text-coffee-text-light text-sm">
                  <div className="flex items-center gap-1">
                    <Building2 className="h-3 w-3" />
                    {tenantName}
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant={u.status === 'active' ? 'success' : 'default'}>
                    {u.status === 'active' ? 'Ativo' : 'Inativo'}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-1">
                    <Button variant="ghost" size="sm" onClick={() => { setEditing(u); setForm({ name: u.name, email: u.email, password: '', role: u.role, tenant_id: u.tenant_id }); setDialogOpen(true) }}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => handleDelete(u.id)}>
                      <Trash2 className="h-4 w-4 text-red-500" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            )
          })}
          {users.length === 0 && (
            <TableRow><TableCell colSpan={6} className="text-center text-coffee-text-light py-8">Nenhum usuário cadastrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Usuário' : 'Novo Usuário'}>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Nome</label>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Email</label>
            <Input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} required />
          </div>
          {!editing && (
            <div>
              <label className="block text-sm font-medium text-coffee-text mb-1">Senha</label>
              <Input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} required />
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Perfil</label>
            <select className="flex h-10 w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-coffee-text" value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })}>
              {Object.entries(roleLabels).map(([value, label]) => (
                <option key={value} value={value}>{label}</option>
              ))}
            </select>
          </div>
          {!editing && (
            <div>
              <label className="block text-sm font-medium text-coffee-text mb-1">Tenant</label>
              <select className="flex h-10 w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-coffee-text" value={form.tenant_id} onChange={(e) => setForm({ ...form, tenant_id: e.target.value })}>
                <option value="">Selecione um tenant</option>
                {tenants.map((t) => (
                  <option key={t.id} value={t.id}>{t.name}</option>
                ))}
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
