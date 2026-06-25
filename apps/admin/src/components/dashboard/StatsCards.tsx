import { Card, CardContent } from '@/components/ui/card'
import { Sprout, DollarSign, MapPin, Grid3X3 } from 'lucide-react'

interface Stats {
  totalFarms: number
  totalPlots: number
  totalProduction: number
  totalCost: number
}

export function StatsCards({ stats }: { stats: Stats }) {
  const items = [
    { icon: MapPin, label: 'Fazendas', value: stats.totalFarms },
    { icon: Grid3X3, label: 'Talhões', value: stats.totalPlots },
    { icon: Sprout, label: 'Produção (sacas)', value: stats.totalProduction.toLocaleString() },
    { icon: DollarSign, label: 'Custo total', value: `R$ ${stats.totalCost.toLocaleString()}` },
  ]

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {items.map((item) => (
        <Card key={item.label}>
          <CardContent className="flex items-center gap-4 p-4">
            <div className="w-10 h-10 rounded-lg bg-coffee-green/10 flex items-center justify-center text-coffee-green">
              <item.icon className="h-5 w-5" />
            </div>
            <div>
              <p className="text-xs text-coffee-text-light">{item.label}</p>
              <p className="text-xl font-bold text-coffee-green-dark">{item.value}</p>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
