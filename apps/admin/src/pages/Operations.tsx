import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Search, Tractor } from 'lucide-react'

interface Operation {
  id: string
  type: string
  date: string
  plot_name: string
  responsible: string
  product_used: string
  quantity: number
  cost: number
}

const typeColors: Record<string, 'info' | 'warning' | 'success' | 'default' | 'danger'> = {
  adubacao: 'info',
  pulverizacao: 'warning',
  irrigacao: 'success',
  poda: 'default',
  colheita: 'danger',
}

export function Operations() {
  const [operations, setOperations] = useState<Operation[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<Operation[]>('/operations')
      setOperations(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const filtered = operations.filter(
    (op) =>
      op.plot_name.toLowerCase().includes(search.toLowerCase()) ||
      op.type.toLowerCase().includes(search.toLowerCase()) ||
      op.responsible.toLowerCase().includes(search.toLowerCase())
  )

  if (loading) return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-coffee-green-dark">Operações</h1>
        <p className="text-sm text-coffee-text-light">Histórico de operações agrícolas</p>
      </div>

      <div className="relative w-full max-w-sm">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
        <Input
          className="pl-9"
          placeholder="Buscar por talhão, tipo, responsável..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Data</TableHeader>
            <TableHeader>Tipo</TableHeader>
            <TableHeader>Talhão</TableHeader>
            <TableHeader>Responsável</TableHeader>
            <TableHeader>Produto</TableHeader>
            <TableHeader>Quantidade</TableHeader>
            <TableHeader className="text-right">Custo</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {filtered.map((op) => (
            <TableRow key={op.id}>
              <TableCell>{op.date}</TableCell>
              <TableCell>
                <div className="flex items-center gap-2">
                  <Tractor className="h-4 w-4 text-coffee-green" />
                  <Badge variant={typeColors[op.type] || 'default'}>{op.type}</Badge>
                </div>
              </TableCell>
              <TableCell>{op.plot_name}</TableCell>
              <TableCell>{op.responsible}</TableCell>
              <TableCell>{op.product_used}</TableCell>
              <TableCell>{op.quantity}</TableCell>
              <TableCell className="text-right">R$ {op.cost.toFixed(2)}</TableCell>
            </TableRow>
          ))}
          {filtered.length === 0 && (
            <TableRow><TableCell colSpan={7} className="text-center text-coffee-text-light py-8">Nenhuma operação encontrada.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  )
}
