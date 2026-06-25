import { useEffect, useState, useCallback } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { MapPin, ArrowLeft, Grid3X3, Ruler, HardHat, Map } from 'lucide-react'

interface Farm {
  id: string
  name: string
  owner: string
  location: string
  total_area_ha: number
  planted_area_ha: number
  created_at: string
}

interface Plot {
  id: string
  name: string
  area_ha: number
  cultivar: string
  soil_type: string
  altitude: number
  planting_year: number
}

export function FarmDetail() {
  const { farmId } = useParams<{ farmId: string }>()
  const navigate = useNavigate()
  const [farm, setFarm] = useState<Farm | null>(null)
  const [plots, setPlots] = useState<Plot[]>([])
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    try {
      const [farmData, plotsData] = await Promise.all([
        apiRequest<Farm>(`/farms/${farmId}`),
        apiRequest<Plot[]>(`/farms/${farmId}/plots`),
      ])
      setFarm(farmData)
      setPlots(plotsData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [farmId])

  useEffect(() => { load() }, [load])

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>
  }

  if (!farm) {
    return (
      <div className="text-center py-16">
        <p className="text-coffee-text-light mb-4">Fazenda não encontrada.</p>
        <Button variant="outline" onClick={() => navigate('/farms')}>Voltar</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" onClick={() => navigate('/farms')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-2xl font-bold text-coffee-green-dark">{farm.name}</h1>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <HardHat className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Proprietário</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{farm.owner || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Map className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Localização</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{farm.location || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Ruler className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Área Total</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{farm.total_area_ha} ha</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Ruler className="h-4 w-4 text-coffee-green" />
            <CardTitle className="text-sm font-medium text-coffee-text-light">Área Plantada</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-coffee-text">{farm.planted_area_ha} ha</p>
          </CardContent>
        </Card>
      </div>

      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-coffee-green-dark flex items-center gap-2">
            <Grid3X3 className="h-5 w-5" />
            Talhões
          </h2>
          <Button variant="outline" size="sm" onClick={() => navigate('/plots')}>
            Gerenciar Talhões
          </Button>
        </div>

        <Table>
          <TableHead>
            <TableRow>
              <TableHeader>Nome</TableHeader>
              <TableHeader>Área</TableHeader>
              <TableHeader>Cultivar</TableHeader>
              <TableHeader>Solo</TableHeader>
              <TableHeader>Altitude</TableHeader>
              <TableHeader>Ano Plantio</TableHeader>
              <TableHeader className="text-right">Ações</TableHeader>
            </TableRow>
          </TableHead>
          <TableBody>
            {plots.map((plot) => (
              <TableRow key={plot.id}>
                <TableCell className="font-medium">{plot.name}</TableCell>
                <TableCell>{plot.area_ha} ha</TableCell>
                <TableCell>{plot.cultivar}</TableCell>
                <TableCell>{plot.soil_type}</TableCell>
                <TableCell>{plot.altitude} m</TableCell>
                <TableCell>{plot.planting_year}</TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="sm" asChild>
                    <Link to={`/plots/${plot.id}`}>Detalhes</Link>
                  </Button>
                </TableCell>
              </TableRow>
            ))}
            {plots.length === 0 && (
              <TableRow>
                <TableCell colSpan={7} className="text-center text-coffee-text-light py-8">
                  Nenhum talhão cadastrado nesta fazenda.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
