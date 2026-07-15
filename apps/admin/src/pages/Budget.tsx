import { useEffect, useState, useCallback } from 'react'
import { useSearchParams } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, Wallet } from 'lucide-react'

interface BudgetItem {
  id: string
  harvest_id: string
  cost_center_id: string
  cost_center_name: string
  planned_amount: number
  realized_amount: number
  variance: number
  description: string
}

interface Harvest { id: string; description: string }
interface CostCenter { id: string; name: string; code: string; type: 'receita' | 'despesa' }

const emptyForm = { harvest_id: '', cost_center_id: '', planned_amount: '', description: '' }

export function Budget() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [budgets, setBudgets] = useState<BudgetItem[]>([])
  const [harvests, setHarvests] = useState<Harvest[]>([])
  const [costCenters, setCostCenters] = useState<CostCenter[]>([])
  const [harvestId, setHarvestId] = useState(searchParams.get('harvest_id') || '')
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<BudgetItem | null>(null)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const [harvestsData, ccData] = await Promise.all([
        apiRequest<Harvest[]>('/harvests'),
        apiRequest<CostCenter[]>('/cost-centers'),
      ])
      setHarvests(harvestsData)
      setCostCenters(ccData)
      const hid = harvestId || harvestsData[0]?.id || ''
      if (hid && !harvestId) setHarvestId(hid)
      if (hid) {
        const budgetsData = await apiRequest<BudgetItem[]>(`/harvests/${hid}/budgets`)
        setBudgets(budgetsData)
      }
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [harvestId])

  useEffect(() => { load() }, [load])

  useEffect(() => {
    if (searchParams.get('harvest_id')) {
      setSearchParams((params) => { params.delete('harvest_id'); return params }, { replace: true })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const handleSave = async () => {
    setSaving(true)
    try {
      const body = { ...form, planned_amount: parseFloat(form.planned_amount) || 0 }
      if (editing) await apiRequest(`/budgets/${editing.id}`, { method: 'PUT', body: { planned_amount: body.planned_amount, description: body.description } })
      else await apiRequest('/budgets', { method: 'POST', body })
      setDialogOpen(false); setEditing(null)
      await load()
      toast.success(editing ? 'Orçamento atualizado' : 'Orçamento criado')
    } catch (err) { console.error(err); toast.error('Erro ao salvar orçamento') }
    finally { setSaving(false) }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover orçamento?', variant: 'danger' }))) return
    try { await apiRequest(`/budgets/${id}`, { method: 'DELETE' }); await load(); toast.success('Orçamento removido') }
    catch (err) { console.error(err); toast.error('Erro ao remover orçamento') }
  }

  const despesaCostCenters = costCenters.filter((cc) => cc.type === 'despesa')

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-primary">Orçamento</h1>
          <p className="text-sm text-muted-foreground">Orçado x realizado por centro de custo</p>
        </div>
        <Button onClick={() => { setEditing(null); setForm({ ...emptyForm, harvest_id: harvestId }); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Orçamento
        </Button>
      </div>

      <Select value={harvestId} onChange={(e) => setHarvestId(e.target.value)} className="w-64">
        {harvests.map((h) => (
          <option key={h.id} value={h.id}>{h.description}</option>
        ))}
      </Select>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Centro de Custo</TableHeader>
            <TableHeader className="text-right">Orçado</TableHeader>
            <TableHeader className="text-right">Realizado</TableHeader>
            <TableHeader className="text-right">Variação</TableHeader>
            <TableHeader className="text-right">Execução</TableHeader>
            <TableHeader>Descrição</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {budgets.map((b) => {
            const pct = b.planned_amount > 0 ? (b.realized_amount / b.planned_amount) * 100 : 0
            return (
              <TableRow key={b.id}>
                <TableCell className="font-medium"><div className="flex items-center gap-2"><Wallet className="h-4 w-4 text-primary" />{b.cost_center_name}</div></TableCell>
                <TableCell className="text-right">R$ {b.planned_amount.toFixed(2)}</TableCell>
                <TableCell className="text-right">R$ {b.realized_amount.toFixed(2)}</TableCell>
                <TableCell className="text-right"><Badge variant={b.variance >= 0 ? 'success' : 'danger'}>R$ {b.variance.toFixed(2)}</Badge></TableCell>
                <TableCell className="text-right">{pct.toFixed(0)}%</TableCell>
                <TableCell className="text-muted-foreground text-sm">{b.description || '-'}</TableCell>
                <TableCell className="text-right"><div className="flex justify-end gap-1">
                  <Button variant="ghost" size="sm" onClick={() => { setEditing(b); setForm({ harvest_id: b.harvest_id, cost_center_id: b.cost_center_id, planned_amount: String(b.planned_amount), description: b.description }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDelete(b.id)}><Trash2 className="h-4 w-4 text-destructive" /></Button>
                </div></TableCell>
              </TableRow>
            )
          })}
          {budgets.length === 0 && <TableRow><TableCell colSpan={7} className="text-center text-muted-foreground py-8">Nenhum orçamento cadastrado para esta safra.</TableCell></TableRow>}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Orçamento' : 'Novo Orçamento'}>
        <div className="space-y-4">
          {!editing && (<>
            <div><label className="block text-sm font-medium text-foreground mb-1">Safra</label>
              <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={form.harvest_id} onChange={(e) => setForm({ ...form, harvest_id: e.target.value })}>
                <option value="">Selecione</option>
                {harvests.map((h) => <option key={h.id} value={h.id}>{h.description}</option>)}
              </select></div>
            <div><label className="block text-sm font-medium text-foreground mb-1">Centro de Custo</label>
              <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={form.cost_center_id} onChange={(e) => setForm({ ...form, cost_center_id: e.target.value })}>
                <option value="">Selecione</option>
                {despesaCostCenters.map((cc) => <option key={cc.id} value={cc.id}>{cc.code} — {cc.name}</option>)}
              </select></div>
          </>)}
          <div><label className="block text-sm font-medium text-foreground mb-1">Valor Planejado (R$)</label><Input type="number" step="0.01" value={form.planned_amount} onChange={(e) => setForm({ ...form, planned_amount: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Descrição</label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
