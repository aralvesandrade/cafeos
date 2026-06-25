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
import { ArrowLeft, Calendar, Package, Plus } from 'lucide-react'

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
  const [plots, setPlots] = useState<Plot[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [finalizing, setFinalizing] = useState(false)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [form, setForm] = useState({ plot_id: '', quantity: '', notes: '' })

  const load = useCallback(async () => {
    try {
      const [harvestData, prodData, plotsData] = await Promise.all([
        apiRequest<Harvest>(`/harvests/${harvestId}`),
        apiRequest<Production[]>(`/harvests/${harvestId}/production`),
        apiRequest<Plot[]>('/plots'),
      ])
      setHarvest(harvestData)
      setProductions(prodData)
      setPlots(plotsData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [harvestId])

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
    return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>
  }

  if (!harvest) {
    return (
      <div className="text-center py-16">
        <p className="text-coffee-text-light mb-4">Safra não encontrada.</p>
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
          <h1 className="text-2xl font-bold text-coffee-green-dark flex items-center gap-2">
            <Calendar className="h-6 w-6" />
            Safra {harvest.year}
          </h1>
          <Badge variant={statusInfo.variant}>{statusInfo.label}</Badge>
        </div>

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

      {harvest.description && (
        <p className="text-coffee-text-light -mt-4">{harvest.description}</p>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Package className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Estimativa</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{harvest.estimated_production} sacas</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Package className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Produzido</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{totalProduced} sacas</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Calendar className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Criada em</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{harvest.created_at}</p>
          </CardContent>
        </Card>
      </div>

      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-coffee-green-dark flex items-center gap-2">
            <Package className="h-5 w-5" />
            Produção
          </h2>
          {harvest.status !== 'finalizada' && (
            <Button onClick={() => setDialogOpen(true)}>
              <Plus className="h-4 w-4" />
              Registrar Produção
            </Button>
          )}
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
                  <TableCell className="text-coffee-text-light">{p.notes || '—'}</TableCell>
                  <TableCell>{p.recorded_at}</TableCell>
                </TableRow>
              )
            })}
            {productions.length === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center text-coffee-text-light py-8">
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
            <label className="block text-sm font-medium text-coffee-text mb-1">Talhão</label>
            <Select value={form.plot_id} onChange={(e) => setForm({ ...form, plot_id: e.target.value })} required>
              <option value="">Selecione...</option>
              {plots.map((pl) => (
                <option key={pl.id} value={pl.id}>{pl.name}</option>
              ))}
            </Select>
          </div>
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Quantidade (sacas)</label>
            <Input
              type="number"
              step="0.01"
              value={form.quantity}
              onChange={(e) => setForm({ ...form, quantity: e.target.value })}
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Observações</label>
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
