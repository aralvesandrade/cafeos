import { useEffect, useState, useCallback } from 'react'
import { apiRequest } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Dialog } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Plus, Pencil, Trash2, CircleDollarSign } from 'lucide-react'

interface CostCenter {
  id: string
  name: string
  code: string
  type: 'receita' | 'despesa'
  cost_group: string
  description: string
}

interface SenarCategory {
  name: string
  cost_group: string
}

const typeLabels: Record<string, string> = { receita: 'Receita', despesa: 'Despesa' }
const typeVariants: Record<string, 'success' | 'danger'> = { receita: 'success', despesa: 'danger' }

const costGroupLabels: Record<string, string> = {
  operacional_efetivo: 'Operacional Efetivo (COE)',
  mao_de_obra_familiar: 'Mão de Obra Familiar (COT)',
  capital_depreciacao: 'Capital / Depreciação (COT)',
  remuneracao_capital: 'Remuneração do Capital (CT)',
}

const emptyForm = { name: '', code: '', type: 'despesa', cost_group: '', description: '' }

export function CostCenters() {
  const [items, setItems] = useState<CostCenter[]>([])
  const [senarCategories, setSenarCategories] = useState<SenarCategory[]>([])
  const [loading, setLoading] = useState(true)
  const [dialogOpen, setDialogOpen] = useState(false)
  const [editing, setEditing] = useState<CostCenter | null>(null)
  const [form, setForm] = useState(emptyForm)
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    try {
      const [ccData, senarData] = await Promise.all([
        apiRequest<CostCenter[]>('/cost-centers'),
        apiRequest<SenarCategory[]>('/cost-centers/senar-categories'),
      ])
      setItems(ccData)
      setSenarCategories(senarData)
    } catch (err) { console.error(err) }
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

  const applySenarCategory = (name: string) => {
    const cat = senarCategories.find((c) => c.name === name)
    if (!cat) {
      setForm({ ...form, name })
      return
    }
    setForm({ ...form, name: cat.name, cost_group: cat.cost_group, type: 'despesa' })
  }

  if (loading) return <div className="flex items-center justify-center h-64 text-coffee-text-light">Carregando...</div>

  return (<div className="space-y-6">
    <div className="flex items-center justify-between">
      <div><h1 className="text-2xl font-bold text-coffee-green-dark">Centros de Custo</h1><p className="text-sm text-coffee-text-light">Plano de contas para receitas e despesas</p></div>
      <Button onClick={() => { setEditing(null); setForm(emptyForm); setDialogOpen(true) }}><Plus className="h-4 w-4" /> Novo Centro de Custo</Button>
    </div>
    <Table>
      <TableHead><TableRow><TableHeader>Código</TableHeader><TableHeader>Nome</TableHeader><TableHeader>Tipo</TableHeader><TableHeader>Classificação SENAR</TableHeader><TableHeader>Descrição</TableHeader><TableHeader className="text-right">Ações</TableHeader></TableRow></TableHead>
      <TableBody>{items.map((cc) => (<TableRow key={cc.id}>
        <TableCell className="font-mono text-sm">{cc.code}</TableCell>
        <TableCell className="font-medium"><div className="flex items-center gap-2"><CircleDollarSign className="h-4 w-4 text-coffee-green" />{cc.name}</div></TableCell>
        <TableCell><Badge variant={typeVariants[cc.type]}>{typeLabels[cc.type]}</Badge></TableCell>
        <TableCell className="text-sm">{cc.cost_group ? costGroupLabels[cc.cost_group] ?? cc.cost_group : <span className="text-coffee-text-light">Não classificado</span>}</TableCell>
        <TableCell className="text-coffee-text-light text-sm">{cc.description}</TableCell>
        <TableCell className="text-right"><div className="flex justify-end gap-1">
          <Button variant="ghost" size="sm" onClick={() => { setEditing(cc); setForm({ name: cc.name, code: cc.code, type: cc.type, cost_group: cc.cost_group ?? '', description: cc.description }); setDialogOpen(true) }}><Pencil className="h-4 w-4" /></Button>
          <Button variant="ghost" size="sm" onClick={() => handleDelete(cc.id)}><Trash2 className="h-4 w-4 text-red-500" /></Button>
        </div></TableCell>
      </TableRow>))}
      {items.length === 0 && <TableRow><TableCell colSpan={6} className="text-center text-coffee-text-light py-8">Nenhum centro de custo cadastrado.</TableCell></TableRow>}
      </TableBody>
    </Table>
    <Dialog open={dialogOpen} onClose={() => { setDialogOpen(false); setEditing(null) }} title={editing ? 'Editar Centro de Custo' : 'Novo Centro de Custo'}>
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-coffee-text mb-1">Categoria SENAR (opcional)</label>
          <Select value="" onChange={(e) => applySenarCategory(e.target.value)}>
            <option value="">Nome personalizado...</option>
            {senarCategories.map((c) => <option key={c.name} value={c.name}>{c.name}</option>)}
          </Select>
          <p className="text-xs text-coffee-text-light mt-1">Selecionar uma categoria preenche nome e classificação automaticamente.</p>
        </div>
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Nome</label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required /></div>
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Código</label><Input value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} placeholder="EX: DESP_ADUBOS" required /></div>
        <div className="grid grid-cols-2 gap-4">
          <div><label className="block text-sm font-medium text-coffee-text mb-1">Tipo</label>
            <Select value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
              <option value="despesa">Despesa</option><option value="receita">Receita</option>
            </Select></div>
          <div><label className="block text-sm font-medium text-coffee-text mb-1">Classificação (COE/COT/CT)</label>
            <Select value={form.cost_group} onChange={(e) => setForm({ ...form, cost_group: e.target.value })}>
              <option value="">Não classificado</option>
              {Object.entries(costGroupLabels).map(([value, label]) => <option key={value} value={value}>{label}</option>)}
            </Select></div>
        </div>
        <div><label className="block text-sm font-medium text-coffee-text mb-1">Descrição</label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
        <div className="flex justify-end gap-3 pt-2">
          <Button variant="outline" onClick={() => { setDialogOpen(false); setEditing(null) }}>Cancelar</Button>
          <Button onClick={handleSave} disabled={saving}>{saving ? 'Salvando...' : 'Salvar'}</Button>
        </div>
      </div>
    </Dialog>
  </div>)
}
