import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { FarmList } from '@/components/farms/FarmList'
import { Plus } from 'lucide-react'
import { useToast } from '@/lib/toast'
import { useConfirm } from '@/lib/confirm'

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

export function Farms() {
  const [farms, setFarms] = useState<Farm[]>([])
  const [loading, setLoading] = useState(true)
  const toast = useToast()
  const confirm = useConfirm()

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

  const handleDelete = async (id: string) => {
    if (!(await confirm({ title: 'Remover fazenda?', variant: 'danger' }))) return
    try {
      await apiRequest(`/farms/${id}`, { method: 'DELETE' })
      await loadFarms()
      toast.success('Fazenda removida')
    } catch (err) {
      console.error(err)
      toast.error('Erro ao remover fazenda')
    }
  }

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="font-display text-2xl font-semibold text-foreground">Fazendas</h1>
          <p className="text-sm text-muted-foreground">Gerencie suas propriedades rurais</p>
        </div>
        <Button asChild>
          <Link to="/farms/new">
            <Plus className="h-4 w-4" />
            Nova Fazenda
          </Link>
        </Button>
      </div>

      <FarmList farms={farms} onDelete={handleDelete} />
    </div>
  )
}
