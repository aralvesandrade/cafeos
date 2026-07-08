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
    return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>
  }

  if (!plot) {
    return (
      <div className="text-center py-16">
        <p className="text-coffee-text-light mb-4">Talhão não encontrado.</p>
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
        <h1 className="text-2xl font-bold text-coffee-green-dark">{plot.name}</h1>
      </div>

      {farm && (
        <p className="text-sm text-coffee-text-light flex items-center gap-1 -mt-4">
          Fazenda: <Link to={`/farms/${farm.id}`} className="text-coffee-green hover:underline font-medium">{farm.name}</Link>
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
            <Ruler className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Área</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{plot.area_ha} ha</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Sprout className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Cultivar</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{plot.cultivar || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Grid3X3 className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Tipo de Solo</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text capitalize">{plot.soil_type || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Mountain className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Altitude</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{plot.altitude} m</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <CalendarDays className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Ano Plantio</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{plot.planting_year}</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center gap-2 pb-2">
          <Droplets className="h-4 w-4 text-coffee-green" />
          <CardTitle className="text-sm font-medium text-coffee-text-light">Informações da atividade</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <p className="text-coffee-text-light">Tipo</p>
              <p className="font-medium text-coffee-text capitalize">{plot.crop_type || '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Irrigação</p>
              <p className="font-medium text-coffee-text">{plot.irrigation || '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Custo de formação/ha</p>
              <p className="font-medium text-coffee-text">{plot.formation_cost_per_ha ? `R$ ${plot.formation_cost_per_ha.toFixed(2)}` : '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Vida útil</p>
              <p className="font-medium text-coffee-text">{plot.useful_life_years ? `${plot.useful_life_years} anos` : '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Espaçamento (linha x planta)</p>
              <p className="font-medium text-coffee-text">{plot.row_spacing_m || plot.plant_spacing_m ? `${plot.row_spacing_m} x ${plot.plant_spacing_m} m` : '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Nº de plantas</p>
              <p className="font-medium text-coffee-text">{plot.plant_count || '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Data de plantio</p>
              <p className="font-medium text-coffee-text">{plot.planting_date ? new Date(plot.planting_date).toLocaleDateString('pt-BR') : '—'}</p>
            </div>
            <div>
              <p className="text-coffee-text-light">Data de ativação</p>
              <p className="font-medium text-coffee-text">{plot.activation_date ? new Date(plot.activation_date).toLocaleDateString('pt-BR') : '—'}</p>
            </div>
          </div>
          {plot.notes && (
            <div className="mt-4 text-sm">
              <p className="text-coffee-text-light">Observações</p>
              <p className="font-medium text-coffee-text">{plot.notes}</p>
            </div>
          )}
        </CardContent>
      </Card>

      <div>
        <h2 className="text-xl font-semibold text-coffee-green-dark flex items-center gap-2 mb-4">
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
                <TableCell colSpan={4} className="text-center text-coffee-text-light py-8">
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
