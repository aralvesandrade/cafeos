import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, Tractor } from 'lucide-react'

interface OperationType {
  id: string
  name: string
  code: string
  color: 'info' | 'warning' | 'success' | 'default' | 'danger'
}

const colorLabels: Record<string, string> = {
  info: 'Azul',
  warning: 'Amarelo',
  success: 'Verde',
  default: 'Cinza',
  danger: 'Vermelho',
}

const emptyForm = { name: '', code: '', color: 'default' }

export function OperationTypes() {
  const [items, setItems] = useState<OperationType[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<OperationType | null>(null)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<OperationType[]>('/operation-types')
      setItems(data)
    } catch (err) { console.error(err) }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      if (editing) await apiRequest(`/operation-types/${editing.id}`, { method: 'PUT', body: form })
      else await apiRequest('/operation-types', { method: 'POST', body: form })
      setDialogOpen(false); setEditing(null); await load()
      toast.success(editing ? 'Tipo de operação atualizado' : 'Tipo de operação criado')
    } catch (err) { console.error(err); toast.error('Erro ao salvar tipo de operação') }
    finally { setSaving(false) }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover tipo de operação?', variant: 'danger' }))) return
    try {
      await apiRequest(`/operation-types/${id}`, { method: 'DELETE' }); await load()
      toast.success('Tipo de operação removido')
    }
    catch (err) { console.error(err); toast.error('Erro ao remover tipo de operação') }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between">
      <div><h1 className="text-2xl font-bold text-primary">Tipos de Operação</h1><p className="text-sm text-muted-foreground">Cadastro de tipos de operações agrícolas</p></div>
      <Button onClick={() => { setEditing(null); setForm(emptyForm); setDialogOpen(true) }}><Plus className="h-4 w-4" /> Novo Tipo</Button>
    </div>
    <Table>
      <TableHead><TableRow><TableHeader>Código</TableHeader><TableHeader>Nome</TableHeader><TableHeader>Cor</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{items.map((ot) => (<TableRow key={ot.id}>
        <TableCell className="font-mono text-sm">{ot.code}</TableCell>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><Tractor className="h-4 w-4 text-primary" />{ot.name}</div></TableCell>
        <TableCell><Badge variant={ot.color || 'default'}>{colorLabels[ot.color] || ot.color}</Badge></TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditing(ot); setForm({ name: ot.name, code: ot.code, color: ot.color }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => handleDelete(ot.id)}><Trash2 className="h-4 w-4 text-destructive" /></Button>
        </div></TableCell>
      </TableRow>))}
      {items.length === 0 && <TableRow><TableCell colSpan={4} className="text-center text-muted-foreground py-8">Nenhum tipo de operação cadastrado.</TableCell></TableRow>}
      </TableBody>
    </Table>
    <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Tipo de Operação' : 'Novo Tipo de Operação'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Nome</label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Código</label><Input value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="EX: adubacao" required /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Cor</label>
          <Select value={form.color} onChange={(e) => setForm({ ...form, color: e.target.value })}>
            {Object.entries(colorLabels).map(([value, label]) => <option key={value} value={value}>{label}</option>)}
          </Select></div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
