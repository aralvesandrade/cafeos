import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, Users, User, Clock } from 'lucide-react'

interface Team { id: string; name: string; leader: string; description: string }
interface Worker { id: string; team_id: string; team: { name: string }; name: string; role: string; phone: string; hourly_rate: number; is_active: boolean }
interface WorkShift { id: string; worker_id: string; worker: { name: string }; date: string; hours: number; cost: number; notes: string }

export function Labor() {
  const [tab, setTab] = useState<'teams' | 'workers' | 'shifts'>('teams')
  const [teams, setTeams] = useState<Team[]>([])
  const [workers, setWorkers] = useState<Worker[]>([])
  const [shifts, setShifts] = useState<WorkShift[]>([])
  const [loading, setLoading] = useState(true)

  const [teamDialog, setTeamDialog] = useState(false)
  const [workerDialog, setWorkerDialog] = useState(false)
  const [shiftDialog, setShiftDialog] = useState(false)
  const [editingTeam, setEditingTeam] = useState<Team | null>(null)
  const [editingWorker, setEditingWorker] = useState<Worker | null>(null)
  const [teamForm, setTeamForm] = useState({ name: '', leader: '', description: '' })
  const [workerForm, setWorkerForm] = useState({ team_id: '', name: '', role: '', phone: '', hourly_rate: '', is_active: true })
  const [shiftForm, setShiftForm] = useState({ worker_id: '', hours: '', cost: '', date: '', notes: '' })
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const [tData, wData, sData] = await Promise.all([
        apiRequest<Team[]>('/labor/teams'),
        apiRequest<Worker[]>('/labor/workers'),
        apiRequest<WorkShift[]>('/labor/shifts'),
      ])
      setTeams(tData); setWorkers(wData); setShifts(sData)
    } catch (err) { console.error(err) }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { load() }, [load])

  const saveTeam = async () => {
    setSaving(true)
    try {
      if (editingTeam) await apiRequest(`/labor/teams/${editingTeam.id}`, { method: 'PUT', body: teamForm })
      else await apiRequest('/labor/teams', { method: 'POST', body: teamForm })
      setTeamDialog(false); setEditingTeam(null); await load()
    } catch (err) { console.error(err) }
    finally { setSaving(false) }
  }

  const saveWorker = async () => {
    setSaving(true)
    try {
      const body = { ...workerForm, hourly_rate: parseFloat(workerForm.hourly_rate) }
      if (editingWorker) await apiRequest(`/labor/workers/${editingWorker.id}`, { method: 'PUT', body })
      else await apiRequest('/labor/workers', { method: 'POST', body })
      setWorkerDialog(false); setEditingWorker(null); await load()
    } catch (err) { console.error(err) }
    finally { setSaving(false) }
  }

  const saveShift = async () => {
    setSaving(true)
    try {
      await apiRequest('/labor/shifts', { method: 'POST', body: { ...shiftForm, hours: parseFloat(shiftForm.hours), cost: parseFloat(shiftForm.cost) } })
      setShiftDialog(false); await load()
    } catch (err) { console.error(err) }
    finally { setSaving(false) }
  }

  const workerTotals: Record<string, { hours: number; cost: number }> = {}
  for (const s of shifts) {
    const t = workerTotals[s.worker_id] || { hours: 0, cost: 0 }
    t.hours += s.hours
    t.cost += s.cost
    workerTotals[s.worker_id] = t
  }

  const teamTotals: Record<string, { hours: number; cost: number }> = {}
  for (const w of workers) {
    if (!w.team_id) continue
    const wt = workerTotals[w.id] || { hours: 0, cost: 0 }
    const t = teamTotals[w.team_id] || { hours: 0, cost: 0 }
    t.hours += wt.hours
    t.cost += wt.cost
    teamTotals[w.team_id] = t
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between">
      <div><h1 className="text-2xl font-bold text-primary">Equipes</h1><p className="text-sm text-muted-foreground">Gestão de mão de obra</p></div>
      {tab === 'teams' && <Button onClick={() => { setEditingTeam(null); setTeamForm({ name: '', leader: '', description: '' }); setTeamDialog(true) }}><Plus className="h-4 w-4" /> Nova Equipe</Button>}
      {tab === 'workers' && <Button onClick={() => { setEditingWorker(null); setWorkerForm({ team_id: '', name: '', role: '', phone: '', hourly_rate: '', is_active: true }); setWorkerDialog(true) }}><Plus className="h-4 w-4" /> Novo Trabalhador</Button>}
      {tab === 'shifts' && <Button onClick={() => { setShiftForm({ worker_id: '', hours: '', cost: '', date: '', notes: '' }); setShiftDialog(true) }}><Plus className="h-4 w-4" /> Registrar Hora</Button>}
    </div>
    <div className="flex gap-2 mb-2">
      <Button variant={tab === 'teams' ? 'primary' : 'outline'} size="sm" onClick={() => setTab('teams')}><Users className="h-4 w-4" /> Equipes</Button>
      <Button variant={tab === 'workers' ? 'primary' : 'outline'} size="sm" onClick={() => setTab('workers')}><User className="h-4 w-4" /> Trabalhadores</Button>
      <Button variant={tab === 'shifts' ? 'primary' : 'outline'} size="sm" onClick={() => setTab('shifts')}><Clock className="h-4 w-4" /> Apontamento</Button>
    </div>

    {tab === 'teams' && (<Table>
      <TableHead><TableRow><TableHeader>Nome</TableHeader><TableHeader>Líder</TableHeader><TableHeader>Descrição</TableHeader><TableHeader className="text-right">Horas</TableHeader><TableHeader className="text-right">Custo Total</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{teams.map((t) => (<TableRow key={t.id}>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><Users className="h-4 w-4 text-primary" />{t.name}</div></TableCell>
        <TableCell>{t.leader || '-'}</TableCell><TableCell className="text-sm text-muted-foreground">{t.description || '-'}</TableCell>
        <TableCell className="text-right">{(teamTotals[t.id]?.hours || 0)}h</TableCell>
        <TableCell className="text-right">R$ {(teamTotals[t.id]?.cost || 0).toFixed(2)}</TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditingTeam(t); setTeamForm({ name: t.name, leader: t.leader, description: t.description }); setTeamDialog(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => { if (confirm('Remover equipe?')) apiRequest(`/labor/teams/${t.id}`, { method: 'DELETE' }).then(load) }}><Trash2 className="h-4 w-4 text-destructive" /></Button>
        </div></TableCell>
      </TableRow>))}
      {teams.length === 0 && <TableRow><TableCell colSpan={6} className="text-center py-8 text-muted-foreground">Nenhuma equipe cadastrada.</TableCell></TableRow>}
      </TableBody>
    </Table>)}

    {tab === 'workers' && (<Table>
      <TableHead><TableRow><TableHeader>Nome</TableHeader><TableHeader>Equipe</TableHeader><TableHeader>Função</TableHeader><TableHeader>Telefone</TableHeader><TableHeader>Valor Hora</TableHeader><TableHeader className="text-right">Horas</TableHeader><TableHeader className="text-right">Custo Total</TableHeader><TableHeader>Status</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{workers.map((w) => (<TableRow key={w.id}>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><User className="h-4 w-4 text-primary" />{w.name}</div></TableCell>
        <TableCell>{w.team?.name || '-'}</TableCell><TableCell>{w.role || '-'}</TableCell><TableCell>{w.phone || '-'}</TableCell>
        <TableCell>R$ {w.hourly_rate.toFixed(2)}</TableCell>
        <TableCell className="text-right">{(workerTotals[w.id]?.hours || 0)}h</TableCell>
        <TableCell className="text-right">R$ {(workerTotals[w.id]?.cost || 0).toFixed(2)}</TableCell>
        <TableCell><Badge variant={w.is_active ? 'success' : 'default'}>{w.is_active ? 'Ativo' : 'Inativo'}</Badge></TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditingWorker(w); setWorkerForm({ team_id: w.team_id, name: w.name, role: w.role, phone: w.phone, hourly_rate: String(w.hourly_rate), is_active: w.is_active }); setWorkerDialog(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => { if (confirm('Remover trabalhador?')) apiRequest(`/labor/workers/${w.id}`, { method: 'DELETE' }).then(load) }}><Trash2 className="h-4 w-4 text-destructive" /></Button>
        </div></TableCell>
      </TableRow>))}
      {workers.length === 0 && <TableRow><TableCell colSpan={9} className="text-center py-8 text-muted-foreground">Nenhum trabalhador cadastrado.</TableCell></TableRow>}
      </TableBody>
    </Table>)}

    {tab === 'shifts' && (<Table>
      <TableHead><TableRow><TableHeader>Trabalhador</TableHeader><TableHeader>Data</TableHeader><TableHeader>Horas</TableHeader><TableHeader>Custo</TableHeader><TableHeader>Observações</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{shifts.map((s) => (<TableRow key={s.id}>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><User className="h-4 w-4 text-primary" />{s.worker?.name || s.worker_id}</div></TableCell>
        <TableCell className="text-sm">{s.date}</TableCell><TableCell>{s.hours}h</TableCell>
        <TableCell>R$ {s.cost.toFixed(2)}</TableCell><TableCell className="text-sm text-muted-foreground">{s.notes || '-'}</TableCell>
        <TableCell className="text-right"><Button variant="ghost" size="sm" onClick={() => { if (confirm('Remover apontamento?')) apiRequest(`/labor/shifts/${s.id}`, { method: 'DELETE' }).then(load) }}><Trash2 className="h-4 w-4 text-destructive" /></Button></TableCell>
      </TableRow>))}
      {shifts.length === 0 && <TableRow><TableCell colSpan={6} className="text-center py-8 text-muted-foreground">Nenhum apontamento registrado.</TableCell></TableRow>}
      </TableBody>
    </Table>)}

    <Dialog open={teamDialog} onClose={() => { setTeamDialog(false); setEditingTeam(null) }} title={editingTeam ? 'Editar Equipe' : 'Nova Equipe'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Nome</label><Input value={teamForm.name} onChange={(e) => setTeamForm({ ...teamForm, name: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Líder</label><Input value={teamForm.leader} onChange={(e) => setTeamForm({ ...teamForm, leader: e.target.value })} /></div>
        <div><label className="block text-sm font-medium text-foreground mb-1">Descrição</label><Input value={teamForm.description} onChange={(e) => setTeamForm({ ...teamForm, description: e.target.value })} /></div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setTeamDialog(false); setEditingTeam(null) }}>Cancelar</Button>
          <Button onClick={saveTeam} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>

    <Dialog open={workerDialog} onClose={() => { setWorkerDialog(false); setEditingWorker(null) }} title={editingWorker ? 'Editar Trabalhador' : 'Novo Trabalhador'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Nome</label><Input value={workerForm.name} onChange={(e) => setWorkerForm({ ...workerForm, name: e.target.value })} required /></div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Equipe</label>
            <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={workerForm.team_id} onChange={(e) => setWorkerForm({ ...workerForm, team_id: e.target.value })}>
              <option value="">Nenhuma</option>{teams.map((t) => <option key={t.id} value={t.id}>{t.name}</option>)}
            </select></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Função</label><Input value={workerForm.role} onChange={(e) => setWorkerForm({ ...workerForm, role: e.target.value })} /></div>
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Telefone</label><Input value={workerForm.phone} onChange={(e) => setWorkerForm({ ...workerForm, phone: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Valor Hora (R$)</label><Input type="number" step="0.01" value={workerForm.hourly_rate} onChange={(e) => setWorkerForm({ ...workerForm, hourly_rate: e.target.value })} /></div>
        </div>
        {editingWorker && <div><label className="block text-sm font-medium text-foreground mb-1">Status</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={workerForm.is_active ? 'true' : 'false'} onChange={(e) => setWorkerForm({ ...workerForm, is_active: e.target.value === 'true' })}>
            <option value="true">Ativo</option><option value="false">Inativo</option>
          </select></div>}
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setWorkerDialog(false); setEditingWorker(null) }}>Cancelar</Button>
          <Button onClick={saveWorker} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>

    <Dialog open={shiftDialog} onClose={() => setShiftDialog(false)} title="Registrar Apontamento">
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-foreground mb-1">Trabalhador</label>
          <select className="flex h-10 w-full rounded-lg border border-input bg-background px-3 py-2 text-sm" value={shiftForm.worker_id} onChange={(e) => setShiftForm({ ...shiftForm, worker_id: e.target.value })}>
            <option value="">Selecione</option>{workers.filter((w) => w.is_active).map((w) => <option key={w.id} value={w.id}>{w.name}</option>)}
          </select></div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Data</label><Input type="date" value={shiftForm.date} onChange={(e) => setShiftForm({ ...shiftForm, date: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Horas</label><Input type="number" step="0.5" value={shiftForm.hours} onChange={(e) => setShiftForm({ ...shiftForm, hours: e.target.value })} /></div>
        </div>
        <div className="grid grid-cols-2 gap-3">
          <div><label className="block text-sm font-medium text-foreground mb-1">Custo (R$)</label><Input type="number" step="0.01" value={shiftForm.cost} onChange={(e) => setShiftForm({ ...shiftForm, cost: e.target.value })} /></div>
          <div><label className="block text-sm font-medium text-foreground mb-1">Observações</label><Input value={shiftForm.notes} onChange={(e) => setShiftForm({ ...shiftForm, notes: e.target.value })} /></div>
        </div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => setShiftDialog(false)}>Cancelar</Button>
          <Button onClick={saveShift} disabled={saving}>{saving ? 'Salvando...' : 'Registrar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
