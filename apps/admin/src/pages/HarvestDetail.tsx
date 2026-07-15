import { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Dialog } from '@/components/ui/dialog'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, Calendar, Package, Plus, TrendingUp, Wallet, SplitSquareHorizontal } from 'lucide-react'

interface Harvest {
  id: string
  year: number
  description: string
  estimated_production: number
  status: string
  created_at: string
}

interface Production {
  id: string
  plot_id: string
  quantity: number
  recorded_at: string
  notes: string
}

interface Plot {
  id: string
  name: string
}

interface Farm { id: string; name: string }

interface Indicator {
  id: string
  type: string
  value: number
}

const indicatorLabels: Record<string, string> = {
  sacas_por_hectare: 'Sacas/ha',
  custo_por_saca: 'Custo/saca',
  producao_total: 'Produção Total',
  custo_total: 'Custo Total',
  area_producao: 'Área de Produção',
  coe: 'COE',
  coe_por_area: 'COE/ha',
  coe_por_saca: 'COE/saca',
  cot: 'COT',
  cot_por_area: 'COT/ha',
  cot_por_saca: 'COT/saca',
  ct_producao: 'CT',
  ct_producao_por_area: 'CT/ha',
  ct_producao_por_saca: 'CT/saca',
}

const indicatorTitles: Record<string, string> = {
  sacas_por_hectare: 'Produtividade: total produzido (sacas) dividido pela área plantada (ha)',
  custo_por_saca: 'Custo total da safra dividido pela produção total (sacas)',
  producao_total: 'Soma de toda a produção registrada nesta safra',
  custo_total: 'Soma de todos os custos (operações, manutenções, mão de obra, despesas e rateios) no período da safra',
  area_producao: 'Soma da área de todos os talhões da organização',
  coe: 'Custo Operacional Efetivo: soma dos desembolsos operacionais diretos (insumos, mão de obra contratada, combustível etc.)',
  coe_por_area: 'COE dividido pela área de produção (ha)',
  coe_por_saca: 'COE dividido pela produção total (sacas)',
  cot: 'Custo Operacional Total: COE + mão de obra familiar + depreciação de máquinas e benfeitorias',
  cot_por_area: 'COT dividido pela área de produção (ha)',
  cot_por_saca: 'COT dividido pela produção total (sacas)',
  ct_producao: 'Custo Total: COT + remuneração do capital investido (terra e outros ativos)',
  ct_producao_por_area: 'CT dividido pela área de produção (ha)',
  ct_producao_por_saca: 'CT dividido pela produção total (sacas)',
}

const currencyIndicators = new Set(['custo_por_saca', 'custo_total', 'coe', 'coe_por_area', 'coe_por_saca', 'cot', 'cot_por_area', 'cot_por_saca', 'ct_producao', 'ct_producao_por_area', 'ct_producao_por_saca'])

function formatIndicator(type: string, value: number): string {
  if (currencyIndicators.has(type)) return `R$ ${value.toFixed(2)}`
  if (type === 'area_producao') return `${value.toFixed(2)} ha`
  if (type === 'producao_total') return `${value.toFixed(2)} sacas`
  return value.toFixed(2)
}

const statusLabels: Record<string, { label: string; variant: 'default' | 'success' | 'warning' }> = {
  planejada: { label: 'Planejada', variant: 'default' },
  em_andamento: { label: 'Em Andamento', variant: 'warning' },
  finalizada: { label: 'Finalizada', variant: 'success' },
}

