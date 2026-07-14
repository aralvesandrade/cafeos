import { useEffect, useState, useCallback } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, CalendarDays, User, Package, Ruler, DollarSign, Pencil, Trash2 } from 'lucide-react'

interface Operation {
  id: string
  plot_id: string
  plot_name: string
  harvest_id: string | null
  cost_center_id: string | null
  type_id: string
  type_name: string
  type_color: 'info' | 'warning' | 'success' | 'default' | 'danger'
  date: string
  responsible: string
  product_used: string
  quantity: number
  cost: number
  notes: string
}

interface Plot { id: string; name: string; farm_id: string }
interface Farm { id: string; name: string }
interface CostCenter { id: string; name: string }
interface Harvest { id: string; year: string }

export function OperationDetail() {
  const { operationId } = useParams<{ operationId: string }>()
  const navigate = useNavigate()
  const [operation, setOperation] = useState<Operation | null>(null)
  const [plot, setPlot] = useState<Plot | null>(null)
  const [farm, setFarm] = useState<Farm | null>(null)
  const [costCenter, setCostCenter] = useState<CostCenter | null>(null)
  const [harvest, setHarvest] = useState<Harvest | null>(null)
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    try {
      const opData = await apiRequest<Operation>(`/operations/${operationId}`)
      setOperation(opData)

      const [plotData, costCenters, harvests] = await Promise.all([
        apiRequest<Plot>(`/plots/${opData.plot_id}`),
        opData.cost_center_id ? apiRequest<CostCenter[]>('/cost-centers') : Promise.resolve([]),
        opData.harvest_id ? apiRequest<Harvest[]>('/harvests') : Promise.resolve([]),
      ])
      setPlot(plotData)
      if (plotData) {
        const farmData = await apiRequest<Farm>(`/farms/${plotData.farm_id}`)
        setFarm(farmData)
      }
      if (opData.cost_center_id) setCostCenter(costCenters.find((c) => c.id === opData.cost_center_id) || null)
      if (opData.harvest_id) setHarvest(harvests.find((h) => h.id === opData.harvest_id) || null)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [operationId])

  useEffect(() => { load() }, [load])

  const handleDelete = async () => {
    if (!operation || !confirm('Remover operação?')) return
    try {
      await apiRequest(`/operations/${operation.id}`, { method: 'DELETE' })
      navigate('/operations')
    } catch (err) { console.error(err) }
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  if (!operation) {
    return (
      <div className="text-center py-16">
        <p className="text-muted-foreground mb-4">Operação não encontrada.</p>
        <Button variant="outline" onClick={() => navigate('/operations')}>Voltar</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="sm" onClick={() => navigate('/operations')}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <h1 className="text-2xl font-bold text-primary">Operação</h1>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={() => navigate(`/operations?edit=${operation.id}`)}>
            <Pencil className="h-4 w-4" />
            Editar
          </Button>
          <Button variant="outline" size="sm" onClick={handleDelete}>
            <Trash2 className="h-4 w-4 text-destructive" />
            Excluir
          </Button>
        </div>
      </div>

      <div className="flex items-center gap-2 -mt-2">
        <Badge variant={operation.type_color || 'default'}>{operation.type_name || '—'}</Badge>
      </div>

      {plot && (
        <p className="text-sm text-muted-foreground flex items-center gap-1 -mt-2">
          Talhão: <Link to={`/plots/${plot.id}`} className="text-primary hover:underline font-medium">{plot.name}</Link>
          {farm && (<>
            {' '}· Fazenda: <Link to={`/farms/${farm.id}`} className="text-primary hover:underline font-medium">{farm.name}</Link>
          </>)}
        </p>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <CalendarDays className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Data</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{new Date(operation.date).toLocaleDateString('pt-BR')}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <User className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Responsável</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{operation.responsible || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Ruler className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Quantidade</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{operation.quantity || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <DollarSign className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Custo</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">R$ {operation.cost.toFixed(2)}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center gap-2 pb-2">
          <Package className="h-4 w-4 text-primary" />
          <CardTitle className="text-sm font-medium text-muted-foreground">Detalhes</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">Produto Usado</p>
              <p className="font-medium text-foreground">{operation.product_used || '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Centro de Custo</p>
              <p className="font-medium text-foreground">{costCenter?.name || '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Safra</p>
              <p className="font-medium text-foreground">{harvest?.year || '—'}</p>
            </div>
          </div>
          {operation.notes && (
            <div className="mt-4 text-sm">
              <p className="text-muted-foreground">Observações</p>
              <p className="font-medium text-foreground">{operation.notes}</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
