import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Field } from '@/components/ui/field'
import { RequiredLegend } from '@/components/ui/required-legend'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Puzzle } from 'lucide-react'
import type { ModuleMeta } from '@/lib/permissions'

const emptyForm = { key: '', name: '', order: '0' }

export function Modules() {
  const [modules, setModules] = useState<ModuleMeta[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<ModuleMeta | null>(null)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)
  const toast = useToast()

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<ModuleMeta[]>('/modules')
      data.sort((a, b) => a.order - b.order)
      setModules(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      if (editing) {
        await apiRequest(`/admin/modules/${editing.key}`, { method: 'PUT', body: { name: form.name, order: Number(form.order) }, admin: true })
      } else {
        await apiRequest('/admin/modules', { method: 'POST', body: { key: form.key, name: form.name, order: Number(form.order) }, admin: true })
      }
      setDialogOpen(false); setEditing(null)
      await load()
      toast.success(editing ? 'Módulo atualizado' : 'Módulo criado')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao salvar módulo')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Módulos</h1>
          <p className="text-sm text-muted-foreground">Catálogo de módulos do sistema</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm(emptyForm); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Módulo
        </Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Ordem</TableHeader>
            <TableHeader>Chave</TableHeader>
            <TableHeader>Nome</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {modules.map((m) => (
            <TableRow key={m.key}>
              <TableCell className="text-muted-foreground text-sm">{m.order}</TableCell>
              <TableCell><Badge variant="default">{m.key}</Badge></TableCell>
              <TableCell className="font-medium">
                <div className="flex items-center gap-2">
                  <Puzzle className="h-4 w-4 text-primary" />
                  {m.name}
                </div>
              </TableCell>
              <TableCell className="text-right">
                <Button variant="ghost" size="sm" onClick={() => { setEditing(m); setForm({ key: m.key, name: m.name, order: String(m.order) }); setDialogOpen(true) }}>
                  <Pencil className="h-4 w-4" />
                </Button>
              </TableCell>
            </TableRow>
          ))}
          {modules.length === 0 && (
            <TableRow><TableCell colSpan={4} className="text-center text-muted-foreground py-8">Nenhum módulo encontrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Módulo' : 'Novo Módulo'}>
        <form onSubmit={handleSave} className="space-y-4">
          <Field label="Chave" required>
            <Input value={form.key} onChange={(e) => setForm({ ...form, key: e.target.value })} required disabled={!!editing} placeholder="ex: custom_module" />
          </Field>
          <Field label="Nome" required>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </Field>
          <Field label="Ordem">
            <Input type="number" value={form.order} onChange={(e) => setForm({ ...form, order: e.target.value })} />
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
