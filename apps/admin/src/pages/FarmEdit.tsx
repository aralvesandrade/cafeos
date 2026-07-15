import { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { FarmForm, emptyProducer, type FarmData } from '@/components/farms/FarmForm'
import { ArrowLeft } from 'lucide-react'
import type { Farm } from '@/pages/Farms'
import { useToast } from '@/lib/toast'

function farmToFormData(farm: Farm): FarmData {
  return {
    name: farm.name,
    owner: farm.owner,
    location: farm.location,
    total_area: String(farm.total_area_ha),
    planted_area: String(farm.planted_area_ha),
    phone: farm.phone ?? '',
    activities: farm.activities ?? '',
    main_crop: farm.main_crop ?? '',
    secondary_crop: farm.secondary_crop ?? '',
    state: farm.state ?? '',
    city: farm.city ?? '',
    address: farm.address ?? '',
    production_system: farm.production_system ?? '',
    commercialization_product: farm.commercialization_product ?? '',
    has_no_cnpj: farm.has_no_cnpj ?? false,
    cnpj: farm.cnpj ?? '',
    has_no_nirf: farm.has_no_nirf ?? false,
    nirf: farm.nirf ?? '',
    has_no_incra: farm.has_no_incra ?? false,
    incra_registration: farm.incra_registration ?? '',
    has_no_state_registration: farm.has_no_state_registration ?? false,
    state_registration: farm.state_registration ?? '',
    has_no_dap: farm.has_no_dap ?? false,
    dap: farm.dap ?? '',
    has_no_car: farm.has_no_car ?? false,
    car: farm.car ?? '',
    fully_leased: farm.fully_leased ?? false,
    land_value_per_ha: String(farm.land_value_per_ha ?? ''),
    dam_area: String(farm.dam_area_ha ?? ''),
    improvements_area: String(farm.improvements_area_ha ?? ''),
    roads_area: String(farm.roads_area_ha ?? ''),
    app_area: String(farm.app_area_ha ?? ''),
    legal_reserve_area: String(farm.legal_reserve_area_ha ?? ''),
    native_vegetation_area: String(farm.native_vegetation_area_ha ?? ''),
    livestock_area_not_covered: String(farm.livestock_area_not_covered_ha ?? ''),
    agriculture_area_not_covered: String(farm.agriculture_area_not_covered_ha ?? ''),
    non_agricultural_area: String(farm.non_agricultural_area_ha ?? ''),
    producer: farm.producer
      ? {
          cpf: farm.producer.cpf ?? '',
          name: farm.producer.name ?? '',
          rg: farm.producer.rg ?? '',
          issuing_body: farm.producer.issuing_body ?? '',
          gender: farm.producer.gender ?? '',
          birth_date: farm.producer.birth_date ? farm.producer.birth_date.slice(0, 10) : '',
          marital_status: farm.producer.marital_status ?? '',
          phone: farm.producer.phone ?? '',
          email: farm.producer.email ?? '',
          education: farm.producer.education ?? '',
        }
      : emptyProducer,
  }
}

function formDataToPayload(form: FarmData) {
  const hasProducerData = Object.entries(form.producer).some(([key, value]) => key !== 'gender' && value !== '')

  return {
    name: form.name,
    owner: form.owner,
    location: form.location,
    total_area_ha: parseFloat(form.total_area) || 0,
    planted_area_ha: parseFloat(form.planted_area) || 0,
    phone: form.phone,
    activities: form.activities,
    main_crop: form.main_crop,
    secondary_crop: form.secondary_crop,
    state: form.state,
    city: form.city,
    address: form.address,
    production_system: form.production_system,
    commercialization_product: form.commercialization_product,
    has_no_cnpj: form.has_no_cnpj,
    cnpj: form.cnpj,
    has_no_nirf: form.has_no_nirf,
    nirf: form.nirf,
    has_no_incra: form.has_no_incra,
    incra_registration: form.incra_registration,
    has_no_state_registration: form.has_no_state_registration,
    state_registration: form.state_registration,
    has_no_dap: form.has_no_dap,
    dap: form.dap,
    has_no_car: form.has_no_car,
    car: form.car,
    fully_leased: form.fully_leased,
    land_value_per_ha: parseFloat(form.land_value_per_ha) || 0,
    dam_area_ha: parseFloat(form.dam_area) || 0,
    improvements_area_ha: parseFloat(form.improvements_area) || 0,
    roads_area_ha: parseFloat(form.roads_area) || 0,
    app_area_ha: parseFloat(form.app_area) || 0,
    legal_reserve_area_ha: parseFloat(form.legal_reserve_area) || 0,
    native_vegetation_area_ha: parseFloat(form.native_vegetation_area) || 0,
    livestock_area_not_covered_ha: parseFloat(form.livestock_area_not_covered) || 0,
    agriculture_area_not_covered_ha: parseFloat(form.agriculture_area_not_covered) || 0,
    non_agricultural_area_ha: parseFloat(form.non_agricultural_area) || 0,
    producer: hasProducerData
      ? {
          cpf: form.producer.cpf,
          name: form.producer.name,
          rg: form.producer.rg,
          issuing_body: form.producer.issuing_body,
          gender: form.producer.gender,
          birth_date: form.producer.birth_date || null,
          marital_status: form.producer.marital_status,
          phone: form.producer.phone,
          email: form.producer.email,
          education: form.producer.education,
        }
      : null,
  }
}

export function FarmEdit() {
  const { farmId } = useParams<{ farmId: string }>()
  const navigate = useNavigate()
  const isEditing = Boolean(farmId)
  const [initial, setInitial] = useState<FarmData | undefined>(undefined)
  const [loading, setLoading] = useState(isEditing)
  const [saving, setSaving] = useState(false)
  const toast = useToast()

  const load = useCallback(async () => {
    if (!farmId) return
    try {
      const farm = await apiRequest<Farm>(`/farms/${farmId}`)
      setInitial(farmToFormData(farm))
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [farmId])

  useEffect(() => { load() }, [load])

  const handleSave = async (formData: FarmData) => {
    setSaving(true)
    try {
      const payload = formDataToPayload(formData)
      if (isEditing) {
        await apiRequest(`/farms/${farmId}`, { method: 'PUT', body: payload })
      } else {
        await apiRequest('/farms', { method: 'POST', body: payload })
      }
      toast.success(isEditing ? 'Fazenda atualizada' : 'Fazenda cadastrada')
      navigate('/farms')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao salvar fazenda')
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
        <Button variant="ghost" size="sm" onClick={() => navigate('/farms')}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <h1 className="font-display text-2xl font-semibold text-foreground">{isEditing ? 'Editar Fazenda' : 'Nova Fazenda'}</h1>
      </div>

      <FarmForm
        initial={initial}
        onSave={handleSave}
        onCancel={() => navigate('/farms')}
        loading={saving}
      />
    </div>
  )
}
