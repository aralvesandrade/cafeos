import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { PlotForm, type PlotData } from '@/components/plots/PlotForm'
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

function plotToFormData(plot: Plot): PlotData {
  return {
    name: plot.name,
    farm_id: plot.farm_id,
    area: String(plot.area_ha),
    cultivar: plot.cultivar ?? '',
    soil_type: plot.soil_type ?? '',
    altitude: String(plot.altitude ?? ''),
    planting_year: String(plot.planting_year ?? ''),
    leased: plot.leased ?? false,
    stage: plot.stage || 'formacao',
    irrigation: plot.irrigation ?? '',
    activation_date: plot.activation_date ? plot.activation_date.slice(0, 10) : '',
    planting_date: plot.planting_date ? plot.planting_date.slice(0, 10) : '',
    deactivation_date: plot.deactivation_date ? plot.deactivation_date.slice(0, 10) : '',
    intercropped: plot.intercropped ?? false,
    secondary_crop: plot.secondary_crop ?? '',
    notes: plot.notes ?? '',
    crop_type: plot.crop_type ?? '',
    formation_cost_per_ha: String(plot.formation_cost_per_ha ?? ''),
    useful_life_years: String(plot.useful_life_years ?? ''),
    row_spacing_m: String(plot.row_spacing_m ?? ''),
    plant_spacing_m: String(plot.plant_spacing_m ?? ''),
    plant_count: String(plot.plant_count ?? ''),
    dam_area: String(plot.dam_area_ha ?? ''),
    improvements_area: String(plot.improvements_area_ha ?? ''),
    roads_area: String(plot.roads_area_ha ?? ''),
    app_area: String(plot.app_area_ha ?? ''),
    legal_reserve_area: String(plot.legal_reserve_area_ha ?? ''),
  }
}

function formDataToPayload(form: PlotData) {
  return {
    name: form.name,
    farm_id: form.farm_id,
    area_ha: parseFloat(form.area) || 0,
    cultivar: form.cultivar,
    soil_type: form.soil_type,
    altitude: parseInt(form.altitude) || 0,
    planting_year: parseInt(form.planting_year) || 0,
    leased: form.leased,
    stage: form.stage,
    irrigation: form.irrigation,
    activation_date: form.activation_date || null,
    planting_date: form.planting_date || null,
    deactivation_date: form.deactivation_date || null,
    intercropped: form.intercropped,
    secondary_crop: form.secondary_crop,
    notes: form.notes,
    crop_type: form.crop_type,
    formation_cost_per_ha: parseFloat(form.formation_cost_per_ha) || 0,
    useful_life_years: parseInt(form.useful_life_years) || 0,
    row_spacing_m: parseFloat(form.row_spacing_m) || 0,
    plant_spacing_m: parseFloat(form.plant_spacing_m) || 0,
    plant_count: parseInt(form.plant_count) || 0,
    dam_area_ha: parseFloat(form.dam_area) || 0,
    improvements_area_ha: parseFloat(form.improvements_area) || 0,
    roads_area_ha: parseFloat(form.roads_area) || 0,
    app_area_ha: parseFloat(form.app_area) || 0,
    legal_reserve_area_ha: parseFloat(form.legal_reserve_area) || 0,
  }
}

export function Plots() {
  const [plots, setPlots] = useState<Plot[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [farmFilter, setFarmFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Plot | null>(null)
  const [saving, setSaving] = useState(false)

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

  const handleSave = async (formData: PlotData) => {
    setSaving(true)
    try {
      const payload = formDataToPayload(formData)

      if (editing) {
        await apiRequest(`/plots/${editing.id}`, { method: 'PUT', body: payload })
      } else {
        await apiRequest('/plots', { method: 'POST', body: payload })
      }

      setDialogOpen(false)
      setEditing(null)
      await loadData()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Remover talhão?')) return
    try {
      await apiRequest(`/plots/${id}`, { method: 'DELETE' })
      await loadData()
    } catch (err) { console.error(err) }
  }

  const openEdit = (plot: Plot) => {
    setEditing(plot)
    setDialogOpen(true)
  }

  const openCreate = () => {
    setEditing(null)
    setDialogOpen(true)
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
          <Button onClick={openCreate}><Plus className="h-4 w-4" /> Novo Talhão</Button>
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
                  <Button variant="ghost" size="sm" onClick={() => openEdit(plot)}><Pencil className="h-4 w-4" /></Button>
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

      <Dialog
        open={dialogOpen}
        onClose={() => { setDialogOpen(false); setEditing(null) }}
        title={editing ? 'Editar Talhão' : 'Novo Talhão'}
        className="max-w-3xl"
      >
        <PlotForm
          initial={editing ? plotToFormData(editing) : undefined}
          farms={farms}
          onSave={handleSave}
          onCancel={() => { setDialogOpen(false); setEditing(null) }}
          loading={saving}
        />
      </Dialog>
    </div>
  )
}
