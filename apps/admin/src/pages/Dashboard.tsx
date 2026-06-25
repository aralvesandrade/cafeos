import { useEffect, useState } from 'react'
import { apiRequest } from '@/lib/api'
import { StatsCards } from '@/components/dashboard/StatsCards'
import { ProductionChart } from '@/components/dashboard/ProductionChart'
import { CostChart } from '@/components/dashboard/CostChart'
import { RecentOperations } from '@/components/dashboard/RecentOperations'

interface DashboardData {
  total_farms: number
  total_plots: number
  total_production: number
  total_cost: number
  production_by_harvest: { year: string; production: number }[]
  cost_per_bag: { year: string; cost: number }[]
  recent_operations: { id: string; type: string; date: string; plot_name: string; cost: number }[]
}

export function Dashboard() {
  const [data, setData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    apiRequest<DashboardData>('/dashboard')
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64 text-coffee-text-light">
        Carregando...
      </div>
    )
  }

  if (!data) {
    return (
      <div className="flex items-center justify-center h-64 text-red-600">
        Erro ao carregar dashboard.
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-coffee-green-dark">Dashboard</h1>
        <p className="text-sm text-coffee-text-light">Visão geral da sua produção</p>
      </div>

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
    </div>
  )
}
