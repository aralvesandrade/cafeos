import { useEffect, useState, useCallback } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { ArrowLeft, Grid3X3, Ruler, Sprout, Mountain, CalendarDays, Tractor, Droplets } from 'lucide-react'

interface Plot {
  id: string
  name: string
  farm_id: string
  area_ha: number
  cultivar: string
  soil_type: string
  altitude: number
  planting_year: number

  leased: boolean
  stage: string
  irrigation: string
  activation_date: string | null
  planting_date: string | null
  deactivation_date: string | null
  intercropped: boolean
  secondary_crop: string
  notes: string
  crop_type: string
  formation_cost_per_ha: number
  useful_life_years: number
  row_spacing_m: number
  plant_spacing_m: number
  plant_count: number
}

interface Farm {
  id: string
  name: string
}

interface Operation {
  id: string
  type: string
  description: string
  date: string
  cost: number
}

export function PlotDetail() {
  const { plotId } = useParams<{ plotId: string }>()
  const navigate = useNavigate()
  const [plot, setPlot] = useState<Plot | null>(null)
  const [farm, setFarm] = useState<Farm | null>(null)
  const [operations, setOperations] = useState<Operation[]>([])
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    try {
      const plotData = await apiRequest<Plot>(`/plots/${plotId}`)
      setPlot(plotData)

      const [farmData, opsData] = await Promise.all([
        apiRequest<Farm>(`/farms/${plotData.farm_id}`),
        apiRequest<Operation[]>(`/plots/${plotId}/operations`),
      ])
      setFarm(farmData)
      setOperations(opsData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [plotId])

  useEffect(() => { load() }, [load])

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  if (!plot) {
    return (
      <div className="text-center py-16">
        <p className="text-muted-foreground mb-4">Talhão não encontrado.</p>
        <Button variant="outline" onClick={() => navigate('/plots')}>Voltar</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" onClick={() => navigate('/plots')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-2xl font-bold text-primary">{plot.name}</h1>
      </div>

      {farm && (
        <p className="text-sm text-muted-foreground flex items-center gap-1 -mt-4">
          Fazenda: <Link to={`/farms/${farm.id}`} className="text-primary hover:underline font-medium">{farm.name}</Link>
        </p>
      )}

      <div className="flex items-center gap-2 -mt-2">
        <Badge variant={plot.stage === 'producao' ? 'success' : 'info'}>
          {plot.stage === 'producao' ? 'Produção' : 'Formação'}
        </Badge>
        {plot.leased && <Badge variant="warning">Arrendado</Badge>}
        {plot.intercropped && <Badge variant="default">Consorciada {plot.secondary_crop && `(${plot.secondary_crop})`}</Badge>}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Ruler className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Área</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{plot.area_ha} ha</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Sprout className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Cultivar</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{plot.cultivar || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Grid3X3 className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Tipo de Solo</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground capitalize">{plot.soil_type || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Mountain className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Altitude</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{plot.altitude} m</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <CalendarDays className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Ano Plantio</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{plot.planting_year}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center gap-2 pb-2">
          <Droplets className="h-4 w-4 text-primary" />
          <CardTitle className="text-sm font-medium text-muted-foreground">Informações da atividade</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <p className="text-muted-foreground">Tipo</p>
              <p className="font-medium text-foreground capitalize">{plot.crop_type || '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Irrigação</p>
              <p className="font-medium text-foreground">{plot.irrigation || '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Custo de formação/ha</p>
              <p className="font-medium text-foreground">{plot.formation_cost_per_ha ? `R$ ${plot.formation_cost_per_ha.toFixed(2)}` : '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Vida útil</p>
              <p className="font-medium text-foreground">{plot.useful_life_years ? `${plot.useful_life_years} anos` : '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Espaçamento (linha x planta)</p>
              <p className="font-medium text-foreground">{plot.row_spacing_m || plot.plant_spacing_m ? `${plot.row_spacing_m} x ${plot.plant_spacing_m} m` : '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Nº de plantas</p>
              <p className="font-medium text-foreground">{plot.plant_count || '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Data de plantio</p>
              <p className="font-medium text-foreground">{plot.planting_date ? new Date(plot.planting_date).toLocaleDateString('pt-BR') : '—'}</p>
            </div>
            <div>
              <p className="text-muted-foreground">Data de ativação</p>
              <p className="font-medium text-foreground">{plot.activation_date ? new Date(plot.activation_date).toLocaleDateString('pt-BR') : '—'}</p>
            </div>
          </div>
          {plot.notes && (
            <div className="mt-4 text-sm">
              <p className="text-muted-foreground">Observações</p>
              <p className="font-medium text-foreground">{plot.notes}</p>
            </div>
          )}
        </CardContent>
      </Card>

      <div>
        <h2 className="text-xl font-semibold text-primary flex items-center gap-2 mb-4">
          <Tractor className="h-5 w-5" />
          Operações
        </h2>

        <Table>
          <TableHead>
            <TableRow>
              <TableHeader>Data</TableHeader>
              <TableHeader>Tipo</TableHeader>
              <TableHeader>Descrição</TableHeader>
              <TableHeader>Custo</TableHeader>
            </TableRow>
          </TableHead>
          <TableBody>
            {operations.map((op) => (
              <TableRow key={op.id}>
                <TableCell>{op.date}</TableCell>
                <TableCell>
                  <Badge variant="info">{op.type}</Badge>
                </TableCell>
                <TableCell>{op.description || '—'}</TableCell>
                <TableCell>
                  {op.cost > 0 ? `R$ ${op.cost.toFixed(2)}` : '—'}
                </TableCell>
              </TableRow>
            ))}
            {operations.length === 0 && (
              <TableRow>
                <TableCell colSpan={4} className="text-center text-muted-foreground py-8">
                  Nenhuma operação registrada para este talhão.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
