import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, CircleDollarSign } from 'lucide-react'

interface CostCenter {
  id: string
  name: string
  code: string
  type: 'receita' | 'despesa'
  description: string
}

const typeLabels: Record<string, string> = { receita: 'Receita', despesa: 'Despesa' }
const typeVariants: Record<string, 'success' | 'danger'> = { receita: 'success', despesa: 'danger' }

export function CostCenters() {
  const [items, setItems] = useState<CostCenter[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<CostCenter | null>(null)
  const [form, setForm] = useState({ name: '', code: '', type: 'despesa', description: '' })
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try { const data = await apiRequest<CostCenter[]>('/cost-centers'); setItems(data) }
    catch (err) { console.error(err) }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { load() }, [load])

  const handleSave = async () => {
    setSaving(true)
    try {
      if (editing) await apiRequest(`/cost-centers/${editing.id}`, { method: 'PUT', body: form })
      else await apiRequest('/cost-centers', { method: 'POST', body: form })
      setDialogOpen(false); setEditing(null); await load()
    } catch (err) { console.error(err) }
    finally { setSaving(false) }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Remover centro de custo?')) return
    try { await apiRequest(`/cost-centers/${id}`, { method: 'DELETE' }); await load() }
    catch (err) { console.error(err) }
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between">
      <div><h1 className="text-2xl font-bold text-coffee-green-dark">Centros de Custo</h1><p className="text-sm text-coffee-text-light">Plano de contas para receitas e despesas</p></div>
      <Button onClick={() => { setEditing(null); setForm({ name: '', code: '', type: 'despesa', description: '' }); setDialogOpen(true) }}><Plus className="h-4 w-4" /> Novo Centro de Custo</Button>
    </div>
    <Table>
      <TableHead><TableRow><TableHeader>Código</TableHeader><TableHeader>Nome</TableHeader><TableHeader>Tipo</TableHeader><TableHeader>Descrição</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{items.map((cc) => (<TableRow key={cc.id}>
        <TableCell className="font-mono text-sm">{cc.code}</TableCell>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><CircleDollarSign className="h-4 w-4 text-coffee-green" />{cc.name}</div></TableCell>
        <TableCell><Badge variant={typeVariants[cc.type]}>{typeLabels[cc.type]}</Badge></TableCell>
        <TableCell className="text-coffee-text-light text-sm">{cc.description}</TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditing(cc); setForm({ name: cc.name, code: cc.code, type: cc.type, description: cc.description }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => handleDelete(cc.id)}><Trash2 className="h-4 w-4 text-red-500" /></Button>
        </div></TableCell>
      </TableRow>))}
      {items.length === 0 && <TableRow><TableCell colSpan={5} className="text-center text-coffee-text-light py-8">Nenhum centro de custo cadastrado.</TableCell></TableRow>}
      </TableBody>
    </Table>
    <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Centro de Custo' : 'Novo Centro de Custo'}>
      <div className="space-y-4">
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Nome</label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Código</label><Input value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="EX: DESP_ADUBOS" required /></div>
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Tipo</label>
          <select className="flex h-10 w-full rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-coffee-text" value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
            <option value="despesa">Despesa</option><option value="receita">Receita</option>
          </select></div>
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Descrição</label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
