import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { StatsCards } from '@/components/dashboard/StatsCards'
import { ProductionChart } from '@/components/dashboard/ProductionChart'
import { CostChart } from '@/components/dashboard/CostChart'
import { RecentOperations } from '@/components/dashboard/RecentOperations'
import { Select } from '@/components/ui/select'

interface DashboardData {
  total_farms: number
  total_plots: number
  total_production: number
  total_cost: number
  production_by_harvest: { year: string; production: number }[]
  cost_per_bag: { year: string; cost: number }[]
  recent_operations: { id: string; type: string; date: string; plot_name: string; cost: number }[]
}

interface Farm { id: string; name: string }
interface Harvest { id: string; year: string; description: string }

export function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [farms, setFarms] = useState<Farm[]>([])
  const [harvests, setHarvests] = useState<Harvest[]>([])
  const [farmId, setFarmId] = useState('')
  const [harvestId, setHarvestId] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      apiRequest<Farm[]>('/farms'),
      apiRequest<Harvest[]>('/harvests'),
    ])
      .then(([farmsData, harvestsData]) => {
        setFarms(farmsData)
        setHarvests(harvestsData)
      })
      .catch(console.error)
  }, [])

  const loadDashboard = useCallback(() => {
    setLoading(true)
    const params: Record<string, string> = {}
    if (farmId) params.farm_id = farmId
    if (harvestId) params.harvest_id = harvestId
    apiRequest<DashboardData>('/dashboard', { params })
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [farmId, harvestId])

  useEffect(() => { loadDashboard() }, [loadDashboard])

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-primary">Dashboard</h1>
          <p className="text-sm text-muted-foreground">Visão geral da sua produção</p>
        </div>
        <div className="flex gap-3">
          <Select value={farmId} onChange={(e) => setFarmId(e.target.value)} className="w-48">
            <option value="">Todas as fazendas</option>
            {farms.map((f) => (
              <option key={f.id} value={f.id}>{f.name}</option>
            ))}
          </Select>
          <Select value={harvestId} onChange={(e) => setHarvestId(e.target.value)} className="w-48">
            <option value="">Todas as safras</option>
            {harvests.map((h) => (
              <option key={h.id} value={h.id}>{h.description}</option>
            ))}
          </Select>
        </div>
      </div>

      {loading || !data ? (
        <div className="flex items-center justify-center h-64 text-muted-foreground">
          {loading ? 'Carregando...' : 'Erro ao carregar dashboard.'}
        </div>
      ) : (
        <>
          <StatsCards
            stats={{
              totalFarms: data.total_farms,
              totalPlots: data.total_plots,
              totalProduction: data.total_production,
              totalCost: data.total_cost,
            }}
          />

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <ProductionChart data={data.production_by_harvest} />
            <CostChart data={data.cost_per_bag} />
          </div>

          <RecentOperations operations={data.recent_operations} />
        </>
      )}
    </div>
  )
}
