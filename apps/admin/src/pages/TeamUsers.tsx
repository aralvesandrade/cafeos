import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { useModuleAccess } from '@/lib/permissions'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, User } from 'lucide-react'
import { ROLE_LABELS, ALL_ROLES } from '@/lib/roles'

interface OrgUser {
  id: string
  name: string
  email: string
  role: string
  status: string
}

export function TeamUsers() {
  const [users, setUsers] = useState<OrgUser[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<OrgUser | null>(null)
  const [form, setForm] = useState({ name: '', email: '', password: '', role: '', status: 'active' })
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()
  const canEdit = useModuleAccess('users') === 'write'

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<OrgUser[]>('/users')
      setUsers(data)
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
        await apiRequest(`/users/${editing.id}`, { method: 'PUT', body: { name: form.name, email: form.email, role: form.role, is_active: form.status === 'active' } })
      } else {
        await apiRequest('/users', { method: 'POST', body: { name: form.name, email: form.email, password: form.password, role: form.role } })
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
      await apiRequest(`/users/${id}`, { method: 'DELETE' })
      await load()
      toast.success('Usuário removido')
    } catch (err) { console.error(err); toast.error('Erro ao remover usuário') }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Minha Equipe</h1>
          <p className="text-sm text-muted-foreground">Usuários da sua organização</p>
        </div>
        {canEdit && (
          <Button onClick={() => { setEditing(null); setForm({ name: '', email: '', password: '', role: '', status: 'active' }); setDialogOpen(true) }}>
            <Plus className="h-4 w-4" /> Novo Usuário
          </Button>
        )}
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Email</TableHeader>
            <TableHeader>Perfil</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((u) => (
            <TableRow key={u.id}>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <User className="h-4 w-4 text-primary" />
                  {u.name}
                </div>
              </TableCell>
              <TableCell>{u.email}</TableCell>
              <TableCell><Badge>{ROLE_LABELS[u.role as keyof typeof ROLE_LABELS] || u.role}</Badge></TableCell>
              <TableCell>
                <Badge variant={u.status === 'active' ? 'success' : 'default'}>
                  {u.status === 'active' ? 'Ativo' : 'Inativo'}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                {canEdit && (
                  <div className="flex justify-end gap-1">
                    <Button variant="ghost" size="sm" onClick={() => { setEditing(u); setForm({ name: u.name, email: u.email, password: '', role: u.role, status: u.status }); setDialogOpen(true) }}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => handleDelete(u.id)}>
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                )}
              </TableCell>
            </TableRow>
          ))}
          {users.length === 0 && (
            <TableRow><TableCell colSpan={5} className="text-center text-muted-foreground py-8">Nenhum usuário cadastrado.</TableCell></TableRow>
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
              <option value="">Selecione...</option>
              {ALL_ROLES.map((role) => (
                <option key={role} value={role}>{ROLE_LABELS[role]}</option>
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
          <Button className="w-full" onClick={handleSave} disabled={saving || !form.name || !form.email || (!editing && !form.password)}>
            {saving ? 'Salvando...' : 'Salvar'}
          </Button>
        </div>
      </Dialog>
    </div>
  )
}
