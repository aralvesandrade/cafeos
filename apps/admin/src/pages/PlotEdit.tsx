import { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { PlotForm, type PlotData } from '@/components/plots/PlotForm'
import { ArrowLeft } from 'lucide-react'
import type { Plot } from '@/pages/Plots'

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

export function PlotEdit() {
  const { plotId } = useParams<{ plotId: string }>()
  const navigate = useNavigate()
  const isEditing = Boolean(plotId)
  const [initial, setInitial] = useState<PlotData | undefined>(undefined)
  const [farms, setFarms] = useState<Farm[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const [farmsData, plot] = await Promise.all([
        apiRequest<Farm[]>('/farms'),
        plotId ? apiRequest<Plot>(`/plots/${plotId}`) : Promise.resolve(null),
      ])
      setFarms(farmsData)
      if (plot) setInitial(plotToFormData(plot))
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [plotId])

  useEffect(() => { load() }, [load])

  const handleSave = async (formData: PlotData) => {
    setSaving(true)
    try {
      const payload = formDataToPayload(formData)
      if (isEditing) {
        await apiRequest(`/plots/${plotId}`, { method: 'PUT', body: payload })
      } else {
        await apiRequest('/plots', { method: 'POST', body: payload })
      }
      navigate('/plots')
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Button variant="ghost" size="sm" onClick={() => navigate('/plots')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="text-2xl font-bold text-primary">{isEditing ? 'Editar Talhão' : 'Novo Talhão'}</h1>
      </div>

      <PlotForm
        initial={initial}
        farms={farms}
        onSave={handleSave}
        onCancel={() => navigate('/plots')}
        loading={saving}
      />
    </div>
  )
}
