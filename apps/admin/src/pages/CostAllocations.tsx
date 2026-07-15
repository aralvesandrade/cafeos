import { useEffect, useState, useCallback } from 'react'
import { useSearchParams } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Trash2, SplitSquareHorizontal } from 'lucide-react'

interface AllocationItem {
  id: string
  plot_id: string
  plot_name: string
  amount: number
  percentage: number
}

interface Allocation {
  id: string
  harvest_id: string
  cost_center_id: string
  cost_center_name: string
  description: string
  total_amount: number
  method: 'area_proportional' | 'custom_percentage'
  date: string
  items: AllocationItem[]
}

interface Harvest { id: string; description: string }
interface CostCenter { id: string; name: string; code: string; type: 'receita' | 'despesa' }
interface Plot { id: string; name: string }

const methodLabels: Record<string, string> = {
  area_proportional: 'Proporcional por Área',
  custom_percentage: 'Percentual Customizado',
}

const emptyForm = { cost_center_id: '', total_amount: '', description: '', date: '', method: 'area_proportional' as 'area_proportional' | 'custom_percentage' }

export function CostAllocations() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [allocations, setAllocations] = useState<Allocation[]>([])
  const [harvests, setHarvests] = useState<Harvest[]>([])
  const [costCenters, setCostCenters] = useState<CostCenter[]>([])
  const [plots, setPlots] = useState<Plot[]>([])
  const [harvestId, setHarvestId] = useState(searchParams.get('harvest_id') || '')
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [form, setForm] = useState(emptyForm)
  const [percentages, setPercentages] = useState<Record<string, string>>({})
  const [saving, setSaving] = useState(false)
  const toast = useToast()
  const confirm = useConfirm()

  const load = useCallback(async () => {
    try {
      const [harvestsData, ccData, plotsData] = await Promise.all([
        apiRequest<Harvest[]>('/harvests'),
        apiRequest<CostCenter[]>('/cost-centers'),
        apiRequest<Plot[]>('/plots'),
      ])
      setHarvests(harvestsData)
      setCostCenters(ccData)
      setPlots(plotsData)
      const hid = harvestId || harvestsData[0]?.id || ''
      if (hid && !harvestId) setHarvestId(hid)
      if (hid) {
        const allocData = await apiRequest<Allocation[]>(`/harvests/${hid}/cost-allocations`)
        setAllocations(allocData)
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
      const body: Record<string, unknown> = {
        harvest_id: harvestId,
        cost_center_id: form.cost_center_id,
        total_amount: parseFloat(form.total_amount) || 0,
        description: form.description,
        method: form.method,
        date: form.date,
      }
      if (form.method === 'custom_percentage') {
        const pctMap: Record<string, number> = {}
        for (const [plotId, value] of Object.entries(percentages)) {
          const n = parseFloat(value)
          if (n > 0) pctMap[plotId] = n
        }
        body.percentages = pctMap
      }
      await apiRequest('/cost-allocations', { method: 'POST', body })
      setDialogOpen(false)
      setForm(emptyForm)
      setPercentages({})
      await load()
      toast.success('Rateio criado')
    } catch (err) { console.error(err); toast.error('Erro ao criar rateio') }
    finally { setSaving(false) }
  }

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover rateio?', variant: 'danger' }))) return
    try { await apiRequest(`/cost-allocations/${id}`, { method: 'DELETE' }); await load(); toast.success('Rateio removido') }
    catch (err) { console.error(err); toast.error('Erro ao remover rateio') }
  }

  const despesaCostCenters = costCenters.filter((cc) => cc.type === 'despesa')
  const percentageTotal = Object.values(percentages).reduce((sum, v) => sum + (parseFloat(v) || 0), 0)

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-primary">Rateio de Custo</h1>
          <p className="text-sm text-muted-foreground">Distribuição de custos por talhão</p>
        </div>
        <Button onClick={() => { setForm(emptyForm); setPercentages({}); setDialogOpen(true) }}>
          <Plus className="h-4 w-4" /> Novo Rateio
        </Button>
      </div>

      <Select value={harvestId} onChange={(e) => setHarvestId(e.target.value)} className="w-64">
        {harvests.map((h) => (
          <option key={h.id} value={h.id}>{h.description}</option>
        ))}
      </Select>

      {allocations.length === 0 && (
        <p className="text-center text-muted-foreground py-8">Nenhum rateio cadastrado para esta safra.</p>
      )}

      {allocations.map((a) => (
        <Card key={a.id}>
          <CardHeader className="flex flex-row items-center justify-between gap-2 pb-2">
            <div className="flex items-center gap-2">
              <SplitSquareHorizontal className="h-4 w-4 text-primary" />
              <CardTitle className="text-sm font-medium text-foreground">{a.cost_center_name}</CardTitle>
              <Badge variant="info">{methodLabels[a.method] || a.method}</Badge>
            </div>
            <div className="flex items-center gap-3">
              <span className="text-sm text-muted-foreground">{new Date(a.date).toLocaleDateString('pt-BR')}</span>
              <span className="font-semibold text-foreground">R$ {a.total_amount.toFixed(2)}</span>
              <Button variant="ghost" size="sm" onClick={() => handleDelete(a.id)}><Trash2 className="h-4 w-4 text-destructive" /></Button>
            </div>
          </CardHeader>
          <CardContent>
            {a.description && <p className="text-sm text-muted-foreground mb-3">{a.description}</p>}
            <Table>
              <TableHead>
                <TableRow>
                  <TableHeader>Talhão</TableHeader>
                  <TableHeader className="text-right">Valor</TableHeader>
                  <TableHeader className="text-right">%</TableHeader>
                </TableRow>
              </TableHead>
              <TableBody>
                {a.items.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell>{item.plot_name || item.plot_id}</TableCell>
                    <TableCell className="text-right">R$ {item.amount.toFixed(2)}</TableCell>
                    <TableCell className="text-right">{item.percentage.toFixed(2)}%</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      ))}

      <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} title="Novo Rateio">
        <div className="space-y-4">
          <div><label className="block text-sm font-medium text-foreground mb-1">Centro de Custo</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={form.cost_center_id} onChange={(e) => setForm({ ...form, cost_center_id: e.target.value })}>
              <option value="">Selecione</option>
              {despesaCostCenters.map((cc) => <option key={cc.id} value={cc.id}>{cc.code} — {cc.name}</option>)}
            </select></div>
          <div className="grid grid-cols-2 gap-3">
            <div><label className="block text-sm font-medium text-foreground mb-1">Valor Total (R$)</label><Input type="number" step="0.01" value={form.total_amount} onChange={(e) => setForm({ ...form, total_amount: e.target.value })} /></div>
            <div><label className="block text-sm font-medium text-foreground mb-1">Data</label><Input type="date" value={form.date} onChange={(e) => setForm({ ...form, date: e.target.value })} /></div>
          </div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Descrição</label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Método</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={form.method} onChange={(e) => setForm({ ...form, method: e.target.value as 'area_proportional' | 'custom_percentage' })}>
              <option value="area_proportional">Proporcional por Área</option>
              <option value="custom_percentage">Percentual Customizado</option>
            </select></div>

          {form.method === 'custom_percentage' && (
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <label className="block text-sm font-medium text-foreground">Percentual por Talhão</label>
                <span className={`text-sm ${Math.abs(percentageTotal - 100) < 0.1 ? 'text-success-foreground' : 'text-destructive'}`}>Total: {percentageTotal.toFixed(1)}%</span>
              </div>
              <div className="space-y-2 max-h-56 overflow-y-auto">
                {plots.map((p) => (
                  <div key={p.id} className="flex items-center gap-3">
                    <span className="flex-1 text-sm text-foreground">{p.name}</span>
                    <Input
                      type="number"
                      step="0.01"
                      className="w-28"
                      value={percentages[p.id] || ''}
                      onChange={(e) => setPercentages({ ...percentages, [p.id]: e.target.value })}
                      placeholder="0"
                    />
                  </div>
                ))}
              </div>
            </div>
          )}

          <div className="flex justify-end gap-3 pt-2">
            <Button variant="outline" onClick={() => setDialogOpen(false)}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
