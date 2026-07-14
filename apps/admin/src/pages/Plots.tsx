import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { Plus, Pencil, Trash2, ExternalLink } from 'lucide-react'

export interface Plot {
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
  dam_area_ha: number
  improvements_area_ha: number
  roads_area_ha: number
  app_area_ha: number
  legal_reserve_area_ha: number
}

interface Farm { id: string; name: string }

export function Plots() {
  const [plots, setPlots] = useState<Plot[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [farmFilter, setFarmFilter] = useState('')
  const [loading, setLoading] = useState(true)

  const loadData = useCallback(async () => {
    try {
      const [plotsData, farmsData] = await Promise.all([
        apiRequest<Plot[]>('/plots', { params: farmFilter ? { farm_id: farmFilter } : undefined }),
        apiRequest<Farm[]>('/farms'),
      ])
      setPlots(plotsData)
      setFarms(farmsData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [farmFilter])

  useEffect(() => { loadData() }, [loadData])

  const handleDelete = async (id: string) => {
    if (!confirm('Remover talhão?')) return
    try {
      await apiRequest(`/plots/${id}`, { method: 'DELETE' })
      await loadData()
    } catch (err) { console.error(err) }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-primary">Talhões</h1>
          <p className="text-sm text-muted-foreground">Gerencie os talhões das suas fazendas</p>
        </div>
        <div className="flex gap-3">
          <Select value={farmFilter} onChange={(e) => setFarmFilter(e.target.value)} className="w-48">
            <option value="">Todas as fazendas</option>
            {farms.map((f) => (
              <option key={f.id} value={f.id}>{f.name}</option>
            ))}
          </Select>
          <Button asChild>
            <Link to="/plots/new">
              <Plus className="h-4 w-4" /> Novo Talhão
            </Link>
          </Button>
        </div>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Fazenda</TableHeader>
            <TableHeader>Área</TableHeader>
            <TableHeader>Estágio</TableHeader>
            <TableHeader>Variedade</TableHeader>
            <TableHeader>Solo</TableHeader>
            <TableHeader>Altitude</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {plots.map((plot) => (
            <TableRow key={plot.id}>
              <TableCell className="font-medium">
                <Link to={`/plots/${plot.id}`} className="text-primary hover:underline">
                  {plot.name}
                </Link>
              </TableCell>
              <TableCell>{farms.find((f) => f.id === plot.farm_id)?.name || plot.farm_id}</TableCell>
              <TableCell>{plot.area_ha} ha</TableCell>
              <TableCell>
                <Badge variant={plot.stage === 'producao' ? 'success' : 'info'}>
                  {plot.stage === 'producao' ? 'Produção' : 'Formação'}
                </Badge>
              </TableCell>
              <TableCell>{plot.cultivar}</TableCell>
              <TableCell>{plot.soil_type}</TableCell>
              <TableCell>{plot.altitude} m</TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Button variant="ghost" size="sm" asChild>
                    <Link to={`/plots/${plot.id}`}>
                      <ExternalLink className="h-4 w-4" />
                    </Link>
                  </Button>
                  <Button variant="ghost" size="sm" asChild>
                    <Link to={`/plots/${plot.id}/edit`}>
                      <Pencil className="h-4 w-4" />
                    </Link>
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDelete(plot.id)}><Trash2 className="h-4 w-4 text-destructive" /></Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
          {plots.length === 0 && (
            <TableRow><TableCell colSpan={8} className="text-center text-muted-foreground py-8">Nenhum talhão cadastrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
