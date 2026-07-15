import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, User, Building2 } from 'lucide-react'

interface AppUser {
  id: string
  organization_id: string
  name: string
  email: string
  role: string
  status: string
}

interface Organization {
  id: string
  name: string
}

const roleLabels: Record<string, string> = {
  platform_owner: 'Platform Owner',
  organization_admin: 'Admin',
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
  const [organizations, setOrganizations] = useState<Organization[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<AppUser | null>(null)
  const [form, setForm] = useState({ name: '', email: '', password: '', role: '', organization_id: '', status: 'active' })
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const [userData, organizationData] = await Promise.all([
        apiRequest<AppUser[]>('/admin/users', { admin: true }),
        apiRequest<Organization[]>('/admin/organizations', { admin: true }),
      ])
      setUsers(userData)
      setOrganizations(organizationData)
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
        await apiRequest(`/admin/users/${editing.id}`, { method: 'PUT', body: { name: form.name, email: form.email, role: form.role, is_active: form.status === 'active' }, admin: true })
      } else {
        await apiRequest('/admin/users', { method: 'POST', body: { name: form.name, email: form.email, password: form.password, role: form.role, organization_id: form.organization_id }, admin: true })
      }
      setDialogOpen(false); setEditing(null)
      await load()
      toast.success(editing ? 'Usuário atualizado' : 'Usuário criado')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao salvar usuário')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover usuário?', variant: 'danger' }))) return
    try {
      await apiRequest(`/admin/users/${id}`, { method: 'DELETE', admin: true })
      await load()
      toast.success('Usuário removido')
    } catch (err) { console.error(err); toast.error('Erro ao remover usuário') }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-primary">Usuários</h1>
          <p className="text-sm text-muted-foreground">Gerenciar usuários do sistema</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm({ name: '', email: '', password: '', role: '', organization_id: '', status: 'active' }); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Usuário
        </Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Email</TableHeader>
            <TableHeader>Perfil</TableHeader>
            <TableHeader>Organização</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((u) => {
            const organizationName = organizations.find((o) => o.id === u.organization_id)?.name || u.organization_id
            return (
              <TableRow key={u.id}>
                <TableCell className="font-medium">
                  <div className="flex items-center gap-2">
                    <User className="h-4 w-4 text-primary" />
                    {u.name}
                  </div>
                </TableCell>
                <TableCell>{u.email}</TableCell>
                <TableCell><Badge>{roleLabels[u.role] || u.role}</Badge></TableCell>
                <TableCell className="text-muted-foreground text-sm">
                  <div className="flex items-center gap-1">
                    <Building2 className="h-3 w-3" />
                    {organizationName}
                  </div>
                </TableCell>
                <TableCell>
                  <Badge variant={u.status === 'active' ? 'success' : 'default'}>
                    {u.status === 'active' ? 'Ativo' : 'Inativo'}
                  </Badge>
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-1">
                    <Button variant="ghost" size="sm" onClick={() => { setEditing(u); setForm({ name: u.name, email: u.email, password: '', role: u.role, organization_id: u.organization_id, status: u.status }); setDialogOpen(true) }}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => handleDelete(u.id)}>
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            )
          })}
          {users.length === 0 && (
            <TableRow><TableCell colSpan={6} className="text-center text-muted-foreground py-8">Nenhum usuário cadastrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Usuário' : 'Novo Usuário'}>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Nome</label>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Email</label>
            <Input type="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} required />
          </div>
          {!editing && (
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">Senha</label>
              <Input type="password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} required />
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Perfil</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })}>
              {Object.entries(roleLabels).map(([value, label]) => (
                <option key={value} value={value}>{label}</option>
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
          {!editing && (
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">Organização</label>
              <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.organization_id} onChange={(e) => setForm({ ...form, organization_id: e.target.value })}>
                <option value="">Selecione uma organização</option>
                {organizations.map((o) => (
                  <option key={o.id} value={o.id}>{o.name}</option>
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
