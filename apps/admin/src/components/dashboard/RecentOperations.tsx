import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tractor } from 'lucide-react'

interface Operation {
  id: string
  type: string
  date: string
  plot_name: string
  cost: number
}

const typeColors: Record<string, 'info' | 'warning' | 'success' | 'default'> = {
  adubacao: 'info',
  pulverizacao: 'warning',
  irrigacao: 'success',
}

export function RecentOperations({ operations }: { operations: Operation[] }) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Operações Recentes</CardTitle>
      </CardHeader>
      <CardContent>
        {operations.length === 0 ? (
          <p className="text-sm text-muted-foreground">Nenhuma operação recente.</p>
        ) : (
          <div className="space-y-3">
            {operations.map((op) => (
              <div key={op.id} className="flex items-center justify-between py-2 border-b border-border last:border-0">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-full bg-muted flex items-center justify-center">
                    <Tractor className="h-4 w-4 text-primary" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-foreground">{op.plot_name}</p>
                    <p className="text-xs text-muted-foreground">{op.date}</p>
                  </div>
                </div>
                <div className="text-right">
                  <Badge variant={typeColors[op.type] || 'default'}>{op.type}</Badge>
                  <p className="text-xs text-muted-foreground mt-1">R$ {op.cost.toFixed(2)}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