export function HarvestDetail() {
  const { harvestId } = useParams<{ harvestId: string }>()
  const navigate = useNavigate()
  const [harvest, setHarvest] = useState<Harvest | null>(null)
  const [productions, setProductions] = useState<Production[]>([])
  const [indicators, setIndicators] = useState<Indicator[]>([])
  const [plots, setPlots] = useState<Plot[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [farmFilter, setFarmFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [finalizing, setFinalizing] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [form, setForm] = useState({ plot_id: '', quantity: '', notes: '' })

  const load = useCallback(async () => {
    try {
      const [harvestData, prodData, plotsData, farmsData] = await Promise.all([
        apiRequest<Harvest>(`/harvests/${harvestId}`),
        apiRequest<Production[]>(`/harvests/${harvestId}/production`, { params: farmFilter ? { farm_id: farmFilter } : undefined }),
        apiRequest<Plot[]>('/plots'),
        apiRequest<Farm[]>('/farms'),
      ])
      setHarvest(harvestData)
      setProductions(prodData)
      setPlots(plotsData)
      setFarms(farmsData)
      if (harvestData.status === 'finalizada') {
        const indicatorData = await apiRequest<Indicator[]>(`/harvests/${harvestId}/indicators`)
        setIndicators(indicatorData)
      }
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [harvestId, farmFilter])

  useEffect(() => { load() }, [load])

  const handleRecordProduction = async () => {
    setSaving(true)
    try {
      await apiRequest(`/harvests/${harvestId}/production`, {
        method: 'POST',
        body: {
          plot_id: form.plot_id,
          quantity: parseFloat(form.quantity),
          notes: form.notes,
        },
      })
      setDialogOpen(false)
      setForm({ plot_id: '', quantity: '', notes: '' })
      await load()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleFinalize = async () => {
    if (!confirm('Finalizar esta safra? Esta ação não pode ser desfeita.')) return
    setFinalizing(true)
    try {
      await apiRequest(`/harvests/${harvestId}/finalize`, { method: 'PUT' })
      await load()
    } catch (err) {
      console.error(err)
    } finally {
      setFinalizing(false)
    }
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  if (!harvest) {
    return (
      <div className="text-center py-16">
        <p className="text-muted-foreground mb-4">Safra não encontrada.</p>
        <Button variant="outline" onClick={() => navigate('/harvests')}>Voltar</Button>
      </div>
    )
  }

  const statusInfo = statusLabels[harvest.status] || statusLabels.planejada
  const totalProduced = productions.reduce((sum, p) => sum + p.quantity, 0)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="sm" onClick={() => navigate('/harvests')}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <h1 className="text-2xl font-bold text-primary flex items-center gap-2">
            <Calendar className="h-6 w-6" />
            Safra {harvest.year}
          </h1>
          <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
        </div>

        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate(`/budgets?harvest_id=${harvest.id}`)}>
            <Wallet className="h-4 w-4" />
            Ver Orçamento
          </Button>
          <Button variant="outline" onClick={() => navigate(`/cost-allocations?harvest_id=${harvest.id}`)}>
            <SplitSquareHorizontal className="h-4 w-4" />
            Ver Rateios
          </Button>
          {harvest.status !== 'finalizada' && (
          <Button
            variant="danger"
            onClick={handleFinalize}
            disabled={finalizing}
          >
            {finalizing ? 'Finalizando...' : 'Finalizar Safra'}
          </Button>
          )}
        </div>
      </div>

      {harvest.description && (
        <p className="text-muted-foreground -mt-4">{harvest.description}</p>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Package className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Estimativa</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{harvest.estimated_production} sacas</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Package className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Produzido</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{totalProduced} sacas</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Calendar className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Criada em</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{harvest.created_at}</p>
          </CardContent>
        </Card>
      </div>

      {harvest.status === 'finalizada' && (
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <TrendingUp className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Indicadores</CardTitle>
          </CardHeader>
          <CardContent>
            {indicators.length === 0 ? (
              <p className="text-sm text-muted-foreground">Nenhum indicador calculado.</p>
            ) : (
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                {indicators.map((ind) => (
                  <div key={ind.id} title={indicatorTitles[ind.type]}>
                    <p className="text-muted-foreground">{indicatorLabels[ind.type] || ind.type}</p>
                    <p className="font-medium text-foreground">{formatIndicator(ind.type, ind.value)}</p>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      <div>
        <div className="flex items-center justify-between mb-4 flex-wrap gap-3">
          <h2 className="text-xl font-semibold text-primary flex items-center gap-2">
            <Package className="h-5 w-5" />
            Produção
          </h2>
          <div className="flex gap-3">
            <Select value={farmFilter} onChange={(e) => setFarmFilter(e.target.value)} className="w-48">
              <option value="">Todas as fazendas</option>
              {farms.map((f) => (
                <option key={f.id} value={f.id}>{f.name}</option>
              ))}
            </Select>
            {harvest.status !== 'finalizada' && (
              <Button onClick={() => setDialogOpen(true)}>
                <Plus className="h-4 w-4" />
                Registrar Produção
              </Button>
            )}
          </div>
        </div>

        <Table>
          <TableHead>
            <TableRow>
              <TableHeader>Talhão</TableHeader>
              <TableHeader>Quantidade (sacas)</TableHeader>
              <TableHeader>Observações</TableHeader>
              <TableHeader>Registrado em</TableHeader>
            </TableRow>
          </TableHead>
          <TableBody>
            {productions.map((p) => {
              const plot = plots.find((pl) => pl.id === p.plot_id)
              return (
                <TableRow key={p.id}>
                  <TableCell className="font-medium">{plot?.name || p.plot_id}</TableCell>
                  <TableCell>{p.quantity}</TableCell>
                  <TableCell className="text-muted-foreground">{p.notes || '—'}</TableCell>
                  <TableCell>{p.recorded_at}</TableCell>
                </TableRow>
              )
            })}
            {productions.length === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center text-muted-foreground py-8">
                  Nenhuma produção registrada.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <Dialog
        open={dialogOpen}
        onClose={() => { setDialogOpen(false); setForm({ plot_id: '', quantity: '', notes: '' }) }}
        title="Registrar Produção"
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Talhão</label>
            <Select value={form.plot_id} onChange={(e) => setForm({ ...form, plot_id: e.target.value })} required>
              <option value="">Selecione...</option>
              {plots.map((pl) => (
                <option key={pl.id} value={pl.id}>{pl.name}</option>
              ))}
            </Select>
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Quantidade (sacas)</label>
            <Input
              type="number"
              step="0.01"
              value={form.quantity}
              onChange={(e) => setForm({ ...form, quantity: e.target.value })}
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Observações</label>
            <Input
              value={form.notes}
              onChange={(e) => setForm({ ...form, notes: e.target.value })}
              placeholder="Opcional"
            />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button
              variant="outline"
              onClick={() => { setDialogOpen(false); setForm({ plot_id: '', quantity: '', notes: '' }) }}
            >
              Cancelar
            </Button>
            <Button onClick={handleRecordProduction} disabled={saving}>
              {saving ? 'Salvando...' : 'Registrar'}
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
