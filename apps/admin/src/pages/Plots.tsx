import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Plus, Pencil, Trash2 } from 'lucide-react'

interface Plot {
  id: string
  name: string
  farm_id: string
  farm_name: string
  area: number
  cultivar: string
  soil_type: string
  altitude: number
  planting_year: number
}

interface Farm { id: string; name: string }

export function Plots() {
  const [plots, setPlots] = useState<Plot[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Plot | null>(null)
  const [saving, setSaving] = useState(false)

  const [form, setForm] = useState({
    name: '', farm_id: '', area: '', cultivar: '', soil_type: '', altitude: '', planting_year: '',
  })

  const loadData = useCallback(async () => {
    try {
      const [plotsData, farmsData] = await Promise.all([
        apiRequest<Plot[]>('/plots'),
        apiRequest<Farm[]>('/farms'),
      ])
      setPlots(plotsData)
      setFarms(farmsData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { loadData() }, [loadData])

  const handleSave = async () => {
    setSaving(true)
    try {
      const payload = {
        name: form.name,
        farm_id: form.farm_id,
        area: parseFloat(form.area),
        cultivar: form.cultivar,
        soil_type: form.soil_type,
        altitude: parseFloat(form.altitude),
        planting_year: parseInt(form.planting_year),
      }

      if (editing) {
        await apiRequest(`/plots/${editing.id}`, { method: 'PUT', body: payload })
      } else {
        await apiRequest('/plots', { method: 'POST', body: payload })
      }

      setDialogOpen(false)
      setEditing(null)
      await loadData()
    } catch (err) {
      console.error(err)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Remover talhão?')) return
    try {
      await apiRequest(`/plots/${id}`, { method: 'DELETE' })
      await loadData()
    } catch (err) { console.error(err) }
  }

  const openEdit = (plot: Plot) => {
    setEditing(plot)
    setForm({
      name: plot.name, farm_id: plot.farm_id, area: String(plot.area),
      cultivar: plot.cultivar, soil_type: plot.soil_type,
      altitude: String(plot.altitude), planting_year: String(plot.planting_year),
    })
    setDialogOpen(true)
  }

  const openCreate = () => {
    setEditing(null)
    setForm({ name: '', farm_id: '', area: '', cultivar: '', soil_type: '', altitude: '', planting_year: '' })
    setDialogOpen(true)
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-coffee-green-dark">Talhões</h1>
          <p className="text-sm text-coffee-text-light">Gerencie os talhões das suas fazendas</p>
        </div>
        <Button onClick={openCreate}><Plus className="h-4 w-4" /> Novo Talhão</Button>
      </div>

      <Table>
        <TableHead>
          <TableRow>
            <TableHeader>Nome</TableHeader>
            <TableHeader>Fazenda</TableHeader>
            <TableHeader>Área</TableHeader>
            <TableHeader>Cultivar</TableHeader>
            <TableHeader>Solo</TableHeader>
            <TableHeader>Altitude</TableHeader>
            <TableHeader>Ano Plantio</TableHeader>
            <TableHeader className="text-right">Ações</TableHeader>
          </TableRow>
        </TableHead>
        <TableBody>
          {plots.map((plot) => (
            <TableRow key={plot.id}>
              <TableCell className="font-medium">{plot.name}</TableCell>
              <TableCell>{plot.farm_name}</TableCell>
              <TableCell>{plot.area} ha</TableCell>
              <TableCell>{plot.cultivar}</TableCell>
              <TableCell>{plot.soil_type}</TableCell>
              <TableCell>{plot.altitude} m</TableCell>
              <TableCell>{plot.planting_year}</TableCell>
              <TableCell className="text-right">
                <div className="flex justify-end gap-1">
                  <Button variant="ghost" size="sm" onClick={() => openEdit(plot)}><Pencil className="h-4 w-4" /></Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDelete(plot.id)}><Trash2 className="h-4 w-4 text-red-500" /></Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
          {plots.length === 0 && (
            <TableRow><TableCell colSpan={8} className="text-center text-coffee-text-light py-8">Nenhum talhão cadastrado.</TableCell></TableRow>
          )}
        </TableBody>
      </Table>

      <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Talhão' : 'Novo Talhão'}>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Nome</label>
            <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </div>
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Fazenda</label>
            <Select value={form.farm_id} onChange={(e) => setForm({ ...form, farm_id: e.target.value })} required>
              <option value="">Selecione...</option>
              {farms.map((f) => <option key={f.id} value={f.id}>{f.name}</option>)}
            </Select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-coffee-text mb-1">Área (ha)</label>
              <Input type="number" step="0.01" value={form.area} onChange={(e) => setForm({ ...form, area: e.target.value })} required />
            </div>
            <div>
              <label className="block text-sm font-medium text-coffee-text mb-1">Altitude (m)</label>
              <Input type="number" value={form.altitude} onChange={(e) => setForm({ ...form, altitude: e.target.value })} required />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-coffee-text mb-1">Cultivar</label>
            <Input value={form.cultivar} onChange={(e) => setForm({ ...form, cultivar: e.target.value })} required />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-coffee-text mb-1">Tipo de Solo</label>
              <Select value={form.soil_type} onChange={(e) => setForm({ ...form, soil_type: e.target.value })} required>
                <option value="">Selecione...</option>
                <option value="argiloso">Argiloso</option>
                <option value="arenoso">Arenoso</option>
                <option value="siltoso">Siltoso</option>
                <option value="organico">Orgânico</option>
              </Select>
            </div>
            <div>
              <label className="block text-sm font-medium text-coffee-text mb-1">Ano Plantio</label>
              <Input type="number" value={form.planting_year} onChange={(e) => setForm({ ...form, planting_year: e.target.value })} required />
            </div>
          </div>
          <div className="flex justify-end gap-3 pt-2">
            <Button type="button" variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
            <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
          </div>
        </div>
      </Dialog>
    </div>
  )
}
