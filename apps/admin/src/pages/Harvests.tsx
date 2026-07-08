import { useEffect, useState, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Calendar, ExternalLink } from 'lucide-react'

interface Harvest {
  id: string
  year: string
  estimated_production: number
  status: string
  created_at: string
}

export function Harvests() {
  const [harvests, setHarvests] = useState<Harvest[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [year, setYear] = useState('')
  const [estimated, setEstimated] = useState('')
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const data = await apiRequest<Harvest[]>('/harvests')
      setHarvests(data)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const handleCreate = async () => {
    setSaving(true)
    try {
      await apiRequest('/harvests', {
        method: 'POST',
        body: { year, estimated_production: parseFloat(estimated) },
      })
      setDialogOpen(false)
      setYear('')
      setEstimated('')
      await load()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleFinalize = async (id: string) => {
    if (!confirm('Finalizar safra?')) return
    try {
      await apiRequest(`/harvests/${id}/finalize`, { method: 'PUT' })
      await load()
    } catch (err) {
      console.error(err)
    }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-primary">Safras</h1>
          <p className="text-sm text-muted-foreground">Gerencie as safras da sua propriedade</p>
        </div>
        <Button onClick={() => setDialogOpen(true)}><Plus className="h-4 w-4" /> Nova Safra</Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Safra</TableHeader>
            <TableHeader>Estimativa (sacas)</TableHeader>
            <TableHeader>Status</TableHeader>
            <TableHeader>Criada em</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {harvests.map((h) => (
            <TableRow key={h.id}>
              <TableCell className="font-medium">
                <Link to={`/harvests/${h.id}`} className="flex items-center gap-2 text-primary hover:underline">
                  <Calendar className="h-4 w-4" />
                  {h.year}
                </Link>
              </TableCell>
              <TableCell>{h.estimated_production}</TableCell>
              <TableCell>
                <Badge variant={h.status === 'finalized' ? 'success' : 'default'}>
                  {h.status === 'finalized' ? 'Finalizada' : 'Ativa'}
                </Badge>
              </TableCell>
              <TableCell>{h.created_at}</TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Button variant="ghost" size="sm" asChild>
                    <Link to={`/harvests/${h.id}`}>
                      <ExternalLink className="h-4 w-4" />
                    </Link>
                  </Button>
                  {h.status !== 'finalized' && (
                    <Button variant="outline" size="sm" onClick={() => handleFinalize(h.id)}>
                      Finalizar
                    </Button>
                  )}
                </div>
              </TableCell>
            </TableRow>
          ))}
          {harvests.length === 0 && (
            <TableRow><TableCell colSpan={5} className="text-center text-muted-foreground py-8">Nenhuma safra cadastrada.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setYear(''); setEstimated('') }} title="Nova Safra">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Ano da Safra</label>
            <Input type="text" value={year} onChange={(e) => setYear(e.target.value)} placeholder="2025/2026" required />
          </div>
          <div>
            <label className="block text-sm font-medium text-foreground mb-1">Estimativa de Produção (sacas)</label>
            <Input type="number" value={estimated} onChange={(e) => setEstimated(e.target.value)} required />
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button variant="outline" onClick={() => { setDialogOpen(false); setYear(''); setEstimated('') }}>Cancelar</Button>
            <Button onClick={handleCreate} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
