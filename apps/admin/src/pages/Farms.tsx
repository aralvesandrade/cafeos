import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { FarmList } from '@/components/farms/FarmList'
import { FarmForm, emptyProducer, type FarmData } from '@/components/farms/FarmForm'
import { Plus } from 'lucide-react'

interface Producer {
  cpf: string
  name: string
  rg: string
  issuing_body: string
  gender: string
  birth_date: string | null
  marital_status: string
  phone: string
  email: string
  education: string
}

export interface Farm {
  id: string
  name: string
  owner: string
  location: string
  total_area_ha: number
  planted_area_ha: number

  phone: string
  activities: string
  main_crop: string
  secondary_crop: string
  state: string
  city: string
  address: string
  production_system: string
  commercialization_product: string

  has_no_cnpj: boolean
  cnpj: string
  has_no_nirf: boolean
  nirf: string
  has_no_incra: boolean
  incra_registration: string
  has_no_state_registration: boolean
  state_registration: string
  has_no_dap: boolean
  dap: string
  has_no_car: boolean
  car: string

  fully_leased: boolean
  land_value_per_ha: number

  dam_area_ha: number
  improvements_area_ha: number
  roads_area_ha: number
  app_area_ha: number
  legal_reserve_area_ha: number
  native_vegetation_area_ha: number
  livestock_area_not_covered_ha: number
  agriculture_area_not_covered_ha: number
  non_agricultural_area_ha: number

  producer?: Producer | null
}

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

export function Farms() {
  const [farms, setFarms] = useState<Farm[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editingFarm, setEditingFarm] = useState<Farm | null>(null)
  const [saving, setSaving] = useState(false)

  const loadFarms = useCallback(async () => {
    try {
      const data = await apiRequest<Farm[]>('/farms')
      setFarms(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    loadFarms()
  }, [loadFarms])

  const handleSave = async (formData: FarmData) => {
    setSaving(true)
    try {
      const payload = formDataToPayload(formData)

      if (editingFarm) {
        await apiRequest(`/farms/${editingFarm.id}`, { method: 'PUT', body: payload })
      } else {
        await apiRequest('/farms', { method: 'POST', body: payload })
      }

      setDialogOpen(false)
      setEditingFarm(null)
      await loadFarms()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Remover fazenda?')) return
    try {
      await apiRequest(`/farms/${id}`, { method: 'DELETE' })
      await loadFarms()
    } catch (err) {
      console.error(err)
    }
  }

  const openEdit = (farm: Farm) => {
    setEditingFarm(farm)
    setDialogOpen(true)
  }

  const openCreate = () => {
    setEditingFarm(null)
    setDialogOpen(true)
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-primary">Fazendas</h1>
          <p className="text-sm text-muted-foreground">Gerencie suas propriedades rurais</p>
        </div>
        <Button onClick={openCreate}>
          <Plus className="h-4 w-4" />
          Nova Fazenda
        </Button>
      </div>

      <FarmList farms={farms} onEdit={openEdit} onDelete={handleDelete} />

      <Dialog
        open={dialogOpen}
        onClose={() => { setDialogOpen(false); setEditingFarm(null) }}
        title={editingFarm ? 'Editar Fazenda' : 'Nova Fazenda'}
        className="max-w-3xl"
      >
        <FarmForm
          initial={editingFarm ? farmToFormData(editingFarm) : undefined}
          onSave={handleSave}
          onCancel={() => { setDialogOpen(false); setEditingFarm(null) }}
          loading={saving}
        />
      </Dialog>
    </div>
  )
}
