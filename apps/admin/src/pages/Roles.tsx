import { useState } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { useRoles } from '@/lib/roles'
import { useModuleAccess } from '@/lib/permissions'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, Shield } from 'lucide-react'

export function Roles() {
  const { roles, loading, refresh } = useRoles()
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<{ id: string; key: string; name: string } | null>(null)
  const [form, setForm] = useState({ key: '', name: '' })
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()
  const canEdit = useModuleAccess('permissions') === 'write'

  const handleSave = async () => {
    setSaving(true)
    try {
      if (editing) {
        await apiRequest(`/roles/${editing.id}`, { method: 'PUT', body: { name: form.name } })
      } else {
        await apiRequest('/roles', { method: 'POST', body: { key: form.key, name: form.name } })
      }
      setDialogOpen(false); setEditing(null)
      await refresh()
      toast.success(editing ? 'Papel atualizado' : 'Papel criado')
    } catch (err) {
      console.error(err)
      toast.error(editing ? 'Erro ao atualizar papel' : 'Erro ao criar papel')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover papel?', description: 'Só é possível remover papéis que não estejam atribuídos a nenhum usuário.', variant: 'danger' }))) return
    try {
      await apiRequest(`/roles/${id}`, { method: 'DELETE' })
      await refresh()
      toast.success('Papel removido')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao remover papel — verifique se ainda está em uso')
    }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Papéis</h1>
          <p className="text-sm text-muted-foreground">
            Papéis de sistema (Platform Owner, Admin) não podem ser alterados. Os demais são específicos desta organização.
          </p>
        </div>
        {canEdit && (
          <Button onClick={() => { setEditing(null); setForm({ key: '', name: '' }); setDialogOpen(true) }}>
            <Plus className="h-4 w-4" /> Novo Papel
          </Button>
        )}
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Chave</TableHeader>
            <TableHeader>Tipo</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {roles.map((role) => (
            <TableRow key={role.id}>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4 text-primary" />
                  {role.name}
                </div>
              </TableCell>
              <TableCell className="text-muted-foreground text-sm">{role.key}</TableCell>
              <TableCell>
                <Badge variant={role.is_system ? 'default' : 'success'}>
                  {role.is_system ? 'Sistema' : 'Organização'}
                </Badge>
              </TableCell>
              <TableCell className="text-right">
                {canEdit && !role.is_system && (
                  <div className="flex justify-end gap-1">
                    <Button variant="ghost" size="sm" onClick={() => { setEditing({ id: role.id, key: role.key, name: role.name }); setForm({ key: role.key, name: role.name }); setDialogOpen(true) }}>
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => handleDelete(role.id)}>
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </div>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Papel' : 'Novo Papel'}>
        <div className="space-y-4">
          {!editing && (
            <div>
              <label className="block text-sm font-medium text-foreground mb-1">Chave</label>
              <Input
                value={form.key}
                onChange={(e) => setForm({ ...form, key: e.target.value.trim().toLowerCase().replace(/\s+/g, '_') })}
                placeholder="ex: colhedor_chefe"
                required
              />
            </div>
          )}
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Nome</label>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="ex: Colhedor Chefe" required />
          </div>
          <Button className="w-full" onClick={handleSave} disabled={saving || !form.name || (!editing && !form.key)}>
            {saving ? 'Salvando...' : 'Salvar'}
          </Button>
        </div>
      </Dialog>
    </div>
  )
}
