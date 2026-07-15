import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { Plus, Pencil, Trash2, DollarSign } from 'lucide-react'

interface CostCenter {
  id: string
  name: string
  code: string
  type: string
}

interface Farm { id: string; name: string }

interface Transaction {
  id: string
  type: string
  cost_center_id: string | null
  farm_id: string | null
  description: string
  amount: number
  date: string
  due_date: string
  status: string
  notes: string
}

const typeLabels: Record<string, string> = { receita: 'Receita', despesa: 'Despesa' }
const typeVariants: Record<string, 'success' | 'danger'> = { receita: 'success', despesa: 'danger' }
const statusLabels: Record<string, string> = { pending: 'Pendente', paid: 'Pago', cancelled: 'Cancelado' }

export function Financial() {
  const [items, setItems] = useState<Transaction[]>([])
  const [costCenters, setCostCenters] = useState<CostCenter[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [farmFilter, setFarmFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Transaction | null>(null)
  const [form, setForm] = useState({ type: 'despesa', cost_center_id: '', farm_id: '', description: '', amount: '', date: '', due_date: '', notes: '', status: 'pending' })
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const [data, ccs, farmsData] = await Promise.all([
        apiRequest<Transaction[]>('/financial', { params: farmFilter ? { farm_id: farmFilter } : undefined }),
        apiRequest<CostCenter[]>('/cost-centers'),
        apiRequest<Farm[]>('/farms'),
      ])
      setItems(data)
      setCostCenters(ccs)
      setFarms(farmsData)
    }
    catch (err) { console.error(err) }
    finally { setLoading(false) }
  }, [farmFilter])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      const body = { ...form, cost_center_id: form.cost_center_id || null, farm_id: form.farm_id || null, amount: parseFloat(form.amount) }
      if (editing) await apiRequest(`/financial/${editing.id}`, { method: 'PUT', body })
      else await apiRequest('/financial', { method: 'POST', body })
      setDialogOpen(false); setEditing(null); await load()
      toast.success(editing ? 'Transação atualizada' : 'Transação registrada')
    } catch (err) { console.error(err); toast.error('Erro ao salvar transação') }
    finally { setSaving(false) }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover transação?', variant: 'danger' }))) return
    try {
      await apiRequest(`/financial/${id}`, { method: 'DELETE' }); await load()
      toast.success('Transação removida')
    }
    catch (err) { console.error(err); toast.error('Erro ao remover transação') }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div><h1 className="font-display text-2xl font-semibold text-foreground">Financeiro</h1><p className="text-sm text-muted-foreground">Contas a pagar e receber</p></div>
      <div className="flex gap-3">
        <Select value={farmFilter} onChange={(e) => setFarmFilter(e.target.value)} className="w-48">
          <option value="">Todas as fazendas</option>
          {farms.map((f) => (
            <option key={f.id} value={f.id}>{f.name}</option>
          ))}
        </Select>
        <Button onClick={() => { setEditing(null); setForm({ type: 'despesa', cost_center_id: '', farm_id: '', description: '', amount: '', date: '', due_date: '', notes: '', status: 'pending' }); setDialogOpen(true) }}><Plus className="h-4 w-4" /> Nova Transação</Button>
      </div>
    </div>
    <Table>
      <TableHead><TableRow><TableHeader>Tipo</TableHeader><TableHeader>Descrição</TableHeader><TableHeader>Centro de Custo</TableHeader><TableHeader>Fazenda</TableHeader><TableHeader>Valor</TableHeader><TableHeader>Data</TableHeader><TableHeader>Status</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{items.map((t) => {
        const cc = costCenters.find(c => c.id === t.cost_center_id)
        const farm = farms.find(f => f.id === t.farm_id)
        return (<TableRow key={t.id}>
        <TableCell><Badge variant={typeVariants[t.type]}>{typeLabels[t.type] || t.type}</Badge></TableCell>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><DollarSign className="h-4 w-4 text-primary" />{t.description}</div></TableCell>
        <TableCell className="text-muted-foreground text-sm">{cc?.name || '-'}</TableCell>
        <TableCell className="text-muted-foreground text-sm">{farm?.name || '-'}</TableCell>
        <TableCell>R$ {t.amount.toFixed(2)}</TableCell>
        <TableCell>{t.date}</TableCell>
        <TableCell><Badge variant={t.status === 'paid' ? 'success' : t.status === 'cancelled' ? 'danger' : 'default'}>{statusLabels[t.status] || t.status}</Badge></TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditing(t); setForm({ type: t.type, cost_center_id: t.cost_center_id || '', farm_id: t.farm_id || '', description: t.description, amount: String(t.amount), date: t.date, due_date: t.due_date, notes: t.notes, status: t.status }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => handleDelete(t.id)}><Trash2 className="h-4 w-4 text-destructive" /></Button>
        </div></TableCell>
      </TableRow>)})}
      {items.length === 0 && <TableRow><TableCell colSpan={8} className="text-center text-muted-foreground py-8">Nenhuma transação cadastrada.</TableCell></TableRow>}
      </TableBody>
    </Table>
    <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Transação' : 'Nova Transação'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Tipo</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
            <option value="despesa">Despesa</option><option value="receita">Receita</option>
          </select></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Descrição</label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Centro de Custo</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.cost_center_id} onChange={(e) => setForm({ ...form, cost_center_id: e.target.value })}>
            <option value="">Selecione</option>
            {costCenters.filter(cc => cc.type === form.type).map(cc => (
              <option key={cc.id} value={cc.id}>{cc.code} — {cc.name}</option>
            ))}
          </select></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Fazenda</label>
          <Select value={form.farm_id} onChange={(e) => setForm({ ...form, farm_id: e.target.value })}>
            <option value="">Sem fazenda vinculada</option>
            {farms.map(f => (
              <option key={f.id} value={f.id}>{f.name}</option>
            ))}
          </Select></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Valor (R$)</label><Input type="number" step="0.01" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Data</label><Input type="date" value={form.date} onChange={(e) => setForm({ ...form, date: e.target.value })} /></div>
        {editing && <div><label className="block text-sm font-medium text-foreground mb-1">Status</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm text-foreground" value={form.status} onChange={(e) => setForm({ ...form, status: e.target.value })}>
            <option value="pending">Pendente</option><option value="paid">Pago</option><option value="cancelled">Cancelado</option>
          </select></div>}
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
