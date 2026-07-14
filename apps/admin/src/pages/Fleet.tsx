import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Select } from '@/components/ui/select'
import { Plus, Pencil, Trash2, Truck, Wrench } from 'lucide-react'

interface Vehicle {
  id: string; name: string; farm_id: string | null; type: string; plate: string; brand: string; model: string; year: number; status: string
}

interface Maintenance {
  id: string; vehicle_id: string; vehicle: { name: string }; date: string; type: string; description: string; cost: number; odometer: number
}

interface Farm { id: string; name: string }

const typeLabels: Record<string, string> = { trator: 'Trator', pulverizador: 'Pulverizador', caminhao: 'Caminhão', utilitario: 'Utilitário', outro: 'Outro' }
const statusLabels: Record<string, string> = { active: 'Ativo', maintenance: 'Em Manutenção', retired: 'Baixado' }

export function Fleet() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [maintenance, setMaintenance] = useState<Maintenance[]>([])
  const [farms, setFarms] = useState<Farm[]>([])
  const [farmFilter, setFarmFilter] = useState('')
  const [loading, setLoading] = useState(true)
  const [tab, setTab] = useState<'vehicles' | 'maintenance'>('vehicles')
  const [dialogOpen, setDialogOpen] = useState(false)
  const [maintDialogOpen, setMaintDialogOpen] = useState(false)
  const [editing, setEditing] = useState<Vehicle | null>(null)
  const [form, setForm] = useState({ name: '', farm_id: '', type: 'trator', plate: '', brand: '', model: '', year: '' })
  const [maintForm, setMaintForm] = useState({ vehicle_id: '', type: 'preventive', description: '', cost: '', odometer: '', date: '', notes: '' })
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const [vData, mData, farmsData] = await Promise.all([
        apiRequest<Vehicle[]>('/fleet/vehicles', { params: farmFilter ? { farm_id: farmFilter } : undefined }),
        apiRequest<Maintenance[]>('/fleet/maintenance', { params: farmFilter ? { farm_id: farmFilter } : undefined }),
        apiRequest<Farm[]>('/farms'),
      ])
      setVehicles(vData); setMaintenance(mData); setFarms(farmsData)
    } catch (err) { console.error(err) }
    finally { setLoading(false) }
  }, [farmFilter])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      const body = { ...form, farm_id: form.farm_id || null, year: parseInt(form.year) || 0 }
      if (editing) await apiRequest(`/fleet/vehicles/${editing.id}`, { method: 'PUT', body })
      else await apiRequest('/fleet/vehicles', { method: 'POST', body })
      setDialogOpen(false); setEditing(null); await load()
    } catch (err) { console.error(err) }
    finally { setSaving(false) }
  }

  const handleMaintSave = async () => {
    setSaving(true)
    try {
      await apiRequest('/fleet/maintenance', { method: 'POST', body: { ...maintForm, cost: parseFloat(maintForm.cost), odometer: parseFloat(maintForm.odometer) } })
      setMaintDialogOpen(false); await load()
    } catch (err) { console.error(err) }
    finally { setSaving(false) }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between flex-wrap gap-4">
      <div><h1 className="text-2xl font-bold text-primary">Frota</h1><p className="text-sm text-muted-foreground">Veículos e manutenções</p></div>
      <div className="flex gap-2 flex-wrap">
        <Select value={farmFilter} onChange={(e) => setFarmFilter(e.target.value)} className="w-48">
          <option value="">Todas as fazendas</option>
          {farms.map((f) => (
            <option key={f.id} value={f.id}>{f.name}</option>
          ))}
        </Select>
        <Button variant="outline" onClick={() => { setMaintForm({ vehicle_id: '', type: 'preventive', description: '', cost: '', odometer: '', date: '', notes: '' }); setMaintDialogOpen(true) }}><Wrench className="h-4 w-4" /> Nova Manutenção</Button>
        <Button onClick={() => { setEditing(null); setForm({ name: '', farm_id: '', type: 'trator', plate: '', brand: '', model: '', year: '' }); setDialogOpen(true) }}><Plus className="h-4 w-4" /> Novo Veículo</Button>
      </div>
    </div>
    <div className="flex gap-2 mb-2">
      <Button variant={tab === 'vehicles' ? 'primary' : 'outline'} size="sm" onClick={() => setTab('vehicles')}>Veículos</Button>
      <Button variant={tab === 'maintenance' ? 'primary' : 'outline'} size="sm" onClick={() => setTab('maintenance')}>Manutenções</Button>
    </div>
    {tab === 'vehicles' ? (<Table>
      <TableHead><TableRow><TableHeader>Nome</TableHeader><TableHeader>Fazenda</TableHeader><TableHeader>Tipo</TableHeader><TableHeader>Placa</TableHeader><TableHeader>Marca/Modelo</TableHeader><TableHeader>Ano</TableHeader><TableHeader>Status</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{vehicles.map((v) => (<TableRow key={v.id}>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><Truck className="h-4 w-4 text-primary" />{v.name}</div></TableCell>
        <TableCell className="text-muted-foreground text-sm">{farms.find(f => f.id === v.farm_id)?.name || '-'}</TableCell>
        <TableCell>{typeLabels[v.type] || v.type}</TableCell><TableCell>{v.plate || '-'}</TableCell>
        <TableCell className="text-sm">{(v.brand || '') + ' ' + (v.model || '')}</TableCell>
        <TableCell>{v.year || '-'}</TableCell>
        <TableCell><Badge variant={v.status === 'maintenance' ? 'warning' : 'success'}>{statusLabels[v.status] || v.status}</Badge></TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditing(v); setForm({ name: v.name, farm_id: v.farm_id || '', type: v.type, plate: v.plate, brand: v.brand, model: v.model, year: String(v.year) }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => { if (confirm('Remover veículo?')) apiRequest(`/fleet/vehicles/${v.id}`, { method: 'DELETE' }).then(load) }}><Trash2 className="h-4 w-4 text-destructive" /></Button>
        </div></TableCell>
      </TableRow>))}
      {vehicles.length === 0 && <TableRow><TableCell colSpan={8} className="text-center text-muted-foreground py-8">Nenhum veículo cadastrado.</TableCell></TableRow>}
      </TableBody>
    </Table>) : (<Table>
      <TableHead><TableRow><TableHeader>Veículo</TableHeader><TableHeader>Data</TableHeader><TableHeader>Tipo</TableHeader><TableHeader>Descrição</TableHeader><TableHeader>Custo</TableHeader><TableHeader>Horímetro</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{maintenance.map((m) => (<TableRow key={m.id}>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><Truck className="h-4 w-4 text-primary" />{m.vehicle?.name || m.vehicle_id}</div></TableCell>
        <TableCell className="text-sm">{m.date}</TableCell>
        <TableCell><Badge variant={m.type === 'corrective' ? 'danger' : 'default'}>{m.type === 'preventive' ? 'Preventiva' : 'Corretiva'}</Badge></TableCell>
        <TableCell className="text-sm">{m.description}</TableCell>
        <TableCell>R$ {m.cost.toFixed(2)}</TableCell>
        <TableCell className="text-sm">{m.odometer}h</TableCell>
        <TableCell className="text-right"><Button variant="ghost" size="sm" onClick={() => { if (confirm('Remover manutenção?')) apiRequest(`/fleet/maintenance/${m.id}`, { method: 'DELETE' }).then(load) }}><Trash2 className="h-4 w-4 text-destructive" /></Button></TableCell>
      </TableRow>))}
      {maintenance.length === 0 && <TableRow><TableCell colSpan={7} className="text-center text-muted-foreground py-8">Nenhuma manutenção registrada.</TableCell></TableRow>}
      </TableBody>
    </Table>)}
    <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Veículo' : 'Novo Veículo'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Nome</label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Fazenda</label>
          <Select value={form.farm_id} onChange={(e) => setForm({ ...form, farm_id: e.target.value })}>
            <option value="">Sem fazenda vinculada</option>
            {farms.map((f) => <option key={f.id} value={f.id}>{f.name}</option>)}
          </Select></div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Tipo</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
              <option value="trator">Trator</option><option value="pulverizador">Pulverizador</option><option value="caminhao">Caminhão</option><option value="utilitario">Utilitário</option><option value="outro">Outro</option>
            </select></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Placa</label><Input value={form.plate} onChange={(e) => setForm({ ...form, plate: e.target.value })} /></div>
        </div>
        <div className="grid grid-cols-3 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Marca</label><Input value={form.brand} onChange={(e) => setForm({ ...form, brand: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Modelo</label><Input value={form.model} onChange={(e) => setForm({ ...form, model: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Ano</label><Input type="number" value={form.year} onChange={(e) => setForm({ ...form, year: e.target.value })} /></div>
        </div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>
    <Dialog open={maintDialogOpen} onClose={() => setMaintDialogOpen(false)} title="Nova Manutenção">
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Veículo</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={maintForm.vehicle_id} onChange={(e) => setMaintForm({ ...maintForm, vehicle_id: e.target.value })}>
            <option value="">Selecione</option>{vehicles.map((v) => <option key={v.id} value={v.id}>{v.name} - {v.plate}</option>)}
          </select></div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Tipo</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={maintForm.type} onChange={(e) => setMaintForm({ ...maintForm, type: e.target.value })}>
              <option value="preventive">Preventiva</option><option value="corrective">Corretiva</option>
            </select></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Data</label><Input type="date" value={maintForm.date} onChange={(e) => setMaintForm({ ...maintForm, date: e.target.value })} /></div>
        </div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Descrição</label><Input value={maintForm.description} onChange={(e) => setMaintForm({ ...maintForm, description: e.target.value })} /></div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Custo (R$)</label><Input type="number" step="0.01" value={maintForm.cost} onChange={(e) => setMaintForm({ ...maintForm, cost: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Horímetro (h)</label><Input type="number" step="0.1" value={maintForm.odometer} onChange={(e) => setMaintForm({ ...maintForm, odometer: e.target.value })} /></div>
        </div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => setMaintDialogOpen(false)}>Cancelar</Button>
          <Button onClick={handleMaintSave} disabled={saving}>{saving ? 'Salvando...' : 'Registrar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
